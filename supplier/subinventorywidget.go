package supplier

import(
	"fmt"
	"github.com/nzin/sws"
)

type SubInventory struct {
	Icon            string
	Title           string
	ButtonPanel     *sws.SWS_CoreWidget
	Widget          sws.SWS_Widget
}

func NewUnallocatedInventorySub(inventory *Inventory) *SubInventory{
	sub:=&SubInventory{
		Icon: "resources/icon-delivery-truck-silhouette.png",
		Title: "Unallocated inventory",
		ButtonPanel: sws.CreateCoreWidget(200,50),
		Widget: sws.CreateScrollWidget(100, 100),
	}
	return sub
}

type UnallocatedServerLineWidget struct {
	sws.SWS_CoreWidget
	Checkbox  *sws.SWS_CheckboxWidget
	desc      *sws.SWS_Label
	placement *sws.SWS_Label
	item      *InventoryItem
}

func NewUnallocatedServerLineWidget(item *InventoryItem) *UnallocatedServerLineWidget{
	ramSizeText:=fmt.Sprintf("%d Mo",item.Serverconf.NbSlotRam*item.Serverconf.RamSize)
	if (item.Serverconf.NbSlotRam*item.Serverconf.RamSize>=2048) {
		ramSizeText=fmt.Sprintf("%d Go",item.Serverconf.NbSlotRam*item.Serverconf.RamSize/1024)
	}
	text:=item.Serverconf.ConfType.ServerName
	placement:=" - "
	if (item.Xplaced!=-1) {
		placement=fmt.Sprintf("%d/%d",item.Xplaced,item.Yplaced)
	}
	line:=&UnallocatedServerLineWidget {
		SWS_CoreWidget: *sws.CreateCoreWidget(625, 25),
		Checkbox: sws.CreateCheckboxWidget(),
		desc: sws.CreateLabel(200,25,text),
		placement: sws.CreateLabel(100,25,placement),
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
	
	cores:=sws.CreateLabel(100,25,fmt.Sprintf("%d",item.Serverconf.NbProcessors*item.Serverconf.NbCore))
	cores.SetColor(0xffffffff)
	cores.Move(325,0)
	line.AddChild(cores)
	
	ram:=sws.CreateLabel(100,25,ramSizeText)
	ram.SetColor(0xffffffff)
	ram.Move(425,0)
	line.AddChild(ram)
	
	diskText:=fmt.Sprintf("%d Mo",item.Serverconf.NbDisks*item.Serverconf.DiskSize)
	if item.Serverconf.NbDisks*item.Serverconf.DiskSize>4096 {
		diskText=fmt.Sprintf("%d Go",item.Serverconf.NbDisks*item.Serverconf.DiskSize/1024)
	}
	if item.Serverconf.NbDisks*item.Serverconf.DiskSize>4*1024*1024 {
		diskText=fmt.Sprintf("%d To",item.Serverconf.NbDisks*item.Serverconf.DiskSize/(1024*1024))
	}
	disk:=sws.CreateLabel(100,25,diskText)
	disk.SetColor(0xffffffff)
	disk.Move(525,0)
	line.AddChild(disk)
	
	return line
}

type UnallocatedServerWidget struct {
	sws.SWS_CoreWidget
	inventory   *Inventory
	scroll      *sws.SWS_ScrollWidget
	vbox        *sws.SWS_VBoxWidget
	globalcheckbox *sws.SWS_CheckboxWidget
	selected    map[*UnallocatedServerLineWidget]bool
	buttonPanel *sws.SWS_CoreWidget
	scrap       *sws.SWS_ButtonWidget
}

func (self *UnallocatedServerWidget) SelectLine(line *UnallocatedServerLineWidget,selected bool) {
	if (selected) {
		if len(self.selected)==0 {
			self.scrap.Move(0,12)
			self.buttonPanel.AddChild(self.scrap)
		}
		self.selected[line]=true
	} else {
		delete(self.selected,line)
		if len(self.selected)==0 {
			self.buttonPanel.RemoveChild(self.scrap)
		}
	}
}

func (self *UnallocatedServerWidget) ItemInTransit(*InventoryItem) {
}

func (self *UnallocatedServerWidget) ItemInStock(item *InventoryItem) {
	if item.Typeitem==PRODUCT_SERVER {
		line:=NewUnallocatedServerLineWidget(item)
		line.Checkbox.SetClicked(func() {
			self.SelectLine(line,line.Checkbox.Selected)
			self.globalcheckbox.SetSelected(false)
		})
		self.vbox.AddChild(line)
		sws.PostUpdate()
	}
}

func (self *UnallocatedServerWidget) Resize(w,h int32) {
	self.SWS_CoreWidget.Resize(w,h)
	h-=25
	self.scroll.Resize(w,h)
}

func NewUnallocatedServerSub(inventory *Inventory) *SubInventory{
	widget:=&UnallocatedServerWidget{
		SWS_CoreWidget: *sws.CreateCoreWidget(100, 100),
		inventory: inventory,
		scroll: sws.CreateScrollWidget(100,100),
		vbox: sws.CreateVBoxWidget(625,0),
		globalcheckbox: sws.CreateCheckboxWidget(),
		selected: make(map[*UnallocatedServerLineWidget]bool),
                buttonPanel: sws.CreateCoreWidget(200,50),
		scrap: sws.CreateButtonWidget(100,25,"Scrap"),
	}
	inventory.AddSubscriber(widget)
	
	widget.AddChild(widget.globalcheckbox)
	widget.globalcheckbox.SetClicked(func() {
		for _,child:=range(widget.vbox.GetChildren()) {
			line:=child.(*UnallocatedServerLineWidget)
			line.Checkbox.SetSelected(widget.globalcheckbox.Selected)
			widget.SelectLine(line,line.Checkbox.Selected)
		}
		sws.PostUpdate()
	})
	
	globaldesc:=sws.CreateLabel(200,25,"Description")
	globaldesc.Move(25,0)
	widget.AddChild(globaldesc)
	
	globalplacement:=sws.CreateLabel(100,25,"Placement")
	globalplacement.Move(225,0)
	widget.AddChild(globalplacement)
	
	globalnbcores:=sws.CreateLabel(100,25,"Nb cores")
	globalnbcores.Move(325,0)
	widget.AddChild(globalnbcores)
	
	globalram:=sws.CreateLabel(100,25,"RAM")
	globalram.Move(425,0)
	widget.AddChild(globalram)
	
	globaldisk:=sws.CreateLabel(100,25,"Disk")
	globaldisk.Move(525,0)
	widget.AddChild(globaldisk)
	
	widget.scroll.Move(0,25)
	widget.scroll.ShowHorizontalScrollbar(false)
	widget.scroll.SetInnerWidget(widget.vbox)
	widget.AddChild(widget.scroll)
	
	sub:=&SubInventory{
		Icon: "resources/icon-hard-drive.png",
		Title: "Unallocated servers",
		ButtonPanel: widget.buttonPanel,
		Widget: widget,
	}
	return sub
}

func NewPoolSub(inventory *Inventory) *SubInventory{
	sub:=&SubInventory{
		Icon: "resources/icon-bucket.png",
		Title: "Server pools",
		ButtonPanel: sws.CreateCoreWidget(200,50),
		Widget: sws.CreateScrollWidget(100, 100),
	}
	return sub
}

func NewOfferSub(inventory *Inventory) *SubInventory{
	sub:=&SubInventory{
		Icon: "resources/icon-paper-bill.png",
		Title: "Server offers",
		ButtonPanel: sws.CreateCoreWidget(200,50),
		Widget: sws.CreateScrollWidget(100, 100),
	}
	return sub
}
