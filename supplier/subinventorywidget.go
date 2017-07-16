package supplier

import(
	//"fmt"
	"github.com/nzin/sws"
)

type SubInventory struct {
	Icon    string
	Title   string
	Buttons []*sws.SWS_ButtonWidget
	Widget  sws.SWS_Widget
}

type UnallocatedInventoryLineWidget struct {
	sws.SWS_CoreWidget
	checkbox *sws.SWS_CheckboxWidget
	desc     *sws.SWS_Label
	item     *InventoryItem
}

func NewUnallocatedInventoryLineWidget(item *InventoryItem) *UnallocatedInventoryLineWidget{
	text:="desc"
	line:=&UnallocatedInventoryLineWidget {
		SWS_CoreWidget: *sws.CreateCoreWidget(225, 25),
		checkbox: sws.CreateCheckboxWidget(),
		desc: sws.CreateLabel(200,25,text),
		item: item,
	}
	line.checkbox.SetColor(0xffffffff)
	line.AddChild(line.checkbox)
	
	line.desc.SetColor(0xffffffff)
	line.desc.Move(25,0)
	line.AddChild(line.desc)
	
	return line
}

type UnallocatedInventoryWidget struct {
	sws.SWS_CoreWidget
	inventory *Inventory
	scroll    *sws.SWS_ScrollWidget
	vbox      *sws.SWS_VBoxWidget
}

func (self *UnallocatedInventoryWidget) ItemInTransit(*InventoryItem) {
}

func (self *UnallocatedInventoryWidget) ItemInStock(item *InventoryItem) {
	self.vbox.AddChild(NewUnallocatedInventoryLineWidget(item))
	sws.PostUpdate()
}

func (self *UnallocatedInventoryWidget) Resize(w,h int32) {
	self.SWS_CoreWidget.Resize(w,h)
	h-=25
	self.scroll.Resize(w,h)
}

func NewUnallocatedInventorySub(inventory *Inventory) *SubInventory{
	widget:=&UnallocatedInventoryWidget{
		SWS_CoreWidget: *sws.CreateCoreWidget(100, 100),
		inventory: inventory,
		scroll: sws.CreateScrollWidget(100,100),
		vbox: sws.CreateVBoxWidget(225,0),
	}
	inventory.AddSubscriber(widget)
	
	globalcheckbox:=sws.CreateCheckboxWidget()
	widget.AddChild(globalcheckbox)
	
	globaldesc:=sws.CreateLabel(200,25,"Description")
	globaldesc.Move(25,0)
	widget.AddChild(globaldesc)
	
	widget.scroll.Move(0,25)
	widget.scroll.SetInnerWidget(widget.vbox)
	widget.AddChild(widget.scroll)
	
	sub:=&SubInventory{
		Icon: "resources/icon-delivery-truck-silhouette.png",
		Title: "Unallocated inventory",
		Buttons: make([]*sws.SWS_ButtonWidget,0),
		Widget: widget,
	}
	return sub
}

func NewUnallocatedServerSub(inventory *Inventory) *SubInventory{
	sub:=&SubInventory{
		Icon: "resources/icon-hard-drive.png",
		Title: "Unallocated servers",
		Buttons: make([]*sws.SWS_ButtonWidget,0),
		Widget: sws.CreateScrollWidget(100, 100),
	}
	return sub
}

func NewPoolSub(inventory *Inventory) *SubInventory{
	sub:=&SubInventory{
		Icon: "resources/icon-bucket.png",
		Title: "Server pools",
		Buttons: make([]*sws.SWS_ButtonWidget,0),
		Widget: sws.CreateScrollWidget(100, 100),
	}
	return sub
}

func NewOfferSub(inventory *Inventory) *SubInventory{
	sub:=&SubInventory{
		Icon: "resources/icon-paper-bill.png",
		Title: "Server offers",
		Buttons: make([]*sws.SWS_ButtonWidget,0),
		Widget: sws.CreateScrollWidget(100, 100),
	}
	return sub
}
