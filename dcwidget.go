package dctycoon

import (
	"strconv"
	"time"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/sws"
	log "github.com/sirupsen/logrus"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	SHOW_MAPWALL = iota
	SHOW_MAP     = iota
	SHOW_HEATMAP = iota
)

//
// This widget allow to display a Datacenter map (and more)
//
type DcWidget struct {
	sws.CoreWidget
	rackwidget    *RackWidget
	rootwindow    *sws.RootWidget
	dcmap         *DatacenterMap
	xRoot, yRoot  int32 // offset of the whole map
	activeTile    *Tile // if we click, to know which element tile we are relocating
	te            *sws.TimerEvent
	activeX       int32
	activeY       int32
	hc            *HardwareChoice // the upper left list of placeable hardware (AC,rack,generator...)
	showMap       int32
	blinkon       bool
	heatmapButton *sws.FlatButtonWidget

	showInventoryManagement func()
}

func (self *DcWidget) DragDrop(x, y int32, payload sws.DragPayload) bool {
	// rack server (-> to install into a rack tower)
	if payload.GetType() == global.DRAG_RACK_SERVER {
		item := payload.(*ServerDragPayload).item
		tile, tx, ty, _ := self.findTile(x, y)
		if tile == nil || tile.element == nil || tile.element.ElementType() != supplier.PRODUCT_RACK {
			return false
		}
		return self.dcmap.InstallItem(item, tx, ty)
	}

	// other: tower, ac, generator, rack
	if payload.GetType() == global.DRAG_ELEMENT_PAYLOAD {
		item := payload.(*ElementDragPayload).item
		height := payload.(*ElementDragPayload).imageheight
		tile, tx, ty := self.findFloorTile(x, y+height/2-24)
		if tile == nil || tile.element != nil {
			return false
		}
		if (!tile.IsFloorOutside() && item.Typeitem == supplier.PRODUCT_GENERATOR) ||
			(!tile.IsFloorInsideNotAirFlow() && item.Typeitem != supplier.PRODUCT_GENERATOR) {
			return false
		}

		return self.dcmap.InstallItem(item, tx, ty)
	}
	return false
}

func (self *DcWidget) Repaint() {
	self.FillRect(0, 0, self.Width(), self.Height(), 0xff000000)

	// prepare heat map tiles
	var heat [10]*sdl.Surface
	for i := 0; i < 10; i++ {
		heat[i], _ = sdl.CreateRGBSurface(0, 105, TILE_HEIGHT, 32, 0x00ff0000, 0x0000ff00, 0x000000ff, 0xff000000)
		floor := getSprite("assets/ui/heat" + strconv.Itoa(i) + ".png")
		rectSrc := sdl.Rect{0, 0, floor.W, floor.H}
		rectDst := sdl.Rect{0, TILE_HEIGHT - floor.H, floor.W, floor.H}
		floor.Blit(&rectSrc, heat[i], &rectDst)
	}
	for y := int32(0); y < self.dcmap.GetHeight(); y++ {
		for x := int32(0); x < self.dcmap.GetWidth(); x++ {
			tile := self.dcmap.GetTile(x, y)
			var surface *sdl.Surface
			if self.showMap != SHOW_HEATMAP {
				if self.showMap == SHOW_MAPWALL {
					surface = (*tile).Draw()
				} else {
					surface = (*tile).DrawWithoutWall()
				}
			} else {
				if tile.floor == "green" {
					surface = (*tile).Draw()
				} else {
					temp := int(self.dcmap.GetTemperature(x, y)-17) / 2
					if temp < 0 {
						temp = 0
					}
					if temp > 9 {
						temp = 9
					}
					surface = heat[temp]
				}
			}
			rectSrc := sdl.Rect{0, 0, surface.W, surface.H}
			rectDst := sdl.Rect{self.xRoot + (self.Surface().W / 2) + (TILE_WIDTH_STEP/2)*int32(x) - (TILE_WIDTH_STEP/2)*int32(y), self.yRoot + (TILE_HEIGHT_STEP/2)*int32(x) + (TILE_HEIGHT_STEP/2)*int32(y), surface.W, surface.H}
			surface.Blit(&rectSrc, self.Surface(), &rectDst)

			var extraIcon *sdl.Surface
			if self.dcmap.GetGeneralOutage() {
				element := tile.TileElement()
				if element != nil && (element.ElementType() == supplier.PRODUCT_AC ||
					element.ElementType() == supplier.PRODUCT_RACK ||
					element.ElementType() == supplier.PRODUCT_SERVER) {
					extraIcon, _ = global.LoadImageAsset("assets/ui/lightning-bolt-shadow.png")
				}
			}
			switch self.dcmap.GetRackStatus(x, y) {
			case RACK_OVER_CURRENT:
				extraIcon, _ = global.LoadImageAsset("assets/ui/lightning-bolt-shadow.png")
			case RACK_HEAT_WARNING:
				extraIcon, _ = global.LoadImageAsset("assets/ui/thermometer.warning.png")
			case RACK_OVER_HEAT:
				extraIcon, _ = global.LoadImageAsset("assets/ui/thermometer.png")
			case RACK_MELTING:
				extraIcon, _ = global.LoadImageAsset("assets/ui/thermometer.png")
			}
			if extraIcon != nil && self.blinkon {
				rectSrc = sdl.Rect{0, 0, extraIcon.W, extraIcon.H}
				rectDst = sdl.Rect{self.xRoot + (self.Surface().W / 2) + (TILE_WIDTH_STEP/2)*int32(x) - (TILE_WIDTH_STEP/2)*int32(y) + (surface.W-extraIcon.W)/2, self.yRoot + (TILE_HEIGHT_STEP/2)*int32(x) + (TILE_HEIGHT_STEP/2)*int32(y) + 20, extraIcon.W, extraIcon.H}
				extraIcon.Blit(&rectSrc, self.Surface(), &rectDst)
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
	self.SetDirtyFalse()
}

//
// this method will return the tile where the cusor point to
// ie if the cursor point to the floor or an element on the tile
//
func (self *DcWidget) findTile(x, y int32) (*Tile, int32, int32, bool) {
	for ty := self.dcmap.GetHeight() - 1; ty >= 0; ty-- {
		for tx := self.dcmap.GetWidth() - 1; tx >= 0; tx-- {
			tile := self.dcmap.GetTile(tx, ty)
			surface := (*tile).Draw()
			xShift := self.xRoot + (self.Surface().W / 2) + (TILE_WIDTH_STEP/2)*int32(tx) - (TILE_WIDTH_STEP/2)*int32(ty)
			yShift := self.yRoot + (TILE_HEIGHT_STEP/2)*int32(tx) + (TILE_HEIGHT_STEP/2)*int32(ty)

			if (x >= xShift) &&
				(y >= yShift) &&
				(x < xShift+surface.W) &&
				(y < yShift+surface.H) {
				_, _, _, alpha := global.GetSurfacePixel(surface, x-xShift, y-yShift)
				if alpha > 0 {
					return tile, int32(tx), int32(ty), tile.TileElement() != nil
				}
			}
		}
	}
	return nil, -1, -1, false
}

func (self *DcWidget) findFloorTile(x, y int32) (*Tile, int32, int32) {
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
	if tilex >= self.dcmap.GetWidth() {
		return nil, -1, -1
	}
	if tiley >= self.dcmap.GetHeight() {
		return nil, -1, -1
	}

	return self.dcmap.GetTile(tilex, tiley), tilex, tiley
}

func (self *DcWidget) MousePressDown(x, y int32, button uint8) {
	if self.showMap == SHOW_HEATMAP {
		return
	}
	self.rackwidget.activeElement = nil
	activeTile, xtile, ytile, isElement := self.findTile(x, y)

	self.activeTile = activeTile
	// now, do we have an active item inside the tile
	if activeTile != nil && isElement && activeTile.TileElement().ElementType() != supplier.PRODUCT_DECORATION {
		//fmt.Println("active element!!!")
		self.activeX = xtile
		self.activeY = ytile

	}
}

func (self *DcWidget) MouseDoubleClick(x, y int32) {
	if self.showMap == SHOW_HEATMAP {
		return
	}
	activeTile, xtile, ytile, isElement := self.findTile(x, y)
	if activeTile != nil && isElement {
		activeElement := activeTile.TileElement()
		if activeElement != nil && activeElement.ElementType() == supplier.PRODUCT_RACK {
			self.rackwidget.activeElement = activeTile.TileElement()
			self.rackwidget.xactiveElement = xtile
			self.rackwidget.yactiveElement = ytile
			self.rackwidget.Show(self.rackwidget.xactiveElement, self.rackwidget.yactiveElement)
		}
		if activeElement != nil && activeElement.ElementType() == supplier.PRODUCT_DECORATION {
			decoration := activeElement.(*DecorationElement)
			if decoration.GetName() == "shelf" && self.showInventoryManagement != nil {
				self.showInventoryManagement()
			}
		}
	}
}

func (self *DcWidget) MousePressUp(x, y int32, button uint8) {
	if self.showMap == SHOW_HEATMAP {
		return
	}
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
				self.PostUpdate()
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
					activeTile.Rotate((0)
					self.dcmap.UninstallItem(rackelement.InventoryItem())
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
				self.PostUpdate()
			}))
			m.AddItem(sws.NewMenuItemLabel("Uninstall", func() {
				activeTile.Rotate((0)
				self.dcmap.UninstallItem(activeTile.TileElement().InventoryItem())
			}))
			m.Move(x, y)
			sws.ShowMenu(m)
		} else if activeElement != nil && (activeElement.ElementType() == supplier.PRODUCT_DECORATION) {
			decoration := activeElement.(*DecorationElement)
			if decoration.GetName() == "shelf" && self.showInventoryManagement != nil {
				m := sws.NewMenuWidget()
				m.AddItem(sws.NewMenuItemLabel("Details", func() {
					self.showInventoryManagement()
				}))
				m.Move(x, y)
				sws.ShowMenu(m)
			}
		} else if activeElement == nil {
			if self.activeTile.IsFloorInsideNotAirFlow() {
				m := sws.NewMenuWidget()
				tile := self.activeTile
				m.AddItem(sws.NewMenuItemLabel("Switch to air flow", func() {
					tile.SwitchToAirFlow()
					self.dcmap.ComputeHeatMap()
					self.dcmap.ComputeOverLimits()
					self.PostUpdate()
				}))
				m.Move(x, y)
				sws.ShowMenu(m)
			}
			if self.activeTile.IsFloorInsideAirFlow() {
				m := sws.NewMenuWidget()
				tile := self.activeTile
				m.AddItem(sws.NewMenuItemLabel("Remove air flow", func() {
					tile.SwitchToNotAirFlow()
					self.dcmap.ComputeHeatMap()
					self.dcmap.ComputeOverLimits()
					self.PostUpdate()
				}))
				m.Move(x, y)
				sws.ShowMenu(m)
			}
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
	if self.activeTile != nil && self.activeTile.TileElement() != nil && self.activeTile.TileElement().ElementType() != supplier.PRODUCT_DECORATION {

		// compute the x-y where the mouse is
		tilex := ((x-self.xRoot-(self.Surface().W/2)-TILE_WIDTH_STEP/2-10)/2 + y - self.yRoot - TILE_HEIGHT + TILE_HEIGHT_STEP + 8) / TILE_HEIGHT_STEP
		tiley := (y - self.yRoot - TILE_HEIGHT + TILE_HEIGHT_STEP + 8 - (x-self.xRoot-(self.Surface().W/2)-TILE_WIDTH_STEP/2-10)/2) / TILE_HEIGHT_STEP

		if tilex < 0 {
			tilex = 0
		}
		if tiley < 0 {
			tiley = 0
		}
		if tilex >= self.dcmap.GetWidth() {
			tilex = self.dcmap.GetWidth() - 1
		}
		if tiley >= self.dcmap.GetHeight() {
			tiley = self.dcmap.GetHeight() - 1
		}

		// move the element
		if self.dcmap.MoveElement(self.activeX, self.activeY, tilex, tiley) == true {
			self.activeX = tilex
			self.activeY = tiley
			self.activeTile = self.dcmap.GetTile(tilex, tiley)
			self.PostUpdate()
		}
	}
}

func (self *DcWidget) MoveLeft() {
	width := (self.dcmap.GetWidth() + self.dcmap.GetHeight() + 1) * TILE_WIDTH_STEP / 2
	if self.xRoot-width/2+self.Width()/2 < 0 {
		self.xRoot += 20
		self.PostUpdate()
	}
}

func (self *DcWidget) MoveUp() {
	height := (self.dcmap.GetWidth()+self.dcmap.GetHeight())*TILE_HEIGHT_STEP/2 + TILE_HEIGHT
	if self.yRoot-height/2+self.Height()/2 < 0 {
		self.yRoot += 20
		self.PostUpdate()
	}
}

func (self *DcWidget) MoveRight() {
	width := (self.dcmap.GetWidth() + self.dcmap.GetHeight() + 1) * TILE_WIDTH_STEP / 2
	if self.xRoot+width > self.Width() {
		self.xRoot -= 20
		self.PostUpdate()
	}
}

func (self *DcWidget) MoveDown() {
	height := (self.dcmap.GetWidth()+self.dcmap.GetHeight())*TILE_HEIGHT_STEP/2 + TILE_HEIGHT
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

func (self *DcWidget) RackStatusChange(x, y int32, rackstate int32) {
	if rackstate == RACK_MELTING {
		flash := uint32(9)
		tilex := x
		tiley := y
		self.dcmap.GetTile(tilex, tiley).SetFlashEffect(flash)
		sws.TimerAddEvent(time.Now(), 75*time.Millisecond, func(evt *sws.TimerEvent) {
			self.dcmap.GetTile(tilex, tiley).SetFlashEffect(flash)
			if flash == 0 {
				evt.StopRepeat()
			} else {
				flash = flash - 1
			}
			self.PostUpdate()
		})
	}
}

func (self *DcWidget) GeneralOutage(bool) {}

func (self *DcWidget) SetGame(inventory *supplier.Inventory, currenttime time.Time, dcmap *DatacenterMap) {
	log.Debug("DcWidget::SetGame(", inventory, ",", currenttime, ",", dcmap, ")")
	self.dcmap = dcmap
	dcmap.AddRackStatusSubscriber(self)
	self.rackwidget.SetGame(inventory, currenttime)
	self.hc.SetGame(inventory, currenttime)
	self.showMap = SHOW_MAPWALL
	self.heatmapButton.SetText("Map without wall")
	self.PostUpdate()
}

func NewDcWidget(w, h int32, rootwindow *sws.RootWidget) *DcWidget {
	corewidget := sws.NewCoreWidget(w, h)
	rackwidget := NewRackWidget(rootwindow)
	widget := &DcWidget{CoreWidget: *corewidget,
		rackwidget:    rackwidget,
		rootwindow:    rootwindow,
		dcmap:         nil,
		xRoot:         0,
		yRoot:         0,
		hc:            NewHardwareChoice(),
		showMap:       SHOW_MAPWALL,
		blinkon:       false,
		heatmapButton: sws.NewFlatButtonWidget(200, 40, "Map without wall"),
	}

	//widget.hc.Move(0,h/2-100)
	widget.hc.Move(0, 0)
	widget.AddChild(widget.hc)

	if icon, err := global.LoadImageAsset("assets/ui/map.png"); err == nil {
		widget.heatmapButton.SetImageSurface(icon)
	}
	widget.heatmapButton.SetClicked(func() {
		var text string
		switch widget.showMap {
		case SHOW_MAPWALL:
			widget.showMap = SHOW_MAP
			text = "HeatMap"
		case SHOW_MAP:
			widget.showMap = SHOW_HEATMAP
			text = "Map"
		case SHOW_HEATMAP:
			widget.showMap = SHOW_MAPWALL
			text = "Map without wall"
		}
		widget.heatmapButton.SetText(text)
		if widget.showMap == SHOW_HEATMAP {
			widget.blinkon = false
			widget.RemoveChild(widget.hc)
		} else {
			widget.AddChild(widget.hc)
		}
		widget.PostUpdate()
	})
	widget.heatmapButton.SetColor(0x80ffffff)
	//	widget.heatmapButton.SetTextColor(sdl.Color{0xff, 0xff, 0xff, 0xff})
	widget.heatmapButton.AlignImageLeft(true)
	widget.heatmapButton.SetCentered(false)
	widget.heatmapButton.Move(30, rootwindow.Height()-65)
	widget.AddChild(widget.heatmapButton)

	sws.TimerAddEvent(time.Now(), 1000*time.Millisecond, func(evt *sws.TimerEvent) {
		// to "blink" on over heat or over current
		if widget.showMap != SHOW_HEATMAP {
			widget.blinkon = !widget.blinkon
		}
		widget.PostUpdate()
	})

	return widget
}

func (self *DcWidget) SetInventoryManagementCallback(showInventoryManagement func()) {
	self.showInventoryManagement = showInventoryManagement
}
