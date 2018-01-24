package dctycoon

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/supplier"
	log "github.com/sirupsen/logrus"
)

// DatacenterMap holds the map information, and is used to
// - show the datacenter on screen
// - show the heatmap on screen
// - spot a rack which is overused electrically
// - spot all racks when there is an electricity outage
// - spot a rack which is overtemperature (>50 degrees)
type DatacenterMap struct {
	inventory    *supplier.Inventory
	externaltemp float64
	tiles        [][]*Tile
	heatmap      [][]float64
	width        int32
	height       int32
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
//       {"x":0, "y":0, "wall0":"","wall1":"","floor":"inside","dcelementtype":"rack","dcelement":{...}},
//       {"x":1, "y":0, "wall0":"","wall1":"","floor":"inside"},
//     ]
//   }
//
func (self *DatacenterMap) LoadMap(dc map[string]interface{}) {
	log.Debug("DatacenterMap::LoadMap(", dc, ")")
	self.width = int32(dc["width"].(float64))
	self.height = int32(dc["height"].(float64))
	self.tiles = make([][]*Tile, self.height)
	self.heatmap = make([][]float64, self.height)
	for y := range self.tiles {
		self.tiles[y] = make([]*Tile, self.width)
		self.heatmap[y] = make([]float64, self.width)
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
				value = fmt.Sprintf(`{"x":%d, "y":%d, "wall0":"%s", "wall1":"%s", "floor":"%s","rotation":%d, "decoration": "%s"}`,
					x,
					y,
					t.wall[0],
					t.wall[1],
					t.floor,
					t.rotation,
					decorationname,
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
	self.ComputeHeatMap()
}

func (self *DatacenterMap) ComputeHeatMap() {
	log.Debug("DatacenterMap::ComputeHeatMap()")
	// to raise 1 m-3 of air by 1 deg C : 1.211kj
	// 1kwh = 3600kJ /h?
	// 1BTU ~ 1kJ /s? or kwh?
	//
	// so
	self.heatmap = make([][]float64, self.height)
	nbairflow := float64(0)
	for y := 0; y < int(self.height); y++ {
		self.heatmap[y] = make([]float64, self.width)
		for x := 0; x < int(self.width); x++ {
			self.heatmap[y][x] = self.externaltemp
			if self.tiles[y][x].floor == "inside.air" {
				nbairflow++
			}
		}
	}

	// removePerAirFlow is in wH
	// GetHotspotValue() is in wH
	_, _, removePerAirflow := self.inventory.GetGlobalPower()
	if nbairflow != 0 {
		removePerAirflow = removePerAirflow / nbairflow
	}

	// we loop 600 times to spread the heat
	for loop := 0; loop < 600; loop++ {
		nextheatmap := make([][]float64, self.height)
		for y := 0; y < int(self.height); y++ {
			nextheatmap[y] = make([]float64, self.width)
			for x := 0; x < int(self.width); x++ {
				nextheatmap[y][x] = self.heatmap[y][x]
				// we flow air only inside
				if self.tiles[y][x].floor != "green" {
					// we add the heat
					nextheatmap[y][x] += self.inventory.GetHotspotValue(int32(y), int32(x)) / 1000
					// we remove via the airflow
					if self.tiles[y][x].floor == "inside.air" {
						nextheatmap[y][x] -= removePerAirflow / 1000
						if nextheatmap[y][x] < 17 {
							nextheatmap[y][x] = 17
						}
					}
					var left, right, top, bottom float64
					var leftfactor, rightfactor, topfactor, bottomfactor float64
					if y > 0 {
						if self.tiles[y-1][x].floor == "green" {
							// building wall are isolating a bit
							bottomfactor = 0.01
						} else {
							bottomfactor = 0.2
						}
						bottom = self.heatmap[y-1][x]
					} else {
						bottomfactor = 0.01
						bottom = self.externaltemp
					}
					if x > 0 {
						if self.tiles[y][x-1].floor == "green" {
							// building wall are isolating a bit
							leftfactor = 0.01
						} else {
							leftfactor = 0.2
						}
						left = self.heatmap[y][x-1]
					} else {
						leftfactor = 0.01
						left = self.externaltemp
					}
					if y < int(self.height)-1 {
						if self.tiles[y+1][x].floor == "green" {
							// building wall are isolating a bit
							topfactor = 0.01
						} else {
							topfactor = 0.2
						}
						top = self.heatmap[y+1][x]
					} else {
						topfactor = 0.01
						top = self.externaltemp
					}
					if x < int(self.width)-1 {
						if self.tiles[y][x+1].floor == "green" {
							// building wall are isolating a bit
							rightfactor = 0.01
						} else {
							rightfactor = 0.2
						}
						right = self.heatmap[y][x+1]
					} else {
						rightfactor = 0.01
						right = self.externaltemp
					}
					nextheatmap[y][x] = (1-leftfactor-rightfactor-topfactor-bottomfactor)*nextheatmap[y][x] +
						leftfactor*left + rightfactor*right + topfactor*top + bottomfactor*bottom
				}
			}
		}
		self.heatmap = nextheatmap
	}
}

func NewDatacenterMap() *DatacenterMap {
	dcmap := &DatacenterMap{
		tiles:        [][]*Tile{{}},
		heatmap:      [][]float64{},
		inventory:    nil,
		width:        0,
		height:       0,
		externaltemp: 0,
	}
	return dcmap
}
