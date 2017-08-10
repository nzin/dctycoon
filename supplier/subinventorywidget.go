package supplier

import (
	"fmt"
	"github.com/nzin/sws"
)

//
// base class for the inventory window sub-widget (i.e. right part)
//
type SubInventory struct {
	Icon        string
	Title       string
	ButtonPanel *sws.CoreWidget
	Widget      sws.Widget
}

type UnallocatedInventoryLineWidget struct {
	sws.CoreWidget
	Checkbox  *sws.CheckboxWidget
	desc      *sws.LabelWidget
	placement *sws.LabelWidget
	item      *InventoryItem
}

func NewUnallocatedInventoryLineWidget(item *InventoryItem) *UnallocatedInventoryLineWidget {
	placement := " - "
	if item.Xplaced != -1 {
		placement = fmt.Sprintf("%d/%d", item.Xplaced, item.Yplaced)
	}
	line := &UnallocatedInventoryLineWidget{
		CoreWidget: *sws.NewCoreWidget(625, 25),
		Checkbox:   sws.NewCheckboxWidget(),
		desc:       sws.NewLabelWidget(200, 25, item.ShortDescription()),
		placement:  sws.NewLabelWidget(100, 25, placement),
		item:       item,
	}
	line.Checkbox.SetColor(0xffffffff)
	line.AddChild(line.Checkbox)

	line.desc.SetColor(0xffffffff)
	line.desc.Move(25, 0)
	line.AddChild(line.desc)

	line.placement.SetColor(0xffffffff)
	line.placement.Move(225, 0)
	line.AddChild(line.placement)

	return line
}

type UnallocatedInventoryWidget struct {
	sws.CoreWidget
	inventory      *Inventory
	scroll         *sws.ScrollWidget
	vbox           *sws.VBoxWidget
	globalcheckbox *sws.CheckboxWidget
	selected       map[*UnallocatedInventoryLineWidget]bool
	buttonPanel    *sws.CoreWidget
	scrap          *sws.ButtonWidget
}

func (self *UnallocatedInventoryWidget) SelectLine(line *UnallocatedInventoryLineWidget, selected bool) {
	if selected {
		if len(self.selected) == 0 {
			self.scrap.Move(0, 2)
			self.buttonPanel.AddChild(self.scrap)
		}
		self.selected[line] = true
	} else {
		delete(self.selected, line)
		if len(self.selected) == 0 {
			self.buttonPanel.RemoveChild(self.scrap)
		}
	}
}

func (self *UnallocatedInventoryWidget) ItemInTransit(*InventoryItem) {
}

func (self *UnallocatedInventoryWidget) ItemInStock(item *InventoryItem) {
	if item.Typeitem == PRODUCT_RACK ||
		item.Typeitem == PRODUCT_AC ||
		item.Typeitem == PRODUCT_GENERATOR {
		line := NewUnallocatedInventoryLineWidget(item)
		line.Checkbox.SetClicked(func() {
			self.SelectLine(line, line.Checkbox.Selected)
			self.globalcheckbox.SetSelected(false)
		})
		self.vbox.AddChild(line)
		sws.PostUpdate()
	}
}

func (self *UnallocatedInventoryWidget) ItemRemoveFromStock(item *InventoryItem) {
	for _, elt := range self.vbox.GetChildren() {
		line := elt.(*UnallocatedInventoryLineWidget)
		if line.item == item {
			self.vbox.RemoveChild(elt)
		}
	}
}

func (self *UnallocatedInventoryWidget) ItemInstalled(*InventoryItem) {
}

func (self *UnallocatedInventoryWidget) ItemUninstalled(*InventoryItem) {
}

func (self *UnallocatedInventoryWidget) Resize(w, h int32) {
	self.CoreWidget.Resize(w, h)
	h -= 25
	self.scroll.Resize(w, h)
}

func NewUnallocatedInventorySub(root *sws.RootWidget, inventory *Inventory) *SubInventory {
	widget := &UnallocatedInventoryWidget{
		CoreWidget:     *sws.NewCoreWidget(100, 100),
		inventory:      inventory,
		scroll:         sws.NewScrollWidget(100, 100),
		vbox:           sws.NewVBoxWidget(625, 0),
		globalcheckbox: sws.NewCheckboxWidget(),
		selected:       make(map[*UnallocatedInventoryLineWidget]bool),
		buttonPanel:    sws.NewCoreWidget(200, 30),
		scrap:          sws.NewButtonWidget(100, 25, "Scrap"),
	}
	inventory.AddInventorySubscriber(widget)

	widget.AddChild(widget.globalcheckbox)
	widget.globalcheckbox.SetClicked(func() {
		for _, child := range widget.vbox.GetChildren() {
			line := child.(*UnallocatedInventoryLineWidget)
			line.Checkbox.SetSelected(widget.globalcheckbox.Selected)
			widget.SelectLine(line, line.Checkbox.Selected)
		}
		sws.PostUpdate()
	})

	globaldesc := sws.NewLabelWidget(200, 25, "Description")
	globaldesc.Move(25, 0)
	widget.AddChild(globaldesc)

	globalplacement := sws.NewLabelWidget(100, 25, "Placement")
	globalplacement.Move(225, 0)
	widget.AddChild(globalplacement)

	widget.scroll.Move(0, 25)
	widget.scroll.ShowHorizontalScrollbar(false)
	widget.scroll.SetInnerWidget(widget.vbox)
	widget.AddChild(widget.scroll)
	
	widget.scrap.SetClicked(func() {
		sws.ShowModalYesNo(root,"Scrap items","resources/icon-triangular-big.png","Do you really want to scrap all of them???",func() {
			for k,_ := range(widget.selected) {
				inventory.DiscardItem(k.item)
				delete(widget.selected,k)
				widget.vbox.RemoveChild(k)
			}
			widget.globalcheckbox.SetSelected(false)
			widget.buttonPanel.RemoveChild(widget.scrap)
		},nil)
	})

	sub := &SubInventory{
		Icon:        "resources/icon-delivery-truck-silhouette.png",
		Title:       "Unallocated inventory",
		ButtonPanel: widget.buttonPanel,
		Widget:      widget,
	}
	return sub
}

type UnallocatedServerLineWidget struct {
	sws.CoreWidget
	Checkbox  *sws.CheckboxWidget
	desc      *sws.LabelWidget
	placement *sws.LabelWidget
	item      *InventoryItem
}

func NewUnallocatedServerLineWidget(item *InventoryItem) *UnallocatedServerLineWidget {
	ramSizeText := fmt.Sprintf("%d Mo", item.Serverconf.NbSlotRam*item.Serverconf.RamSize)
	if item.Serverconf.NbSlotRam*item.Serverconf.RamSize >= 2048 {
		ramSizeText = fmt.Sprintf("%d Go", item.Serverconf.NbSlotRam*item.Serverconf.RamSize/1024)
	}
	text := item.Serverconf.ConfType.ServerName
	placement := " - "
	if item.Xplaced != -1 {
		placement = fmt.Sprintf("%d/%d", item.Xplaced, item.Yplaced)
	}
	line := &UnallocatedServerLineWidget{
		CoreWidget: *sws.NewCoreWidget(625, 25),
		Checkbox:   sws.NewCheckboxWidget(),
		desc:       sws.NewLabelWidget(200, 25, text),
		placement:  sws.NewLabelWidget(100, 25, placement),
		item:       item,
	}
	line.Checkbox.SetColor(0xffffffff)
	line.AddChild(line.Checkbox)

	line.desc.SetColor(0xffffffff)
	line.desc.Move(25, 0)
	line.AddChild(line.desc)

	line.placement.SetColor(0xffffffff)
	line.placement.Move(225, 0)
	line.AddChild(line.placement)

	cores := sws.NewLabelWidget(100, 25, fmt.Sprintf("%d", item.Serverconf.NbProcessors*item.Serverconf.NbCore))
	cores.SetColor(0xffffffff)
	cores.Move(325, 0)
	line.AddChild(cores)

	ram := sws.NewLabelWidget(100, 25, ramSizeText)
	ram.SetColor(0xffffffff)
	ram.Move(425, 0)
	line.AddChild(ram)

	diskText := fmt.Sprintf("%d Mo", item.Serverconf.NbDisks*item.Serverconf.DiskSize)
	if item.Serverconf.NbDisks*item.Serverconf.DiskSize > 4096 {
		diskText = fmt.Sprintf("%d Go", item.Serverconf.NbDisks*item.Serverconf.DiskSize/1024)
	}
	if item.Serverconf.NbDisks*item.Serverconf.DiskSize > 4*1024*1024 {
		diskText = fmt.Sprintf("%d To", item.Serverconf.NbDisks*item.Serverconf.DiskSize/(1024*1024))
	}
	disk := sws.NewLabelWidget(100, 25, diskText)
	disk.SetColor(0xffffffff)
	disk.Move(525, 0)
	line.AddChild(disk)

	return line
}

type UnallocatedServerWidget struct {
	sws.CoreWidget
	inventory      *Inventory
	scroll         *sws.ScrollWidget
	vbox           *sws.VBoxWidget
	globalcheckbox *sws.CheckboxWidget
	selected       map[*UnallocatedServerLineWidget]bool
	buttonPanel    *sws.CoreWidget
	scrap          *sws.ButtonWidget
}

func (self *UnallocatedServerWidget) SelectLine(line *UnallocatedServerLineWidget, selected bool) {
	if selected {
		if len(self.selected) == 0 {
			self.scrap.Move(0, 2)
			self.buttonPanel.AddChild(self.scrap)
		}
		self.selected[line] = true
	} else {
		delete(self.selected, line)
		if len(self.selected) == 0 {
			self.buttonPanel.RemoveChild(self.scrap)
		}
	}
}

func (self *UnallocatedServerWidget) ItemInTransit(*InventoryItem) {
}

func (self *UnallocatedServerWidget) ItemInStock(item *InventoryItem) {
	if item.Typeitem == PRODUCT_SERVER {
		line := NewUnallocatedServerLineWidget(item)
		line.Checkbox.SetClicked(func() {
			self.SelectLine(line, line.Checkbox.Selected)
			self.globalcheckbox.SetSelected(false)
		})
		self.vbox.AddChild(line)
		sws.PostUpdate()
	}
}

func (self *UnallocatedServerWidget) ItemRemoveFromStock(item *InventoryItem) {
	for _, elt := range self.vbox.GetChildren() {
		line := elt.(*UnallocatedServerLineWidget)
		if line.item == item {
			self.vbox.RemoveChild(elt)
		}
	}
}

func (self *UnallocatedServerWidget) ItemInstalled(*InventoryItem) {
}

func (self *UnallocatedServerWidget) ItemUninstalled(*InventoryItem) {
}

func (self *UnallocatedServerWidget) Resize(w, h int32) {
	self.CoreWidget.Resize(w, h)
	h -= 25
	self.scroll.Resize(w, h)
}

func NewUnallocatedServerSub(root *sws.RootWidget, inventory *Inventory) *SubInventory {
	widget := &UnallocatedServerWidget{
		CoreWidget:     *sws.NewCoreWidget(100, 100),
		inventory:      inventory,
		scroll:         sws.NewScrollWidget(100, 100),
		vbox:           sws.NewVBoxWidget(625, 0),
		globalcheckbox: sws.NewCheckboxWidget(),
		selected:       make(map[*UnallocatedServerLineWidget]bool),
		buttonPanel:    sws.NewCoreWidget(200, 30),
		scrap:          sws.NewButtonWidget(100, 25, "Scrap"),
	}
	inventory.AddInventorySubscriber(widget)

	widget.AddChild(widget.globalcheckbox)
	widget.globalcheckbox.SetClicked(func() {
		for _, child := range widget.vbox.GetChildren() {
			line := child.(*UnallocatedServerLineWidget)
			line.Checkbox.SetSelected(widget.globalcheckbox.Selected)
			widget.SelectLine(line, line.Checkbox.Selected)
		}
		sws.PostUpdate()
	})

	globaldesc := sws.NewLabelWidget(200, 25, "Description")
	globaldesc.Move(25, 0)
	widget.AddChild(globaldesc)

	globalplacement := sws.NewLabelWidget(100, 25, "Placement")
	globalplacement.Move(225, 0)
	widget.AddChild(globalplacement)

	globalnbcores := sws.NewLabelWidget(100, 25, "Nb cores")
	globalnbcores.Move(325, 0)
	widget.AddChild(globalnbcores)

	globalram := sws.NewLabelWidget(100, 25, "RAM")
	globalram.Move(425, 0)
	widget.AddChild(globalram)

	globaldisk := sws.NewLabelWidget(100, 25, "Disk")
	globaldisk.Move(525, 0)
	widget.AddChild(globaldisk)

	widget.scroll.Move(0, 25)
	widget.scroll.ShowHorizontalScrollbar(false)
	widget.scroll.SetInnerWidget(widget.vbox)
	widget.AddChild(widget.scroll)

	widget.scrap.SetClicked(func() {
		sws.ShowModalYesNo(root,"Scrap items","resources/icon-triangular-big.png","Do you really want to scrap all of them???",func() {
			for k,_ := range(widget.selected) {
				inventory.DiscardItem(k.item)
				delete(widget.selected,k)
				widget.vbox.RemoveChild(k)
			}
			widget.globalcheckbox.SetSelected(false)
			widget.buttonPanel.RemoveChild(widget.scrap)
		},nil)
	})

	sub := &SubInventory{
		Icon:        "resources/icon-hard-drive.png",
		Title:       "Unallocated servers",
		ButtonPanel: widget.buttonPanel,
		Widget:      widget,
	}
	return sub
}

type PoolCreateWidget struct {
	rootwindow *sws.RootWidget
	mainwidget *sws.MainWidget
	nameL      *sws.LabelWidget
	name       *sws.InputWidget
	vps        *sws.CheckboxWidget
	vpsL       *sws.LabelWidget
	vpsNote    *sws.LabelWidget
	cpuOverL   *sws.LabelWidget
	cpuOver    *sws.DropdownWidget
	ramOverL   *sws.LabelWidget
	ramOver    *sws.DropdownWidget
	create     *sws.ButtonWidget
	cancel     *sws.ButtonWidget
}

func (self *PoolCreateWidget) Show() {
        self.rootwindow.AddChild(self.mainwidget)
        self.rootwindow.SetFocus(self.mainwidget)
}

func (self *PoolCreateWidget) Hide() {
        self.rootwindow.RemoveChild(self.mainwidget)
        children := self.rootwindow.GetChildren()
        if len(children) > 0 {
                self.rootwindow.SetFocus(children[0])
        }
}

func NewPoolCreateWidget(root *sws.RootWidget,inventory *Inventory) *PoolCreateWidget {
	mainwidget := sws.NewMainWidget(400, 220, " Create new pool ", false, false)
	mainwidget.Move(root.Width()/2-200,root.Height()/2-100)
	widget := &PoolCreateWidget {
		rootwindow:  root,
		mainwidget: mainwidget,
		nameL:   sws.NewLabelWidget(100,25,"Pool name:"),
		name:    sws.NewInputWidget(100,25,""),
		vps:     sws.NewCheckboxWidget(),
		vpsL:    sws.NewLabelWidget(100,25,"VPS pool (*)"),
		vpsNote: sws.NewLabelWidget(200,25,"(*) only if you can have VT processors"),
		cpuOverL:sws.NewLabelWidget(100,25,"Cpu overcommit"),
		cpuOver: sws.NewDropdownWidget(80,25,[]string{"0%","10%","20%","30%","40%","50%"}),
		ramOverL:sws.NewLabelWidget(100,25,"RAM overcommit"),
		ramOver: sws.NewDropdownWidget(80,25,[]string{"0%","10%","20%","30%","40%","50%"}),
		create:  sws.NewButtonWidget(100,25,"Create"),
		cancel:  sws.NewButtonWidget(100,25,"Cancel"),
	}
	
	widget.nameL.Move(40,20)
	mainwidget.AddChild(widget.nameL)
	
	widget.name.Move(150,20)
	mainwidget.AddChild(widget.name)
	
	widget.vpsL.Move(40,45)
	mainwidget.AddChild(widget.vpsL)
	
	widget.vps.Move(150,45)
	mainwidget.AddChild(widget.vps)
	widget.vps.SetClicked(func() {
		if widget.vps.Selected {
			mainwidget.AddChild(widget.cpuOverL)
			mainwidget.AddChild(widget.cpuOver)
			mainwidget.AddChild(widget.ramOverL)
			mainwidget.AddChild(widget.ramOver)
		} else {
			mainwidget.RemoveChild(widget.cpuOverL)
			mainwidget.RemoveChild(widget.cpuOver)
			mainwidget.RemoveChild(widget.ramOverL)
			mainwidget.RemoveChild(widget.ramOver)
		}
	})
	
	widget.vpsNote.Move(40,70)
	widget.vpsNote.SetFont(sws.LatoRegular12)
	mainwidget.AddChild(widget.vpsNote)
	
	widget.cpuOverL.Move(40,95)
	widget.cpuOver.Move(150,95)
	widget.ramOverL.Move(40,120)
	widget.ramOver.Move(150,120)
	
	widget.create.Move(150,150)
	mainwidget.AddChild(widget.create)
	
	widget.cancel.Move(260,150)
	mainwidget.AddChild(widget.cancel)
	
	widget.create.SetClicked(func() {
		if widget.vps.Selected {
			pool := &VpsServerPool{
				Name: widget.name.GetText(),
				pool: make(map[int32]*InventoryItem),
				cpuoverallocation: 1.0+0.1*float64(widget.cpuOver.ActiveChoice),
				ramoverallocation: 1.0+0.1*float64(widget.ramOver.ActiveChoice),
			}
			inventory.AddPool(pool)
		} else {
			pool := &HardwareServerPool{
				Name: widget.name.GetText(),
				pool: make(map[int32]*InventoryItem),
			}
			inventory.AddPool(pool)
		}
		widget.Hide()
	})
	widget.cancel.SetClicked(func() {
		widget.Hide()
	})
	mainwidget.SetCloseCallback(func() {
		widget.Hide()
 	})

	return widget
}

type PoolWidget struct {
	sws.CoreWidget
	inventory      *Inventory
	//scroll         *sws.ScrollWidget
	//vbox           *sws.VBoxWidget
	//globalcheckbox *sws.CheckboxWidget
	//selected       map[*UnallocatedServerLineWidget]bool
	buttonPanel    *sws.CoreWidget
	create         *sws.ButtonWidget
	poolCreateWidget *PoolCreateWidget
}

func NewPoolSub(root *sws.RootWidget, inventory *Inventory) *SubInventory {
	widget := &PoolWidget{
		CoreWidget:  *sws.NewCoreWidget(100, 100),
		inventory:   inventory,
		buttonPanel: sws.NewCoreWidget(200, 30),
		create:      sws.NewButtonWidget(100,25,"New Pool"),
		poolCreateWidget: NewPoolCreateWidget(root,inventory),
	}
	
	widget.create.Move(0,2)
	widget.create.SetClicked(func() {
		widget.poolCreateWidget.Show()
	})
	widget.buttonPanel.AddChild(widget.create)
	
	sub := &SubInventory{
		Icon:        "resources/icon-bucket.png",
		Title:       "Server pools",
		ButtonPanel: widget.buttonPanel,
		Widget:      widget,
	}
	
	return sub
}

func NewOfferSub(inventory *Inventory) *SubInventory {
	sub := &SubInventory{
		Icon:        "resources/icon-paper-bill.png",
		Title:       "Server offers",
		ButtonPanel: sws.NewCoreWidget(200, 30),
		Widget:      sws.NewScrollWidget(100, 100),
	}
	return sub
}
