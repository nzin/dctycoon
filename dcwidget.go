package dctycoon

import (
	"fmt"
	"github.com/nzin/sws"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/nzin/dctycoon/supplier"
	"time"
)

type ServerDragPayload struct {
	item *supplier.InventoryItem
}

func (self *ServerDragPayload) GetType() int32 {
	return 1
}

type RackWidgetLine struct {
	sws.SWS_Label
	item *supplier.InventoryItem
}

func (self *RackWidgetLine) MousePressDown(x, y int32, button uint8) {
	payload:=&ServerDragPayload{
		item: self.item,
	}
	var parent sws.SWS_Widget
	parent=self
	for (parent!=nil) {
		x+=parent.X()
		y+=parent.Y()
		parent=parent.Parent()
	}
	sws.NewDragEvent(x,y,"resources/"+self.item.Serverconf.ConfType.ServerSprite+"0.png", payload)
}

func NewRackWidgetLine(item *supplier.InventoryItem) *RackWidgetLine{
	label:=sws.CreateLabel(200,25,"server")
	return &RackWidgetLine{
		SWS_Label: *label,
		item: item,
	}
}

type RackChassisWidget struct {
	sws.SWS_CoreWidget
	inventory  *supplier.Inventory
	ydrag      int32
	items      []*supplier.InventoryItem
	comingitem *supplier.InventoryItem
}

func (self *RackChassisWidget) computeComingPos(ypos int32) int32 {
	ypos=(ypos-10)/10
	nbu:=self.comingitem.Serverconf.ConfType.NbU
	var busy [42]bool
	
	// first create the map of what is filled
	for _,item := range(self.items) {
		itemNbU:=item.Serverconf.ConfType.NbU
		for j:=0;j<int(itemNbU);j++ {
			if j+int(item.Zplaced)<42 {
				busy[j]=true
			}
		}
	}
	
	// now try to find a nbu empty space
	found:=false
	for i:=0;i<42;i++ {
		//upper
		found=true
		for j:=int(ypos);j<int(ypos+nbu);j++ {
			if i+j<0 || i+j>=42 || busy[i+j]==true { found=false}
		}
		if (found==true) {
			return int32(i)+ypos
		}
		
		//lower
		found=true
		for j:=int(ypos);j<int(ypos+nbu);j++ {
			if j-i<0 || j-i>=42 || busy[j-i]==true { found=false}
		}
		if (found==true) {
			return ypos-int32(i)
		}
	}
	return -1
}

func (self *RackChassisWidget) Repaint() {
	self.SWS_CoreWidget.Repaint()
	
	self.SetDrawColor(0, 0, 0, 255)
	self.DrawLine(9,9,110,9)
	self.DrawLine(110,9,110,430)
	self.DrawLine(110,430,9,430)
	self.DrawLine(9,430,9,9)
	
	self.FillRect(10, 10, 100,420, 0xffaaaaaa)
	
	if (self.ydrag!=-1) {
		ypos:=self.computeComingPos(self.ydrag)
		if ypos!=-1 {
			nbu:=self.comingitem.Serverconf.ConfType.NbU
			self.SetDrawColor(255, 0, 0, 255)
			self.DrawLine(10,10+ypos*10,19,10+ypos*10)
			self.DrawLine(10,10+ypos*10,10,19+ypos*10)

			self.DrawLine(109,10+ypos*10,100,10+ypos*10)
			self.DrawLine(109,10+ypos*10,109,19+ypos*10)

			self.DrawLine(10,9+(ypos+nbu)*10,19,9+(ypos+nbu)*10)
			self.DrawLine(10,9+(ypos+nbu)*10,10,0+(ypos+nbu)*10)

			self.DrawLine(109,9+(ypos+nbu)*10,100,9+(ypos+nbu)*10)
			self.DrawLine(109,9+(ypos+nbu)*10,109,0+(ypos+nbu)*10)
		}
	}
}

func (self *RackChassisWidget) DragMove(x,y int32, payload sws.DragPayload) {
	if payload.GetType()==1 {
		self.ydrag=y
		sws.PostUpdate()
	}
}

func (self *RackChassisWidget) DragEnter(x,y int32, payload sws.DragPayload) {
	if payload.GetType()==1 {
		self.ydrag=y
		self.comingitem=payload.(*ServerDragPayload).item
		sws.PostUpdate()
	}
}

func (self *RackChassisWidget) DragLeave() {
	self.ydrag=-1
	sws.PostUpdate()
}

func (self *RackChassisWidget) DragDrop(x,y int32, payload sws.DragPayload) {
	if payload.GetType()==1 {
		// TODO

		self.ydrag=-1
		sws.PostUpdate()
	}
}

func NewRackChassisWidget(inventory *supplier.Inventory) *RackChassisWidget {
	chassis:=&RackChassisWidget {
		SWS_CoreWidget: *sws.CreateCoreWidget(300,600),
		inventory: inventory,
		ydrag: -1,
		items: make([]*supplier.InventoryItem,0),
	}
	return chassis
}

//
// 2 zones:
// - one left creating DragEvent to populate the rack
//   -> create row that dont send back GetChildren()
// - one on the right that receive dragevent AND create DragEvent (to move, and to trash?)
//
type RackWidget struct {
	sws.SWS_CoreWidget
	mainwidget     *sws.SWS_MainWidget
	rootwindow     *sws.SWS_RootWidget
	inventory      *supplier.Inventory
	xactiveElement int32
	yactiveElement int32
	activeElement  DcElement
	vbox           *sws.SWS_VBoxWidget
	splitview      *sws.SWS_SplitviewWidget
}

func NewRackWidget(rootwindow *sws.SWS_RootWidget,inventory *supplier.Inventory) *RackWidget {
	mainwidget:=sws.CreateMainWidget(650,400," Rack info ",false,true)
	sv:=sws.CreateSplitviewWidget(400,300,true)
	
	rack:=&RackWidget{
		SWS_CoreWidget: *sws.CreateCoreWidget(650,400),
		mainwidget: mainwidget,
		rootwindow: rootwindow,
		inventory: inventory,
		xactiveElement: -1,
		yactiveElement: -1,
		activeElement: nil,
		vbox: sws.CreateVBoxWidget(200,300),
		splitview: sv,
	}
	inventory.AddSubscriber(rack)
	
	scrollleft:=sws.CreateScrollWidget(200,300)
	scrollleft.ShowHorizontalScrollbar(false)
	scrollleft.SetInnerWidget(rack.vbox)
	
	sv.PlaceSplitBar(200)
	sv.SetLeftWidget(scrollleft)

	scrollright:=sws.CreateScrollWidget(200,300)
	scrollright.SetInnerWidget(NewRackChassisWidget(inventory))
	sv.SetRightWidget(scrollright)
	
	rack.AddChild(sv)

	mainwidget.SetInnerWidget(rack)
	mainwidget.SetCloseCallback(func() {
		rack.Hide()
	})
	return rack
}

func (self *RackWidget) Resize(w,h int32) {
	self.SWS_CoreWidget.Resize(w,h)
	self.splitview.Resize(w,h)
}

func (self *RackWidget) Show() {
	self.rootwindow.AddChild(self.mainwidget)
	self.rootwindow.SetFocus(self.mainwidget)
}

func (self *RackWidget) Hide() {
	self.rootwindow.RemoveChild(self.mainwidget)
	children:=self.rootwindow.GetChildren()
	if len(children)>0 {
		self.rootwindow.SetFocus(children[0])
	}
}

func (self *RackWidget) ItemInTransit(item *supplier.InventoryItem) {
}

func (self *RackWidget) ItemInStock(item *supplier.InventoryItem) {
	fmt.Println("RackWidget::ItemInStock")
	if item.Typeitem==supplier.PRODUCT_SERVER && item.Serverconf.ConfType.NbU>0 && item.Xplaced==-1 {
		fmt.Println("RackWidget::ItemInStock b")
		self.vbox.AddChild(NewRackWidgetLine(item))
	}
}

//
// This widget allow to display a Datacenter map (and more)
//
type DcWidget struct {
	sws.SWS_CoreWidget
	rackwidget    *RackWidget
	tiles         [][]*Tile
	xRoot, yRoot  int32
	activeTile    *Tile
	te            *sws.TimerEvent
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
		self.rackwidget.mainwidget.SetTitle(fmt.Sprintf("Rack details %d/%d",self.rackwidget.xactiveElement,self.rackwidget.yactiveElement))
		self.rackwidget.Show()
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
			self.tiles[y][x] = CreateGrassTile()
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
		var dcelement map[string]interface{}
		if tile["dcelementtype"] != nil {
			dcelementtype = tile["dcelementtype"].(string)
		}
		if tile["dcelement"] != nil {
			dcelement = tile["dcelement"].(map[string]interface{})
		}
		if dcelementtype == "" || dcelementtype == "rack" {
			// basic floor
			self.tiles[y][x] = CreateElectricalTile(wall0, wall1, floor, rotation, dcelementtype, dcelement)
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
				value = fmt.Sprintf(`{"x":%d, "y":%d, "wall0":"%s", "wall1":"%s", "floor":"%s", "rotation":%d, "dcelementtype":"%s", "dcelement":%s}`,
					x,
					y,
					t.wall[0],
					t.wall[1],
					t.floor,
					t.rotation,
					t.element.ElementType(),
					t.element.Save(),
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

func CreateDcWidget(w, h int32,rootwindow *sws.SWS_RootWidget,inventory *supplier.Inventory) *DcWidget {
	corewidget := sws.CreateCoreWidget(w, h)
	rackwidget:=NewRackWidget(rootwindow,inventory)
	widget := &DcWidget{SWS_CoreWidget: *corewidget,
		rackwidget: rackwidget,
		tiles: [][]*Tile{{}},
		xRoot: 0,
		yRoot: 0,
	}
	return widget
}
