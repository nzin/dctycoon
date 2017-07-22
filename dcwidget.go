package dctycoon

import (
	"fmt"
	"github.com/nzin/sws"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/nzin/dctycoon/supplier"
	"time"
)

//
// This widget allow to display a Datacenter map (and more)
//
type DcWidget struct {
	sws.CoreWidget
	rackwidget    *RackWidget
	tiles         [][]*Tile
	xRoot, yRoot  int32
	activeTile    *Tile
	te            *sws.TimerEvent
	inventory     *supplier.Inventory
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
	sws.PostUpdate()
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

func (self *DcWidget) MousePressDown(x, y int32, button uint8) {
	self.activeTile = nil
	self.rackwidget.activeElement = nil
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
					//fmt.Println("activeTile: "+tile.floor,tx,ty)
					self.activeTile = tile
					// now, do we have an active item inside the tile
					if tile.IsElementAt(x-xShift, y-yShift) {
						//fmt.Println("active element!!!")
						self.rackwidget.activeElement = tile.DcElement()
						self.rackwidget.xactiveElement = int32(tx)
						self.rackwidget.yactiveElement = int32(ty)
					}
					break
				}
			}
		}
		if self.activeTile != nil {
			break
		}
	}
}

func (self *DcWidget) MousePressUp(x, y int32, button uint8) {
	if self.rackwidget.activeElement != nil {
		self.rackwidget.Show(self.rackwidget.xactiveElement,self.rackwidget.yactiveElement)
	}
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
		self.te = sws.TimerAddEvent(time.Now(), 50*time.Millisecond, func() {
			self.MoveLeft()
		})
		return
	}
	if y < 10 {
		self.te = sws.TimerAddEvent(time.Now(), 50*time.Millisecond, func() {
			self.MoveUp()
		})
		return
	}
	if x > self.Width()-10 {
		self.te = sws.TimerAddEvent(time.Now(), 50*time.Millisecond, func() {
			self.MoveRight()
		})
		return
	}
	if y > self.Height()-10 {
		self.te = sws.TimerAddEvent(time.Now(), 50*time.Millisecond, func() {
			self.MoveDown()
		})
		return
	}
}

func (self *DcWidget) MoveLeft() {
	width := int32(len(self.tiles[0])+len(self.tiles)+1) * TILE_WIDTH_STEP / 2
	if self.xRoot-width/2+self.Width()/2 < 0 {
		self.xRoot += 20
		sws.PostUpdate()
	}
}

func (self *DcWidget) MoveUp() {
	height := int32(len(self.tiles[0])+len(self.tiles))*TILE_HEIGHT_STEP/2 + TILE_HEIGHT
	if self.yRoot-height/2+self.Height()/2 < 0 {
		self.yRoot += 20
		sws.PostUpdate()
	}
}

func (self *DcWidget) MoveRight() {
	width := int32(len(self.tiles[0])+len(self.tiles)+1) * TILE_WIDTH_STEP / 2
	if self.xRoot+width > self.Width() {
		self.xRoot -= 20
		sws.PostUpdate()
	}
}

func (self *DcWidget) MoveDown() {
	height := int32(len(self.tiles[0])+len(self.tiles))*TILE_HEIGHT_STEP/2 + TILE_HEIGHT
	if self.yRoot+height > self.Height() {
		self.yRoot -= 20
		sws.PostUpdate()
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

func (self *DcWidget) ItemInstalled(item *supplier.InventoryItem) {
	mapheight := len(self.tiles)
	mapwidth := len(self.tiles[0])
	if (item.Xplaced<=int32(mapwidth) && item.Yplaced<=int32(mapheight) && item.Xplaced>=0 && item.Yplaced>=0 ) {
		self.tiles[item.Yplaced][item.Xplaced].ItemInstalled(item)
	}
}

func (self *DcWidget) ItemUninstalled(item *supplier.InventoryItem) {
	mapheight := len(self.tiles)
	mapwidth := len(self.tiles[0])
	if (item.Xplaced<=int32(mapwidth) && item.Yplaced<=int32(mapheight) && item.Xplaced>=0 && item.Yplaced>=0 ) {
		self.tiles[item.Yplaced][item.Xplaced].ItemUninstalled(item)
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
		var dcelementtype string
		if tile["dcelementtype"] != nil {
			dcelementtype = tile["dcelementtype"].(string)
		}
		if dcelementtype == "" || dcelementtype == "rack" {
			// basic floor
			self.tiles[y][x] = NewElectricalTile(wall0, wall1, floor, rotation, dcelementtype)
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
			if t.element == nil {
				if t.wall[0] != "" || t.wall[1] != "" || t.floor != "green" {
					value = fmt.Sprintf(`{"x":%d, "y":%d, "wall0":"%s", "wall1":"%s", "floor":"%s","rotation":%d}`,
						x,
						y,
						t.wall[0],
						t.wall[1],
						t.floor,
						t.rotation,
					)
				}
			} else {
				value = fmt.Sprintf(`{"x":%d, "y":%d, "wall0":"%s", "wall1":"%s", "floor":"%s", "rotation":%d, "dcelementtype":"%s"}`,
					x,
					y,
					t.wall[0],
					t.wall[1],
					t.floor,
					t.rotation,
					t.element.ElementType(),
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

func NewDcWidget(w, h int32,rootwindow *sws.RootWidget,inventory *supplier.Inventory) *DcWidget {
	corewidget := sws.NewCoreWidget(w, h)
	rackwidget:=NewRackWidget(rootwindow,inventory)
	widget := &DcWidget{CoreWidget: *corewidget,
		rackwidget: rackwidget,
		tiles:      [][]*Tile{{}},
		xRoot:      0,
		yRoot:      0,
		inventory:  inventory,
	}
	inventory.AddSubscriber(widget)
	return widget
}
