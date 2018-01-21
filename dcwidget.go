package dctycoon

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/sws"
	log "github.com/sirupsen/logrus"
	"github.com/veandco/go-sdl2/sdl"
)

//
// This widget allow to display a Datacenter map (and more)
//
type DcWidget struct {
	sws.CoreWidget
	rackwidget   *RackWidget
	rootwindow   *sws.RootWidget
	tiles        [][]*Tile
	xRoot, yRoot int32 // offset of the whole map
	activeTile   *Tile
	te           *sws.TimerEvent
	inventory    *supplier.Inventory
	activeX      int32
	activeY      int32
	// the upper left list of placeable hardware (AC,rack,generator...)
	hc *HardwareChoice
}

func (self *DcWidget) DragDrop(x, y int32, payload sws.DragPayload) bool {
	// rack server (-> to install into a rack tower)
	if payload.GetType() == global.DRAG_RACK_SERVER {
		item := payload.(*ServerDragPayload).item
		tile, tx, ty, _ := self.findTile(x, y)
		if tile == nil || tile.element == nil || tile.element.ElementType() != supplier.PRODUCT_RACK {
			return false
		}
		nbu := item.Serverconf.ConfType.NbU
		var busy [42]bool

		// first create the map of what is filled
		for _, i := range self.inventory.Items {
			if i.Xplaced == tx && i.Yplaced == ty && i.Typeitem == supplier.PRODUCT_SERVER {
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
				self.inventory.InstallItem(item, tx, ty, int32(i))
				return true
			}
		}
		return false
	}

	// other: tower, ac, generator, rack
	if payload.GetType() == global.DRAG_ELEMENT_PAYLOAD {
		item := payload.(*ElementDragPayload).item
		height := payload.(*ElementDragPayload).imageheight
		tile, tx, ty := self.findFloorTile(x, y+height/2-24)
		if tile == nil || tile.element != nil {
			return false
		}

		self.inventory.InstallItem(item, tx, ty, -1)
		return true
	}
	return false
}

func (self *DcWidget) Repaint() {
	mapheight := len(self.tiles)
	mapwidth := len(self.tiles[0])
	self.FillRect(0, 0, self.Width(), self.Height(), 0xff000000)
	for y := 0; y < mapheight; y++ {
		for x := 0; x < mapwidth; x++ {
			tile := self.tiles[y][x]
			if tile != nil {
				surface := (*tile).Draw()
				rectSrc := sdl.Rect{0, 0, surface.W, surface.H}
				rectDst := sdl.Rect{self.xRoot + (self.Surface().W / 2) + (TILE_WIDTH_STEP/2)*int32(x) - (TILE_WIDTH_STEP/2)*int32(y), self.yRoot + (TILE_HEIGHT_STEP/2)*int32(x) + (TILE_HEIGHT_STEP/2)*int32(y), surface.W, surface.H}
				surface.Blit(&rectSrc, self.Surface(), &rectDst)
			}
		}
	}
	for _, child := range self.GetChildren() {
		// adjust the clipping to the current child
		child.Repaint()
		rectSrc := sdl.Rect{0, 0, child.Width(), child.Height()}
		rectDst := sdl.Rect{child.X(), child.Y(), child.Width(), child.Height()}
		child.Surface().Blit(&rectSrc, self.Surface(), &rectDst)
	}
}

//
// helper function, to know which pixel is in (x.y)
//
// It is mainly used to know if we are on a transparent pixel
//
func GetSurfacePixel(surface *sdl.Surface, x, y int32) (red, green, blue, alpha uint8) {
	if x < 0 || x >= surface.W || y < 0 || y >= surface.H {
		return 0, 0, 0, 0
	}
	err := surface.Lock()
	if err != nil {
		panic(err)
	}
	bpp := surface.Format.BytesPerPixel
	bytes := surface.Pixels()
	red = bytes[int(y)*int(surface.Pitch)+int(x)*int(bpp)]
	green = bytes[int(y)*int(surface.Pitch)+int(x)*int(bpp)+1]
	blue = bytes[int(y)*int(surface.Pitch)+int(x)*int(bpp)+2]
	alpha = bytes[int(y)*int(surface.Pitch)+int(x)*int(bpp)+3]

	surface.Unlock()
	return
}

//
// this method will return the tile where the cusor point to
// ie if the cursor point to the floor or an element on the tile
//
func (self *DcWidget) findTile(x, y int32) (*Tile, int32, int32, bool) {
	mapheight := len(self.tiles)
	mapwidth := len(self.tiles[0])
	for ty := mapheight - 1; ty >= 0; ty-- {
		for tx := mapwidth - 1; tx >= 0; tx-- {
			tile := self.tiles[ty][tx]
			surface := (*tile).Draw()
			xShift := self.xRoot + (self.Surface().W / 2) + (TILE_WIDTH_STEP/2)*int32(tx) - (TILE_WIDTH_STEP/2)*int32(ty)
			yShift := self.yRoot + (TILE_HEIGHT_STEP/2)*int32(tx) + (TILE_HEIGHT_STEP/2)*int32(ty)

			if (x >= xShift) &&
				(y >= yShift) &&
				(x < xShift+surface.W) &&
				(y < yShift+surface.H) {
				_, _, _, alpha := GetSurfacePixel(surface, x-xShift, y-yShift)
				if alpha > 0 {
					return tile, int32(tx), int32(ty), tile.IsElementAt(x-xShift, y-yShift)
				}
			}
		}
	}
	return nil, -1, -1, false
}

func (self *DcWidget) findFloorTile(x, y int32) (*Tile, int32, int32) {
	mapheight := len(self.tiles)
	mapwidth := len(self.tiles[0])
	// compute the x-y where the mouse is
	tilex := ((x-self.xRoot-(self.Surface().W/2)-TILE_WIDTH_STEP/2-10)/2 + y - self.yRoot - TILE_HEIGHT + TILE_HEIGHT_STEP + 8) / TILE_HEIGHT_STEP
	tiley := (y - self.yRoot - TILE_HEIGHT + TILE_HEIGHT_STEP + 8 - (x-self.xRoot-(self.Surface().W/2)-TILE_WIDTH_STEP/2-10)/2) / TILE_HEIGHT_STEP

	//fmt.Println("DcWidget::MouveMove",tilex,tiley)
	if tilex < 0 {
		return nil, -1, -1
	}
	if tiley < 0 {
		return nil, -1, -1
	}
	if tilex >= int32(mapwidth) {
		return nil, -1, -1
	}
	if tiley >= int32(mapheight) {
		return nil, -1, -1
	}

	return self.tiles[tiley][tilex], tilex, tiley
}

func (self *DcWidget) MousePressDown(x, y int32, button uint8) {
	self.rackwidget.activeElement = nil
	activeTile, xtile, ytile, isElement := self.findTile(x, y)
	self.activeTile = activeTile

	// now, do we have an active item inside the tile
	if activeTile != nil && isElement {
		//fmt.Println("active element!!!")
		self.activeX = xtile
		self.activeY = ytile

	}
}

func (self *DcWidget) MouseDoubleClick(x, y int32) {
	activeTile, xtile, ytile, isElement := self.findTile(x, y)
	if activeTile != nil && isElement {
		activeElement := activeTile.TileElement()
		if activeElement != nil && activeElement.ElementType() == supplier.PRODUCT_RACK {
			self.rackwidget.activeElement = activeTile.TileElement()
			self.rackwidget.xactiveElement = xtile
			self.rackwidget.yactiveElement = ytile
			self.rackwidget.Show(self.rackwidget.xactiveElement, self.rackwidget.yactiveElement)
		}
	}
}

func (self *DcWidget) MousePressUp(x, y int32, button uint8) {
	if button == sdl.BUTTON_LEFT {
	}
	if button == sdl.BUTTON_RIGHT && self.activeTile != nil {
		// if we are on a rack
		activeElement := self.activeTile.TileElement()
		if activeElement != nil && activeElement.ElementType() == supplier.PRODUCT_RACK {
			m := sws.NewMenuWidget()
			activeTile := self.activeTile
			m.AddItem(sws.NewMenuItemLabel("Rotate", func() {
				activeTile.Rotate((activeTile.rotation + 1) % 4)
			}))
			// prepare rackwidget
			self.rackwidget.activeElement = self.activeTile.TileElement()
			self.rackwidget.xactiveElement = self.activeX
			self.rackwidget.yactiveElement = self.activeY
			m.AddItem(sws.NewMenuItemLabel("Details", func() {
				self.rackwidget.Show(self.rackwidget.xactiveElement, self.rackwidget.yactiveElement)
			}))
			m.AddItem(sws.NewMenuItemLabel("Uninstall", func() {
				rackelement := activeTile.TileElement().(*RackElement)
				if len(rackelement.items) > 0 {
					iconsurface, _ := global.LoadImageAsset("assets/ui/icon-triangular-big.png")
					sws.ShowModalErrorSurfaceicon(self.rootwindow, "Uninstall action", iconsurface, "It is not possible to uninstall a rack unless it is empty", nil)
				} else {
					self.inventory.UninstallItem(rackelement.InventoryItem())
				}
			}))
			m.Move(x, y)
			sws.ShowMenu(m)
		} else if activeElement != nil && (activeElement.ElementType() == supplier.PRODUCT_GENERATOR ||
			activeElement.ElementType() == supplier.PRODUCT_AC ||
			activeElement.ElementType() == supplier.PRODUCT_SERVER) {
			m := sws.NewMenuWidget()
			activeTile := self.activeTile
			m.AddItem(sws.NewMenuItemLabel("Rotate", func() {
				activeTile.Rotate((activeTile.rotation + 1) % 4)
			}))
			m.AddItem(sws.NewMenuItemLabel("Uninstall", func() {
				self.inventory.UninstallItem(activeTile.TileElement().InventoryItem())
			}))
			m.Move(x, y)
			sws.ShowMenu(m)
		}
	}
	self.activeTile = nil
}

func (self *DcWidget) HasFocus(focus bool) {
	if focus == false {
		self.te.StopRepeat()
		self.te = nil
	}
}

func (self *DcWidget) MouseMove(x, y, xrel, yrel int32) {
	if self.te != nil {
		self.te.StopRepeat()
		self.te = nil
	}
	if x < 10 {
		self.te = sws.TimerAddEvent(time.Now(), 50*time.Millisecond, func(evt *sws.TimerEvent) {
			self.MoveLeft()
		})
		return
	}
	if y < 10 {
		self.te = sws.TimerAddEvent(time.Now(), 50*time.Millisecond, func(evt *sws.TimerEvent) {
			self.MoveUp()
		})
		return
	}
	if x > self.Width()-10 {
		self.te = sws.TimerAddEvent(time.Now(), 50*time.Millisecond, func(evt *sws.TimerEvent) {
			self.MoveRight()
		})
		return
	}
	if y > self.Height()-10 {
		self.te = sws.TimerAddEvent(time.Now(), 50*time.Millisecond, func(evt *sws.TimerEvent) {
			self.MoveDown()
		})
		return
	}

	// if we are moving an element's tile
	if self.activeTile != nil && self.activeTile.TileElement() != nil {

		// compute the x-y where the mouse is
		mapheight := len(self.tiles)
		mapwidth := len(self.tiles[0])
		tilex := ((x-self.xRoot-(self.Surface().W/2)-TILE_WIDTH_STEP/2-10)/2 + y - self.yRoot - TILE_HEIGHT + TILE_HEIGHT_STEP + 8) / TILE_HEIGHT_STEP
		tiley := (y - self.yRoot - TILE_HEIGHT + TILE_HEIGHT_STEP + 8 - (x-self.xRoot-(self.Surface().W/2)-TILE_WIDTH_STEP/2-10)/2) / TILE_HEIGHT_STEP

		//fmt.Println("DcWidget::MouveMove",tilex,tiley)
		if tilex < 0 {
			tilex = 0
		}
		if tiley < 0 {
			tiley = 0
		}
		if tilex >= int32(mapwidth) {
			tilex = int32(mapwidth) - 1
		}
		if tiley >= int32(mapheight) {
			tiley = int32(mapheight) - 1
		}

		if (tilex != self.activeX || tiley != self.activeY) && self.tiles[tiley][tilex].element == nil {
			element := self.activeTile.element
			rotation := self.activeTile.rotation
			self.activeTile.element = nil
			self.activeTile.surface = nil

			// update the tiles
			self.activeX = tilex
			self.activeY = tiley
			self.activeTile = self.tiles[tiley][tilex]
			self.activeTile.element = element
			self.activeTile.rotation = rotation
			self.activeTile.surface = nil

			// update the InventoryItem
			if element != nil {
				element.InventoryItem().Xplaced = tilex
				element.InventoryItem().Yplaced = tiley
			}
			if element.ElementType() == supplier.PRODUCT_RACK {
				rack := element.(*RackElement)
				for _, i := range rack.items {
					i.Xplaced = tilex
					i.Yplaced = tiley
				}
			}

			self.PostUpdate()
		}
	}
}

func (self *DcWidget) MoveLeft() {
	width := int32(len(self.tiles[0])+len(self.tiles)+1) * TILE_WIDTH_STEP / 2
	if self.xRoot-width/2+self.Width()/2 < 0 {
		self.xRoot += 20
		self.PostUpdate()
	}
}

func (self *DcWidget) MoveUp() {
	height := int32(len(self.tiles[0])+len(self.tiles))*TILE_HEIGHT_STEP/2 + TILE_HEIGHT
	if self.yRoot-height/2+self.Height()/2 < 0 {
		self.yRoot += 20
		self.PostUpdate()
	}
}

func (self *DcWidget) MoveRight() {
	width := int32(len(self.tiles[0])+len(self.tiles)+1) * TILE_WIDTH_STEP / 2
	if self.xRoot+width > self.Width() {
		self.xRoot -= 20
		self.PostUpdate()
	}
}

func (self *DcWidget) MoveDown() {
	height := int32(len(self.tiles[0])+len(self.tiles))*TILE_HEIGHT_STEP/2 + TILE_HEIGHT
	if self.yRoot+height > self.Height() {
		self.yRoot -= 20
		self.PostUpdate()
	}
}

func (self *DcWidget) KeyDown(key sdl.Keycode, mod uint16) {
	if key == sdl.K_LEFT {
		self.MoveLeft()
	}
	if key == sdl.K_UP {
		self.MoveUp()
	}
	if key == sdl.K_RIGHT {
		self.MoveRight()
	}
	if key == sdl.K_DOWN {
		self.MoveDown()
	}

}

func (self *DcWidget) ItemInTransit(*supplier.InventoryItem) {
}

func (self *DcWidget) ItemInStock(*supplier.InventoryItem) {
}

func (self *DcWidget) ItemRemoveFromStock(*supplier.InventoryItem) {
}

func (self *DcWidget) ItemChangedPool(*supplier.InventoryItem) {
}

func (self *DcWidget) ItemInstalled(item *supplier.InventoryItem) {
	mapheight := len(self.tiles)
	mapwidth := len(self.tiles[0])
	if item.Xplaced <= int32(mapwidth) && item.Yplaced <= int32(mapheight) && item.Xplaced >= 0 && item.Yplaced >= 0 {
		//if item.Typeitem==supplier.PRODUCT_SERVER {
		self.tiles[item.Yplaced][item.Xplaced].ItemInstalled(item)
		//}
	}
}

func (self *DcWidget) ItemUninstalled(item *supplier.InventoryItem) {
	mapheight := len(self.tiles)
	mapwidth := len(self.tiles[0])
	if item.Xplaced <= int32(mapwidth) && item.Yplaced <= int32(mapheight) && item.Xplaced >= 0 && item.Yplaced >= 0 {
		//if item.Typeitem==supplier.PRODUCT_SERVER {
		self.tiles[item.Yplaced][item.Xplaced].ItemUninstalled(item)
		//}
	}
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
func (self *DcWidget) LoadMap(dc map[string]interface{}) {
	log.Debug("DcWidget::LoadMap(", dc, ")")
	width := int32(dc["width"].(float64))
	height := int32(dc["height"].(float64))
	self.tiles = make([][]*Tile, height)
	for y := range self.tiles {
		self.tiles[y] = make([]*Tile, width)
		for x := range self.tiles[y] {
			self.tiles[y][x] = NewGrassTile()
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
}

func (self *DcWidget) InitMap(assetdcmap string) {
	log.Debug("DcWidget::InitMap(", assetdcmap, ")")
	if data, err := global.Asset("assets/dcmap/" + assetdcmap); err == nil {
		var dcmap map[string]interface{}
		if json.Unmarshal(data, &dcmap) == nil {
			width := int32(dcmap["width"].(float64))
			height := int32(dcmap["height"].(float64))
			self.tiles = make([][]*Tile, height)
			for y := range self.tiles {
				self.tiles[y] = make([]*Tile, width)
				for x := range self.tiles[y] {
					self.tiles[y][x] = NewGrassTile()
				}
			}
			tiles := dcmap["tiles"].([]interface{})
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
		}
	}
}

func (self *DcWidget) SaveMap() string {
	s := fmt.Sprintf(`{"width":%d, "height":%d, "tiles": [`, len(self.tiles[0]), len(self.tiles))
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

func NewDcWidget(w, h int32, rootwindow *sws.RootWidget) *DcWidget {
	corewidget := sws.NewCoreWidget(w, h)
	rackwidget := NewRackWidget(rootwindow)
	widget := &DcWidget{CoreWidget: *corewidget,
		rackwidget: rackwidget,
		rootwindow: rootwindow,
		tiles:      [][]*Tile{{}},
		xRoot:      0,
		yRoot:      0,
		inventory:  nil,
		hc:         NewHardwareChoice(),
	}

	//widget.hc.Move(0,h/2-100)
	widget.hc.Move(0, 0)
	widget.AddChild(widget.hc)
	return widget
}

func (self *DcWidget) SetGame(inventory *supplier.Inventory, currenttime time.Time) {
	log.Debug("DcWidget::SetGame()")
	if self.inventory != nil {
		self.inventory.RemoveInventorySubscriber(self)
	}
	self.inventory = inventory
	inventory.AddInventorySubscriber(self)

	self.rackwidget.SetGame(inventory, currenttime)
	self.hc.SetGame(inventory, currenttime)
}
