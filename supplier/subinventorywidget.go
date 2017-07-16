package supplier

import(
	"fmt"
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
	Checkbox  *sws.SWS_CheckboxWidget
	desc      *sws.SWS_Label
	placement *sws.SWS_Label
	item      *InventoryItem
}

func NewUnallocatedInventoryLineWidget(item *InventoryItem) *UnallocatedInventoryLineWidget{
//	ramSizeText:=fmt.Sprintf("%d Mo",item.Serverconf.NbSlotRam*item.Serverconf.RamSize)
//	if (item.Serverconf.NbSlotRam*item.Serverconf.RamSize>=2048) {
//		ramSizeText=fmt.Sprintf("%d Go",item.Serverconf.NbSlotRam*item.Serverconf.RamSize/1024)
//	}
//	text:=fmt.Sprintf("%dx %d cores\n%s RAM\n%d disks",item.Serverconf.NbProcessors,item.Serverconf.NbCore,ramSizeText,item.Serverconf.NbDisks)
	text:=item.Serverconf.ConfType.ServerName
	placement:=" - "
	if (item.xplaced!=-1) {
		placement=fmt.Sprintf("%d/%d",item.xplaced,item.yplaced)
	}
	line:=&UnallocatedInventoryLineWidget {
		SWS_CoreWidget: *sws.CreateCoreWidget(325, 25),
		Checkbox: sws.CreateCheckboxWidget(),
		desc: sws.CreateLabel(200,25,text),
		placement: sws.CreateLabel(200,25,placement),
		item: item,
	}
	line.Checkbox.SetColor(0xffffffff)
	line.AddChild(line.Checkbox)
	
	line.desc.SetColor(0xffffffff)
	line.desc.Move(25,0)
	line.AddChild(line.desc)
	
	line.placement.SetColor(0xffffffff)
	line.placement.Move(225,0)
	line.AddChild(line.placement)
	
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
	if item.Typeitem==PRODUCT_SERVER {
		self.vbox.AddChild(NewUnallocatedInventoryLineWidget(item))
		sws.PostUpdate()
	}
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
		vbox: sws.CreateVBoxWidget(325,0),
	}
	inventory.AddSubscriber(widget)
	
	globalcheckbox:=sws.CreateCheckboxWidget()
	widget.AddChild(globalcheckbox)
	globalcheckbox.SetClicked(func() {
		for _,child:=range(widget.vbox.GetChildren()) {
			line:=child.(*UnallocatedInventoryLineWidget)
			line.Checkbox.SetSelected(globalcheckbox.Selected)
		}
		sws.PostUpdate()
	})
	
	globaldesc:=sws.CreateLabel(200,25,"Description")
	globaldesc.Move(25,0)
	widget.AddChild(globaldesc)
	
	globalplacement:=sws.CreateLabel(100,25,"Placement")
	globalplacement.Move(225,0)
	widget.AddChild(globalplacement)
	
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
