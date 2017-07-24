package dctycoon

import (
	"fmt"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/sws"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	CHASSIS_OFFSET = 35
	RACK_SIZE      = 15
)

type ServerDragPayload struct {
	item *supplier.InventoryItem
}

func (self *ServerDragPayload) GetType() int32 {
	return 1
}

func (self *ServerDragPayload) PayloadAccepted(bool) {
}

type ServerMovePayload struct {
	inventory *supplier.Inventory
	item      *supplier.InventoryItem
}

func (self *ServerMovePayload) GetType() int32 {
	return 2
}

func (self *ServerMovePayload) PayloadAccepted(accepted bool) {
	if accepted == false {
		self.inventory.UninstallItem(self.item)
	}
}

type RackWidgetLine struct {
	sws.LabelWidget
	item *supplier.InventoryItem
}

func (self *RackWidgetLine) MousePressDown(x, y int32, button uint8) {
	payload := &ServerDragPayload{
		item: self.item,
	}
	var parent sws.Widget
	parent = self
	for parent != nil {
		x += parent.X()
		y += parent.Y()
		parent = parent.Parent()
	}
	sws.NewDragEvent(x, y, "resources/"+self.item.Serverconf.ConfType.ServerSprite+"0.png", payload)
}

func NewRackWidgetLine(item *supplier.InventoryItem) *RackWidgetLine {
	label := sws.NewLabelWidget(300, 45, item.ShortDescription())
	label.SetImage("resources/" + item.Serverconf.ConfType.ServerSprite + "half.png")
	label.AlignImageLeft(true)
	label.SetColor(0xffffffff)

	return &RackWidgetLine{
		LabelWidget: *label,
		item:        item,
	}
}

type RackWidgetItems struct {
	sws.CoreWidget
	vbox   *sws.VBoxWidget
	scroll *sws.ScrollWidget
}

func NewRackWidgetItems(inventory *supplier.Inventory) *RackWidgetItems {
	widgetitems := &RackWidgetItems{
		CoreWidget: *sws.NewCoreWidget(300, 100),
		vbox:       sws.NewVBoxWidget(300, 0),
		scroll:     sws.NewScrollWidget(300, 300),
	}
	inventory.AddSubscriber(widgetitems)

	label := sws.NewLabelWidget(300, 25, "Available server to place: ")
	widgetitems.AddChild(label)

	widgetitems.scroll.ShowHorizontalScrollbar(false)
	widgetitems.scroll.SetInnerWidget(widgetitems.vbox)
	widgetitems.scroll.Move(0, 25)
	widgetitems.AddChild(widgetitems.scroll)

	return widgetitems
}

func (self *RackWidgetItems) Resize(w, h int32) {
	self.CoreWidget.Resize(w, h)
	self.scroll.Resize(w, h-25)
}

func (self *RackWidgetItems) ItemInTransit(item *supplier.InventoryItem) {
}

func (self *RackWidgetItems) ItemInStock(item *supplier.InventoryItem) {
	if item.Typeitem == supplier.PRODUCT_SERVER && item.Serverconf.ConfType.NbU > 0 && item.Xplaced == -1 {
		self.vbox.AddChild(NewRackWidgetLine(item))
	}
}

func (self *RackWidgetItems) ItemRemoveFromStock(item *supplier.InventoryItem) {
	for _, elt := range self.vbox.GetChildren() {
		line := elt.(*RackWidgetLine)
		if line.item == item {
			self.vbox.RemoveChild(elt)
		}
	}
}

func (self *RackWidgetItems) ItemInstalled(*supplier.InventoryItem) {
}

func (self *RackWidgetItems) ItemUninstalled(*supplier.InventoryItem) {
}

type RackChassisWidget struct {
	sws.CoreWidget
	inventory  *supplier.Inventory
	xpos       int32
	ypos       int32
	ydrag      int32
	items      []*supplier.InventoryItem
	inmove     *supplier.InventoryItem
	comingitem *supplier.InventoryItem
}

func (self *RackChassisWidget) ItemInTransit(*supplier.InventoryItem) {
}

func (self *RackChassisWidget) ItemInStock(*supplier.InventoryItem) {
}

func (self *RackChassisWidget) ItemRemoveFromStock(*supplier.InventoryItem) {
}

func (self *RackChassisWidget) ItemInstalled(item *supplier.InventoryItem) {
	if item.Xplaced == self.xpos && item.Yplaced == self.ypos && item.Typeitem == supplier.PRODUCT_SERVER {
		self.items = append(self.items, item)
		sws.PostUpdate()
	}
}

func (self *RackChassisWidget) ItemUninstalled(item *supplier.InventoryItem) {
	if item.Xplaced == self.xpos && item.Yplaced == self.ypos {
		for p, i := range self.items {
			if i == item {
				self.items = append(self.items[:p], self.items[p+1:]...)
			}
		}
		sws.PostUpdate()
	}
}

func (self *RackChassisWidget) SetLocation(x, y int32) {
	self.xpos = x
	self.ypos = y
	self.items = make([]*supplier.InventoryItem, 0)
	for _, item := range self.inventory.Items {
		if item.Xplaced == x && item.Yplaced == y && item.Typeitem == supplier.PRODUCT_SERVER {
			self.items = append(self.items, item)
		}
	}
}

func (self *RackChassisWidget) computeComingPos(zpos int32) int32 {
	zpos = (zpos - CHASSIS_OFFSET) / RACK_SIZE
	nbu := self.comingitem.Serverconf.ConfType.NbU
	var busy [42]bool

	// first create the map of what is filled
	for _, item := range self.items {
		if item == self.inmove {
			continue
		}
		itemNbU := item.Serverconf.ConfType.NbU
		for j := 0; j < int(itemNbU); j++ {
			if j+int(item.Zplaced) < 42 {
				busy[j+int(item.Zplaced)] = true
			}
		}
	}

	// now try to find a nbu empty space
	found := false
	for i := 0; i < 42; i++ {
		//upper
		found = true
		for j := int(zpos); j < int(zpos+nbu); j++ {
			if i+j < 0 || i+j >= 42 || busy[i+j] == true {
				found = false
			}
		}
		if found == true {
			return int32(i) + zpos
		}

		//lower
		found = true
		for j := int(zpos); j < int(zpos+nbu); j++ {
			if j-i < 0 || j-i >= 42 || busy[j-i] == true {
				found = false
			}
		}
		if found == true {
			return zpos - int32(i)
		}
	}
	return -1
}

func (self *RackChassisWidget) Repaint() {
	self.CoreWidget.Repaint()

	var watts float64
	for _, i := range self.items {
		watts += i.Serverconf.PowerConsumption()
	}
	rackwatt := fmt.Sprintf("Rack consumption: %.2f W", watts)
	self.WriteText(10, 10, rackwatt, sdl.Color{0, 0, 0, 255})

	self.SetDrawColor(0, 0, 0, 255)
	self.DrawLine(9, CHASSIS_OFFSET-1, 110, CHASSIS_OFFSET-1)
	self.DrawLine(110, CHASSIS_OFFSET-1, 110, 42*RACK_SIZE+CHASSIS_OFFSET)
	self.DrawLine(110, 42*RACK_SIZE+CHASSIS_OFFSET, 9, 42*RACK_SIZE+CHASSIS_OFFSET)
	self.DrawLine(9, 42*RACK_SIZE+CHASSIS_OFFSET, 9, CHASSIS_OFFSET-1)

	self.FillRect(10, CHASSIS_OFFSET, 100, 42*RACK_SIZE, 0xffaaaaaa)

	for _, i := range self.items {
		if i == self.inmove {
			continue
		}
		nbu := i.Serverconf.ConfType.NbU
		self.FillRect(10, CHASSIS_OFFSET+i.Zplaced*RACK_SIZE, 100, nbu*RACK_SIZE, 0xff000000)
		self.WriteText(120, CHASSIS_OFFSET+i.Zplaced*RACK_SIZE, i.ShortDescription(), sdl.Color{0, 0, 0, 255})
	}

	if self.ydrag != -1 {
		zpos := self.computeComingPos(self.ydrag)
		if zpos != -1 {
			nbu := self.comingitem.Serverconf.ConfType.NbU
			self.SetDrawColor(255, 0, 0, 255)
			self.DrawLine(10, CHASSIS_OFFSET+zpos*RACK_SIZE, 19, CHASSIS_OFFSET+zpos*RACK_SIZE)
			self.DrawLine(10, CHASSIS_OFFSET+zpos*RACK_SIZE, 10, CHASSIS_OFFSET+9+zpos*RACK_SIZE)

			self.DrawLine(109, CHASSIS_OFFSET+zpos*RACK_SIZE, 100, CHASSIS_OFFSET+zpos*RACK_SIZE)
			self.DrawLine(109, CHASSIS_OFFSET+zpos*RACK_SIZE, 109, CHASSIS_OFFSET+9+zpos*RACK_SIZE)

			self.DrawLine(10, CHASSIS_OFFSET-1+(zpos+nbu)*RACK_SIZE, 19, CHASSIS_OFFSET-1+(zpos+nbu)*RACK_SIZE)
			self.DrawLine(10, CHASSIS_OFFSET-1+(zpos+nbu)*RACK_SIZE, 10, CHASSIS_OFFSET-10+(zpos+nbu)*RACK_SIZE)

			self.DrawLine(109, CHASSIS_OFFSET-1+(zpos+nbu)*RACK_SIZE, 100, CHASSIS_OFFSET-1+(zpos+nbu)*RACK_SIZE)
			self.DrawLine(109, CHASSIS_OFFSET-1+(zpos+nbu)*RACK_SIZE, 109, CHASSIS_OFFSET-10+(zpos+nbu)*RACK_SIZE)
		}
	}
}

func (self *RackChassisWidget) MousePressDown(x, y int32, button uint8) {
	if x >= 10 && x < 110 && y >= CHASSIS_OFFSET && y < 42*RACK_SIZE+CHASSIS_OFFSET {
		var item *supplier.InventoryItem
		for _, i := range self.items {
			if y >= i.Zplaced*RACK_SIZE+CHASSIS_OFFSET && y < (i.Zplaced+i.Serverconf.ConfType.NbU)*RACK_SIZE+CHASSIS_OFFSET {
				item = i
			}
		}

		if item != nil {
			self.inmove = item
			payload := &ServerMovePayload{
				item:      item,
				inventory: self.inventory,
			}
			var parent sws.Widget
			parent = self
			for parent != nil {
				x += parent.X()
				y += parent.Y()
				parent = parent.Parent()
			}
			sws.NewDragEvent(x, y, "resources/"+item.Serverconf.ConfType.ServerSprite+"0.png", payload)
		}
	}
}

func (self *RackChassisWidget) MousePressUp(x, y int32, button uint8) {
	self.inmove = nil
}

func (self *RackChassisWidget) DragMove(x, y int32, payload sws.DragPayload) {
	if payload.GetType() == 1 || payload.GetType() == 2 {
		self.ydrag = y
		sws.PostUpdate()
	}
}

func (self *RackChassisWidget) DragEnter(x, y int32, payload sws.DragPayload) {
	if payload.GetType() == 1 {
		self.comingitem = payload.(*ServerDragPayload).item
	}
	if payload.GetType() == 1 || payload.GetType() == 2 {
		self.ydrag = y
		sws.PostUpdate()
	}
}

func (self *RackChassisWidget) DragLeave(payload sws.DragPayload) {
	if payload.GetType() == 1 || payload.GetType() == 2 {
		self.ydrag = -1
		sws.PostUpdate()
	}
}

func (self *RackChassisWidget) DragDrop(x, y int32, payload sws.DragPayload) bool {
	if payload.GetType() == 1 {
		zpos := self.computeComingPos(self.ydrag)
		if zpos != -1 {
			self.inventory.InstallItem(self.comingitem, self.xpos, self.ypos, zpos)
			self.ydrag = -1
			sws.PostUpdate()
			return true
		}

		self.ydrag = -1
		sws.PostUpdate()
	}
	if payload.GetType() == 2 {
		// we reset sel.inmove because MousePressUp disabled it
		item := payload.(*ServerMovePayload).item
		self.inmove = item
		zpos := self.computeComingPos(self.ydrag)
		self.inmove = nil
		if zpos == -1 {
			panic("not able to find back a place")
		}
		xpos := item.Xplaced
		ypos := item.Yplaced
		self.inventory.UninstallItem(item)
		self.inventory.InstallItem(item, xpos, ypos, zpos)
		self.ydrag = -1
		sws.PostUpdate()
		return true
	}
	return false
}

func NewRackChassisWidget(inventory *supplier.Inventory) *RackChassisWidget {
	chassis := &RackChassisWidget{
		CoreWidget: *sws.NewCoreWidget(420, 42*RACK_SIZE+10+CHASSIS_OFFSET),
		inventory:  inventory,
		ydrag:      -1,
		xpos:       -1,
		ypos:       -1,
		items:      make([]*supplier.InventoryItem, 0),
	}
	inventory.AddSubscriber(chassis)
	return chassis
}

//
// 2 zones:
// - one left creating DragEvent to populate the rack
//   -> create row that dont send back GetChildren()
// - one on the right that receive dragevent AND create DragEvent (to move, and to trash?)
//
type RackWidget struct {
	sws.CoreWidget
	mainwidget     *sws.MainWidget
	rootwindow     *sws.RootWidget
	inventory      *supplier.Inventory
	xactiveElement int32
	yactiveElement int32
	activeElement  TileElement
	splitview      *sws.SplitviewWidget
	rackchassis    *RackChassisWidget
}

func NewRackWidget(rootwindow *sws.RootWidget, inventory *supplier.Inventory) *RackWidget {
	mainwidget := sws.NewMainWidget(650, 400, " Rack info ", false, true)
	svBottom := sws.NewSplitviewWidget(400, 300, true)

	rack := &RackWidget{
		mainwidget:     mainwidget,
		rootwindow:     rootwindow,
		inventory:      inventory,
		xactiveElement: -1,
		yactiveElement: -1,
		activeElement:  nil,
		splitview:      svBottom,
		rackchassis:    NewRackChassisWidget(inventory),
	}

	sv := sws.NewSplitviewWidget(200, 200, false)
	sv.PlaceSplitBar(50)
	sv.SplitBarMovable(false)
	mainwidget.SetInnerWidget(sv)

	menu := sws.NewCoreWidget(500, 50)
	menu.SetColor(0xffffffff)
	sv.SetLeftWidget(menu)

	menuservers := sws.NewLabelWidget(300, 50, "Servers")
	menuservers.SetColor(0xffffffff)
	menuservers.SetCentered(true)
	menu.AddChild(menuservers)

	menurack := sws.NewLabelWidget(300, 50, "Rack chassis")
	menurack.SetColor(0xffffffff)
	menurack.Move(300, 0)
	menu.AddChild(menurack)

	svBottom.PlaceSplitBar(300)
	svBottom.SplitBarMovable(false)
	svBottom.SetLeftWidget(NewRackWidgetItems(inventory))

	scrollright := sws.NewScrollWidget(320, 300)
	scrollright.SetInnerWidget(rack.rackchassis)
	svBottom.SetRightWidget(scrollright)

	sv.SetRightWidget(svBottom)

	mainwidget.SetCloseCallback(func() {
		rack.Hide()
	})
	return rack
}

func (self *RackWidget) Show(x, y int32) {
	self.rootwindow.AddChild(self.mainwidget)
	self.rootwindow.SetFocus(self.mainwidget)
	self.mainwidget.SetTitle(fmt.Sprintf(" Rack details %d-%d ", x, y))
	self.xactiveElement = x
	self.yactiveElement = y
	self.rackchassis.SetLocation(x, y)
}

func (self *RackWidget) Hide() {
	self.rootwindow.RemoveChild(self.mainwidget)
	children := self.rootwindow.GetChildren()
	if len(children) > 0 {
		self.rootwindow.SetFocus(children[0])
	}
}