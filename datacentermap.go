package dctycoon

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/supplier"
	log "github.com/sirupsen/logrus"
)

const (
	RACK_NORMAL_STATE = iota
	RACK_OVER_CURRENT = iota
	RACK_HEAT_WARNING = iota
	RACK_OVER_HEAT    = iota
	RACK_MELTING      = iota
)

type RackStatusSubscriber interface {
	// when a rack change from status (RACK_NORMAL_STATE, RACK_OVER_CURRENT, RACK_HEAT_WARNING, RACK_OVER_HEAT, RACK_MELTING)
	RackStatusChange(x, y int32, rackstate int32)

	// outage is true when there is a global outage, and false when we recover
	GeneralOutage(outage bool)
}

// DatacenterMap holds the map information, and is used to
// - show the datacenter on screen
// - show the heatmap on screen
// - spot a rack which is overused electrically
// - spot all racks when there is an electricity outage
// - spot a rack which is overtemperature (>50 degrees)
type DatacenterMap struct {
	inventory             *supplier.Inventory
	externaltemp          float64
	tiles                 [][]*Tile
	heatmap               [][]float64
	overheating           [][]int32
	width                 int32
	height                int32
	rackstatusSubscribers []RackStatusSubscriber
	inoutage              bool
}

func (self *DatacenterMap) AddRackStatusSubscriber(subscriber RackStatusSubscriber) {
	for _, s := range self.rackstatusSubscribers {
		if s == subscriber {
			return
		}
	}
	self.rackstatusSubscribers = append(self.rackstatusSubscribers, subscriber)
}

func (self *DatacenterMap) RemoveRackStatusSubscriber(subscriber RackStatusSubscriber) {
	for i, s := range self.rackstatusSubscribers {
		if s == subscriber {
			self.rackstatusSubscribers = append(self.rackstatusSubscribers[:i], self.rackstatusSubscribers[i+1:]...)
			break
		}
	}
}

func (self *DatacenterMap) GetWidth() int32 {
	return self.width
}

func (self *DatacenterMap) GetHeight() int32 {
	return self.width
}

func (self *DatacenterMap) GetTile(x, y int32) *Tile {
	return self.tiles[y][x]
}

func (self *DatacenterMap) GetTemperature(x, y int32) float64 {
	return self.heatmap[y][x]
}

func (self *DatacenterMap) GetGeneralOutage() bool {
	return self.inoutage
}

func (self *DatacenterMap) GetRackStatus(x, y int32) int32 {
	element := self.tiles[y][x].TileElement()
	if self.inoutage && element != nil {
		return RACK_OVER_CURRENT
	}
	if self.overheating[y][x] < 0 {
		return RACK_OVER_CURRENT
	}
	if self.overheating[y][x] >= 16 {
		return RACK_MELTING
	}
	if self.overheating[y][x] > 8 {
		return RACK_OVER_HEAT
	}
	if self.overheating[y][x] > 0 {
		return RACK_HEAT_WARNING
	}
	return RACK_NORMAL_STATE
}

func (self *DatacenterMap) ItemInTransit(*supplier.InventoryItem)       {}
func (self *DatacenterMap) ItemInStock(*supplier.InventoryItem)         {}
func (self *DatacenterMap) ItemRemoveFromStock(*supplier.InventoryItem) {}
func (self *DatacenterMap) ItemChangedPool(*supplier.InventoryItem)     {}

func (self *DatacenterMap) ItemInstalled(item *supplier.InventoryItem) {
	if item.Xplaced <= self.width && item.Yplaced <= self.height && item.Xplaced >= 0 && item.Yplaced >= 0 {
		self.tiles[item.Yplaced][item.Xplaced].ItemInstalled(item)
	}
}

func (self *DatacenterMap) ItemUninstalled(item *supplier.InventoryItem) {
	if item.Xplaced <= self.width && item.Yplaced <= self.height && item.Xplaced >= 0 && item.Yplaced >= 0 {
		self.tiles[item.Yplaced][item.Xplaced].ItemUninstalled(item)
	}
}

/*
func (self *DatacenterMap) GetRackItems(x, y int32) []*supplier.InventoryItem {
	items := make([]*supplier.InventoryItem, 0, 0)
	for _, i := range self.inventory.Items {
		if i.Xplaced == x && i.Yplaced == y && i.Typeitem == supplier.PRODUCT_SERVER {
			items = append(items, i)
		}
	}
	return items
}*/

// InstallItem is used when we drag&drop something into a tile
func (self *DatacenterMap) InstallItem(item *supplier.InventoryItem, x, y int32) bool {
	// if it is a rack server we drop onto a rack tower
	if item.Typeitem == supplier.PRODUCT_SERVER && item.Serverconf.ConfType.NbU > 0 {
		nbu := item.Serverconf.ConfType.NbU
		var busy [42]bool

		// first create the map of what is filled
		for _, i := range self.inventory.Items {
			if i.Xplaced == x && i.Yplaced == y && i.Typeitem == supplier.PRODUCT_SERVER {
				iNbU := i.Serverconf.ConfType.NbU
				for j := 0; j < int(iNbU); j++ {
					if j+int(i.Zplaced) < 42 {
						busy[j+int(i.Zplaced)] = true
					}
				}
			}
		}

		// try to find a place
		for i := 0; i < int(42-nbu); i++ {
			found := true
			for j := 0; j < int(nbu); j++ {
				if busy[i+j] == true {
					found = false
					break
				}
			}
			if found == true {
				self.inventory.InstallItem(item, x, y, int32(i))
				return true
			}
		}
		return false
	} else { // we drag&drop a NON rack server on a tile
		return self.inventory.InstallItem(item, x, y, -1)
	}
}

func (self *DatacenterMap) UninstallItem(item *supplier.InventoryItem) {
	self.inventory.UninstallItem(item)
}

//
// LoadMap typically load a map like:
//   {
//     "width": 10,
//     "height": 10,
//     "tiles": [
//       {"x":0, "y":0, "wall0":"","wall1":"","floor":"inside","rotation":0, "decoration":""},
//       {"x":1, "y":0, "wall0":"","wall1":"","floor":"inside","rotation":0, "decoration":"chair"}
//     ]
//   }
//
func (self *DatacenterMap) LoadMap(dc map[string]interface{}) {
	log.Debug("DatacenterMap::LoadMap(", dc, ")")
	self.width = int32(dc["width"].(float64))
	self.height = int32(dc["height"].(float64))
	self.tiles = make([][]*Tile, self.height)
	self.heatmap = make([][]float64, self.height)
	self.overheating = make([][]int32, self.height)
	self.inoutage = false
	for y := range self.tiles {
		self.tiles[y] = make([]*Tile, self.width)
		self.heatmap[y] = make([]float64, self.width)
		self.overheating[y] = make([]int32, self.width)
		for x := range self.tiles[y] {
			self.tiles[y][x] = NewGrassTile()
			self.heatmap[y][x] = self.externaltemp
		}
	}
	tiles := dc["tiles"].([]interface{})
	for _, t := range tiles {
		tile := t.(map[string]interface{})
		x := int32(tile["x"].(float64))
		y := int32(tile["y"].(float64))
		wall0 := tile["wall0"].(string)
		wall1 := tile["wall1"].(string)
		floor := tile["floor"].(string)
		rotation := uint32(tile["rotation"].(float64))
		var decorationname string
		if data, ok := tile["decoration"]; ok {
			decorationname = data.(string)
		}
		if data, ok := tile["heat"]; ok {
			self.heatmap[y][x] = data.(float64)
		}
		if data, ok := tile["overheating"]; ok {
			self.overheating[y][x] = int32(data.(float64))
		}
		self.tiles[y][x] = NewTile(wall0, wall1, floor, rotation, decorationname)
	}
	// place everything except servers
	for _, item := range self.inventory.Items {
		if item.IsPlaced() && item.Typeitem != supplier.PRODUCT_SERVER {
			self.ItemInstalled(item)
		}
	}
	// place servers
	for _, item := range self.inventory.Items {
		if item.IsPlaced() && item.Typeitem == supplier.PRODUCT_SERVER {
			self.ItemInstalled(item)
		}
	}
	self.ComputeHeatMap()
	self.ComputeOverLimits()
}

func (self *DatacenterMap) InitMap(assetdcmap string) {
	log.Debug("DatacenterMap::InitMap(", assetdcmap, ")")
	if data, err := global.Asset("assets/dcmap/" + assetdcmap); err == nil {
		var dcmap map[string]interface{}
		if json.Unmarshal(data, &dcmap) == nil {
			self.LoadMap(dcmap)
		}
	}
}

func (self *DatacenterMap) SaveMap() string {
	s := fmt.Sprintf(`{"width":%d, "height":%d, "tiles": [`, self.width, self.height)
	previous := false
	for y, _ := range self.tiles {
		for x, _ := range self.tiles[y] {
			t := self.tiles[y][x]
			value := ""
			decorationname := ""
			if t.TileElement() != nil && t.TileElement().ElementType() == supplier.PRODUCT_DECORATION {
				decoration := t.TileElement().(*DecorationElement)
				decorationname = decoration.GetName()
			}
			if t.wall[0] != "" || t.wall[1] != "" || t.floor != "green" {
				value = fmt.Sprintf(`{"x":%d, "y":%d, "wall0":"%s", "wall1":"%s", "floor":"%s","rotation":%d, "decoration": "%s", "heat":%f, "overheating": %d}`,
					x,
					y,
					t.wall[0],
					t.wall[1],
					t.floor,
					t.rotation,
					decorationname,
					self.heatmap[y][x],
					self.overheating[y][x],
				)
			}
			if value != "" {
				if previous == true {
					s += ",\n"
				}
				previous = true
				s += value
			}
		}
	}
	s += "]}"
	return s
}

func (self *DatacenterMap) SetGame(inventory *supplier.Inventory, location *supplier.LocationType, currenttime time.Time) {
	log.Debug("DatacenterMap::SetGame(", inventory, ",", location, ",", currenttime, ")")
	if self.inventory != nil {
		self.inventory.RemoveInventorySubscriber(self)
		self.inventory.RemovePowerChangeSubscriber(self)
	}
	self.inventory = inventory
	inventory.AddInventorySubscriber(self)
	inventory.AddPowerStatSubscriber(self)
	self.externaltemp = location.Temperatureaverage
}

func (self *DatacenterMap) PowerChange(time time.Time, consumed, generated, delivered, cooler float64) {
	outagesituation := consumed > generated && consumed > delivered
	if outagesituation != self.inoutage {
		self.inoutage = outagesituation
		for _, s := range self.rackstatusSubscribers {
			s.GeneralOutage(self.inoutage)
		}
	}

	if self.inoutage {
		for y := 0; y < int(self.height); y++ {
			for x := 0; x < int(self.width); x++ {
				self.heatmap[y][x] = self.externaltemp
				self.overheating[y][x] = 0
			}
		}
	}
}

func (self *DatacenterMap) getHeatSpread(x, y int32) (temperature, spreadfactor float64) {
	if y < 0 || y >= self.height || x < 0 || x >= self.width {
		return self.externaltemp, 0.002
	}
	if self.tiles[y][x].floor == "green" {
		// building wall are isolating a bit
		return self.heatmap[y][x], 0.002
	} else {
		return self.heatmap[y][x], 0.2
	}
}

func (self *DatacenterMap) ComputeHeatMap() {
	log.Debug("DatacenterMap::ComputeHeatMap()")
	// to raise 1 m-3 of air by 1 deg C : 1.211kj
	// 1kwh = 3600kJ /h?
	// 1BTU ~ 1kJ /s? or kwh?
	//
	// so
	//	self.heatmap = make([][]float64, self.height)
	nbairflow := float64(0)
	for y := 0; y < int(self.height); y++ {
		//		self.heatmap[y] = make([]float64, self.width)
		for x := 0; x < int(self.width); x++ {
			//			self.heatmap[y][x] = self.externaltemp
			if self.tiles[y][x].floor == "inside.air" {
				nbairflow++
			}
			if self.tiles[y][x].floor == "green" {
				self.heatmap[y][x] = self.externaltemp
			}
		}
	}

	// removePerAirFlow is in wH
	// GetHotspotValue() is in wH
	_, _, removePerAirflow := self.inventory.GetGlobalPower()
	if nbairflow != 0 {
		removePerAirflow = removePerAirflow / nbairflow
	}

	// we loop 60 times to spread the heat
	for loop := 0; loop < 60; loop++ {
		nextheatmap := make([][]float64, self.height)
		for y := int32(0); y < self.height; y++ {
			nextheatmap[y] = make([]float64, self.width)
			for x := int32(0); x < self.width; x++ {
				nextheatmap[y][x] = self.heatmap[y][x]
				// we flow air only inside
				if self.tiles[y][x].floor != "green" {
					// we spread the heat
					bottom, bottomfactor := self.getHeatSpread(x, y-1)
					left, leftfactor := self.getHeatSpread(x-1, y)
					top, topfactor := self.getHeatSpread(x, y+1)
					right, rightfactor := self.getHeatSpread(x+1, y)
					nextheatmap[y][x] = (1-leftfactor-rightfactor-topfactor-bottomfactor)*nextheatmap[y][x] +
						leftfactor*left + rightfactor*right + topfactor*top + bottomfactor*bottom
					// we add the heat
					nextheatmap[y][x] += self.inventory.GetHotspotValue(int32(y), int32(x)) / 1000
					// we remove via the airflow
					if self.tiles[y][x].floor == "inside.air" {
						nextheatmap[y][x] -= removePerAirflow / 1000
						if nextheatmap[y][x] < 17 {
							nextheatmap[y][x] = 17
						}
					}
				}
			}
		}
		self.heatmap = nextheatmap
	}
}

func (self *DatacenterMap) triggerRackStatus(x, y int32, status int32) {
	log.Debug("DatacenterMap::triggerRackStatus(", x, ",", y, ",", status, ")")
	for _, s := range self.rackstatusSubscribers {
		s.RackStatusChange(x, y, status)
	}
}

func (self *DatacenterMap) ComputeOverLimits() {
	log.Debug("DatacenterMap::ComputeOverLimits()")
	if self.inoutage {
		for y := int32(0); y < self.height; y++ {
			for x := int32(0); x < self.width; x++ {
				if self.overheating[y][x] > 0 {
					self.triggerRackStatus(x, y, RACK_NORMAL_STATE)
				}
				self.overheating[y][x] = 0
			}
		}
	} else {
		for y := int32(0); y < self.height; y++ {
			for x := int32(0); x < self.width; x++ {
				previousState := self.GetRackStatus(x, y)
				newState := int32(RACK_NORMAL_STATE)
				element := self.tiles[y][x].TileElement()

				if element != nil && element.ElementType() == supplier.PRODUCT_RACK {
					if self.heatmap[y][x] > 40 {
						if self.heatmap[y][x] > 45 {
							self.overheating[y][x]++
						}
						// if we over heat since 16 days
						if self.overheating[y][x] >= 16 {
							newState = RACK_MELTING
							self.overheating[y][x] = 8 // to repeat over heating in 8 days
						} else if self.overheating[y][x] > 8 {
							newState = RACK_OVER_HEAT
						} else {
							// if we begin to over heat
							newState = RACK_HEAT_WARNING
							if self.overheating[y][x] <= 0 {
								self.overheating[y][x] = 1
							}
						}
					} else {
						self.overheating[y][x] = 0
						// if we go over 64 A
						if self.inventory.GetHotspotValue(y, x) > 64.0*110.0 {
							newState = RACK_OVER_CURRENT
							self.overheating[y][x] = -1
						}
					}
				}
				if previousState != newState {
					self.triggerRackStatus(x, y, newState)
				}
			}
		}
	}
}

func (self *DatacenterMap) MoveElement(xfrom, yfrom, xto, yto int32) bool {
	tileFrom := self.GetTile(xfrom, yfrom)
	element := tileFrom.element
	tileTo := self.GetTile(xto, yto)

	if (tileTo.IsFloorOutside() && element.ElementType() == supplier.PRODUCT_GENERATOR) ||
		(tileTo.IsFloorInsideNotAirFlow() && element.ElementType() != supplier.PRODUCT_GENERATOR) {
		rotation := tileFrom.rotation
		tileFrom.rotation = 0
		tileFrom.element = nil
		tileFrom.freeSurface()

		// update the tiles
		tileTo.element = element
		tileTo.rotation = rotation
		tileTo.freeSurface()

		// update the InventoryItem
		if element != nil {
			element.InventoryItem().Xplaced = xto
			element.InventoryItem().Yplaced = yto
		}
		if element.ElementType() == supplier.PRODUCT_RACK {
			rack := element.(*RackElement)
			for _, i := range rack.items {
				i.Xplaced = xto
				i.Yplaced = yto
			}
		}

		// update the electrical/heat situation
		rackstatus := self.GetRackStatus(xfrom, yfrom)
		overheat := self.overheating[yfrom][xfrom]
		if rackstatus != RACK_NORMAL_STATE {
			self.overheating[yfrom][xfrom] = 0
			self.overheating[yto][xto] = overheat
			self.triggerRackStatus(xfrom, yfrom, RACK_NORMAL_STATE)
			self.triggerRackStatus(xto, yto, rackstatus)
		}
		self.inventory.ComputeGlobalPower()

		return true
	}
	return false
}

func NewDatacenterMap() *DatacenterMap {
	dcmap := &DatacenterMap{
		tiles:                 [][]*Tile{{}},
		heatmap:               [][]float64{},
		overheating:           [][]int32{},
		inventory:             nil,
		width:                 0,
		height:                0,
		externaltemp:          0,
		rackstatusSubscribers: make([]RackStatusSubscriber, 0, 0),
		inoutage:              false,
	}
	return dcmap
}
