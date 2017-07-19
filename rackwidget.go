package dctycoon

import (
	"fmt"
	"github.com/nzin/sws"
	"github.com/nzin/dctycoon/supplier"
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
		vbox: sws.CreateVBoxWidget(300,0),
		splitview: sv,
	}
	inventory.AddSubscriber(rack)
	
	scrollleft:=sws.CreateScrollWidget(300,300)
	scrollleft.ShowHorizontalScrollbar(false)
	scrollleft.SetInnerWidget(rack.vbox)
	scrollleft.SetColor(0xffffffff)
	
	sv.PlaceSplitBar(300)
	sv.SplitBarMovable(false)
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

