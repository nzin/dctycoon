package dctycoon

import (
	"fmt"
	"github.com/nzin/sws"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/nzin/dctycoon/supplier"
)

const (
	CHASSIS_OFFSET=35
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
	label:=sws.CreateLabel(300,45,item.ShortDescription())
	label.SetImage("resources/"+item.Serverconf.ConfType.ServerSprite+"half.png")
	label.AlignImageLeft(true)
	label.SetColor(0xffffffff)

	return &RackWidgetLine{
		SWS_Label: *label,
		item: item,
	}
}

type RackWidgetItems struct {
	sws.SWS_CoreWidget
	vbox   *sws.SWS_VBoxWidget
	scroll *sws.SWS_ScrollWidget
}

func NewRackWidgetItems(inventory *supplier.Inventory) *RackWidgetItems {
	widgetitems:=&RackWidgetItems{
		SWS_CoreWidget: *sws.CreateCoreWidget(300,100),
		vbox: sws.CreateVBoxWidget(300,0),
		scroll: sws.CreateScrollWidget(300,300),
	}
	inventory.AddSubscriber(widgetitems)
	
	label:=sws.CreateLabel(300,25,"Available server to place: ")
	widgetitems.AddChild(label)
	
	widgetitems.scroll.ShowHorizontalScrollbar(false)
	widgetitems.scroll.SetInnerWidget(widgetitems.vbox)
	widgetitems.scroll.Move(0,25)
	widgetitems.AddChild(widgetitems.scroll)
	
	return widgetitems
}

func (self *RackWidgetItems) Resize(w,h int32) {
	self.scroll.Resize(w,h-25)
}

func (self *RackWidgetItems) ItemInTransit(item *supplier.InventoryItem) {
}

func (self *RackWidgetItems) ItemInStock(item *supplier.InventoryItem) {
	if item.Typeitem==supplier.PRODUCT_SERVER && item.Serverconf.ConfType.NbU>0 && item.Xplaced==-1 {
		self.vbox.AddChild(NewRackWidgetLine(item))
	}
}

func (self *RackWidgetItems) ItemRemoveFromStock(item *supplier.InventoryItem) {
	for _,elt := range(self.vbox.GetChildren()) {
		line := elt.(*RackWidgetLine)
		if line.item==item {
			self.vbox.RemoveChild(elt)
		}
	}
}

func (self *RackWidgetItems) ItemInstalled(*supplier.InventoryItem) {
}

func (self *RackWidgetItems) ItemUninstalled(*supplier.InventoryItem) {
}

type RackChassisWidget struct {
	sws.SWS_CoreWidget
	inventory  *supplier.Inventory
	xpos       int32
	ypos       int32
	ydrag      int32
	items      []*supplier.InventoryItem
	comingitem *supplier.InventoryItem
}

func (self *RackChassisWidget) SetLocation(x,y int32) {
	self.xpos=x
	self.ypos=y
}

func (self *RackChassisWidget) computeComingPos(zpos int32) int32 {
	zpos=(zpos-CHASSIS_OFFSET)/10
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
		for j:=int(zpos);j<int(zpos+nbu);j++ {
			if i+j<0 || i+j>=42 || busy[i+j]==true { found=false}
		}
		if (found==true) {
			return int32(i)+zpos
		}
		
		//lower
		found=true
		for j:=int(zpos);j<int(zpos+nbu);j++ {
			if j-i<0 || j-i>=42 || busy[j-i]==true { found=false}
		}
		if (found==true) {
			return zpos-int32(i)
		}
	}
	return -1
}

func (self *RackChassisWidget) Repaint() {
	self.SWS_CoreWidget.Repaint()
	
	var watts float64
	for _,i := range(self.items) {
		watts+=i.Serverconf.PowerConsumption()
	}
	rackwatt:=fmt.Sprintf("Rack consumption: %.2f W",watts)
	self.WriteText(10,10,rackwatt,sdl.Color{0,0,0,255})
	
	self.SetDrawColor(0, 0, 0, 255)
	self.DrawLine(9,CHASSIS_OFFSET-1,110,CHASSIS_OFFSET-1)
	self.DrawLine(110,CHASSIS_OFFSET-1,110,420+CHASSIS_OFFSET)
	self.DrawLine(110,420+CHASSIS_OFFSET,9,420+CHASSIS_OFFSET)
	self.DrawLine(9,420+CHASSIS_OFFSET,9,CHASSIS_OFFSET-1)
	
	self.FillRect(10, CHASSIS_OFFSET, 100,420, 0xffaaaaaa)
	
	if (self.ydrag!=-1) {
		zpos:=self.computeComingPos(self.ydrag)
		if zpos!=-1 {
			nbu:=self.comingitem.Serverconf.ConfType.NbU
			self.SetDrawColor(255, 0, 0, 255)
			self.DrawLine(10,CHASSIS_OFFSET+zpos*10,19,CHASSIS_OFFSET+zpos*10)
			self.DrawLine(10,CHASSIS_OFFSET+zpos*10,10,CHASSIS_OFFSET+9+zpos*10)

			self.DrawLine(109,CHASSIS_OFFSET+zpos*10,100,CHASSIS_OFFSET+zpos*10)
			self.DrawLine(109,CHASSIS_OFFSET+zpos*10,109,CHASSIS_OFFSET+9+zpos*10)

			self.DrawLine(10,CHASSIS_OFFSET-1+(zpos+nbu)*10,19,CHASSIS_OFFSET-1+(zpos+nbu)*10)
			self.DrawLine(10,CHASSIS_OFFSET-1+(zpos+nbu)*10,10,CHASSIS_OFFSET-10+(zpos+nbu)*10)

			self.DrawLine(109,CHASSIS_OFFSET-1+(zpos+nbu)*10,100,CHASSIS_OFFSET-1+(zpos+nbu)*10)
			self.DrawLine(109,CHASSIS_OFFSET-1+(zpos+nbu)*10,109,CHASSIS_OFFSET-10+(zpos+nbu)*10)
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
		zpos:=self.computeComingPos(self.ydrag)
		if zpos!=-1 {
			self.inventory.InstallItem(self.comingitem,self.xpos,self.ypos,zpos)
		}

		self.ydrag=-1
		sws.PostUpdate()
	}
}

func NewRackChassisWidget(inventory *supplier.Inventory) *RackChassisWidget {
	chassis:=&RackChassisWidget {
		SWS_CoreWidget: *sws.CreateCoreWidget(220,430+CHASSIS_OFFSET),
		inventory: inventory,
		ydrag: -1,
		xpos: -1,
		ypos: -1,
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
	splitview      *sws.SWS_SplitviewWidget
	rackchassis    *RackChassisWidget
}

func NewRackWidget(rootwindow *sws.SWS_RootWidget,inventory *supplier.Inventory) *RackWidget {
	mainwidget:=sws.CreateMainWidget(650,400," Rack info ",false,true)
	svBottom:=sws.CreateSplitviewWidget(400,300,true)
	
	rack:=&RackWidget{
		mainwidget: mainwidget,
		rootwindow: rootwindow,
		inventory: inventory,
		xactiveElement: -1,
		yactiveElement: -1,
		activeElement: nil,
		splitview: svBottom,
		rackchassis: NewRackChassisWidget(inventory),
	}

	sv := sws.CreateSplitviewWidget(200,200,false)
	sv.PlaceSplitBar(50)
	sv.SplitBarMovable(false)
	mainwidget.SetInnerWidget(sv)

	menu:=sws.CreateCoreWidget(500,50)
	menu.SetColor(0xffffffff)
	sv.SetLeftWidget(menu)
	
	menuservers:=sws.CreateLabel(300,50,"Servers")
	menuservers.SetColor(0xffffffff)
	menuservers.SetCentered(true)
	menu.AddChild(menuservers)
	
	menurack:=sws.CreateLabel(300,50,"Rack chassis")
	menurack.SetColor(0xffffffff)
	menurack.Move(300,0)
	menu.AddChild(menurack)
	
	svBottom.PlaceSplitBar(300)
	svBottom.SplitBarMovable(false)
	svBottom.SetLeftWidget(NewRackWidgetItems(inventory))

	scrollright:=sws.CreateScrollWidget(200,300)
	scrollright.SetInnerWidget(rack.rackchassis)
	svBottom.SetRightWidget(scrollright)
	
	sv.SetRightWidget(svBottom)
	
	mainwidget.SetCloseCallback(func() {
		rack.Hide()
	})
	return rack
}

func (self *RackWidget) Show(x,y int32) {
	self.rootwindow.AddChild(self.mainwidget)
	self.rootwindow.SetFocus(self.mainwidget)
	self.mainwidget.SetTitle(fmt.Sprintf(" Rack details %d-%d ",x,y))
	self.xactiveElement=x
	self.yactiveElement=y
	self.rackchassis.SetLocation(x,y)
}

func (self *RackWidget) Hide() {
	self.rootwindow.RemoveChild(self.mainwidget)
	children:=self.rootwindow.GetChildren()
	if len(children)>0 {
		self.rootwindow.SetFocus(children[0])
	}
}
