package dctycoon

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/nzin/sws"
	"github.com/nzin/dctycoon/supplier"
)

const (
	CHOICE_WIDTH = 230
	
	CATEGORY_SERVER_TOWER = 0
	CATEGORY_SERVER_RACK = 1
	CATEGORY_RACK = 2
	CATEGORY_AC = 3
	CATEGORY_GENERATOR = 4
)

type HardwareChoiceItem struct {
	sws.LabelWidget
	item *supplier.InventoryItem
}

func NewHardwareChoiceItem(item *supplier.InventoryItem) *HardwareChoiceItem{
	i := &HardwareChoiceItem {
		LabelWidget: *sws.NewLabelWidget(200,50,item.UltraShortDescription()),
		item: item,
	}
	i.AlignImageLeft(true)
	i.SetImage("resources/icon."+item.GetSprite()+".png")
	return i
}

type HardwareChoiceCategory struct {
	sws.LabelWidget
	Category int32
	main *HardwareChoice
	subpanel *sws.ScrollWidget
	vbox *sws.VBoxWidget
	items map[int32]*supplier.InventoryItem
}

func NewHardwareChoiceCategory(category int32, main *HardwareChoice) * HardwareChoiceCategory{
	c := &HardwareChoiceCategory {
		LabelWidget: *sws.NewLabelWidget(50,75,"0x"),
		Category: category,
		main: main,
		items: make(map[int32]*supplier.InventoryItem),
		subpanel: sws.NewScrollWidget(200,0),
		vbox: sws.NewVBoxWidget(200,0),
	}
	c.SetTextColor(sdl.Color{0xff,0xff,0xff,0xff})
	c.SetColor(0)
	c.SetCentered(true)
	switch category {
		case CATEGORY_SERVER_TOWER:
			c.SetImage("resources/icon.tower.png")
		case CATEGORY_SERVER_RACK:
			c.SetImage("resources/icon.rackserver.png")
		case CATEGORY_RACK:
			c.SetImage("resources/icon.rack.png")
		case CATEGORY_AC:
			c.SetImage("resources/icon.ac.png")
		case CATEGORY_GENERATOR:
			c.SetImage("resources/icon.generator.png")
	}
	c.subpanel.SetInnerWidget(c.vbox)
	c.subpanel.ShowHorizontalScrollbar(false)
	
	return c
}

func (self *HardwareChoiceCategory) addItem(item *supplier.InventoryItem) {
	self.items[item.Id]=item
	self.SetText(fmt.Sprintf("%dx",len(self.items)))
	
	if (self.Category == CATEGORY_SERVER_TOWER || self.Category == CATEGORY_SERVER_RACK) {
		if (len(self.items)<=10) {
			self.subpanel.Resize(200,int32(50*len(self.items)))
			self.vbox.AddChild(NewHardwareChoiceItem(item))
		}
	}
}

func (self *HardwareChoiceCategory) removeItem(item *supplier.InventoryItem) {
	delete(self.items, item.Id)
	self.SetText(fmt.Sprintf("%dx",len(self.items)))

	if (self.Category == CATEGORY_SERVER_TOWER || self.Category == CATEGORY_SERVER_RACK) {
		for _,i := range(self.vbox.GetChildren()) {
			hci := i.(*HardwareChoiceItem)
			if hci.item == item {
				self.vbox.RemoveChild(hci)
			}
		}
		if (len(self.items)==0) {
			self.main.closeItemPanel()
		} else if (len(self.items)<=10) {
			self.subpanel.Resize(200,int32(50*len(self.items)))
		}
	}
}

func (self *HardwareChoiceCategory) MousePressDown(x, y int32, button uint8) {
}

func (self *HardwareChoiceCategory) MousePressUp(x, y int32, button uint8) {
	switch self.Category {
		case CATEGORY_SERVER_TOWER:
			if button == sdl.BUTTON_LEFT && len(self.items)>0 {
				self.main.switchItemPanel(CATEGORY_SERVER_TOWER,self.subpanel)
			}
		case CATEGORY_SERVER_RACK:
			if button == sdl.BUTTON_LEFT && len(self.items)>0 {
				self.main.switchItemPanel(CATEGORY_SERVER_RACK,self.subpanel)
			}
	}
}

type HardwareChoice struct {
	sws.CoreWidget
	categories [5]*HardwareChoiceCategory
	inventory *supplier.Inventory
	currentPanel sws.Widget
	currentPanelCategory int32
}

func (self *HardwareChoice) addItem(item *supplier.InventoryItem) {
	switch item.Typeitem {
		case supplier.PRODUCT_RACK:
			self.categories[CATEGORY_RACK].addItem(item)
		case supplier.PRODUCT_AC:
			self.categories[CATEGORY_AC].addItem(item)
		case supplier.PRODUCT_GENERATOR:
			self.categories[CATEGORY_GENERATOR].addItem(item)
		case supplier.PRODUCT_SERVER:
			if item.Serverconf.ConfType.NbU < 0 {
				self.categories[CATEGORY_SERVER_TOWER].addItem(item)
			} else {
				self.categories[CATEGORY_SERVER_RACK].addItem(item)
			}
	}
}

func (self *HardwareChoice) removeItem(item *supplier.InventoryItem) {
	switch item.Typeitem {
		case supplier.PRODUCT_RACK:
			self.categories[CATEGORY_RACK].removeItem(item)
		case supplier.PRODUCT_AC:
			self.categories[CATEGORY_AC].removeItem(item)
		case supplier.PRODUCT_GENERATOR:
			self.categories[CATEGORY_GENERATOR].removeItem(item)
		case supplier.PRODUCT_SERVER:
			if item.Serverconf.ConfType.NbU < 0 {
				self.categories[CATEGORY_SERVER_TOWER].removeItem(item)
			} else {
				self.categories[CATEGORY_SERVER_RACK].removeItem(item)
			}
	}
}

func (self *HardwareChoice) ItemInTransit(*supplier.InventoryItem) {
}

func (self *HardwareChoice) ItemInStock(item *supplier.InventoryItem) {
	self.addItem(item)
}

func (self *HardwareChoice) ItemRemoveFromStock(item *supplier.InventoryItem) {
	self.removeItem(item)
}

func (self *HardwareChoice) ItemInstalled(item *supplier.InventoryItem) {
	self.removeItem(item)
}

func (self *HardwareChoice) ItemUninstalled(item *supplier.InventoryItem) {
	self.addItem(item)
}

func (self *HardwareChoice) switchItemPanel(category int32, widget sws.Widget) {
	if self.currentPanel != nil {
		self.RemoveChild(self.currentPanel)
		for _,w := range self.categories {
			w.SetColor(0)
		}
		
		if self.currentPanelCategory == category {
			self.currentPanel = nil
			self.currentPanelCategory = -1
			return
		}
	}
	self.currentPanelCategory = category
	self.currentPanel = widget
	self.currentPanel.Move(50,75*category)
	
	self.AddChild(self.currentPanel)
	self.categories[category].SetColor(0xdddddddd)

	height := self.Height()
	if (self.currentPanel.Y()+self.currentPanel.Height()) > height {
		height = self.currentPanel.Y()+self.currentPanel.Height()
	}
	self.Resize(50+self.currentPanel.Width(),height)
	sws.PostUpdate()
}

func (self *HardwareChoice) closeItemPanel() {
	self.RemoveChild(self.currentPanel)
	for _,w := range self.categories {
		w.SetColor(0)
	}
	self.Resize(50,375)
	self.currentPanel = nil
	self.currentPanelCategory = -1
	sws.PostUpdate()
}

func NewHardwareChoice(inventory *supplier.Inventory) *HardwareChoice{
	hc := &HardwareChoice {
		CoreWidget: *sws.NewCoreWidget(50,375),
		inventory: inventory,
		currentPanelCategory: -1,
	}
	hc.SetColor(0)
	hc.categories[CATEGORY_SERVER_TOWER] = NewHardwareChoiceCategory(CATEGORY_SERVER_TOWER,hc)
	hc.AddChild(hc.categories[CATEGORY_SERVER_TOWER])
	hc.categories[CATEGORY_SERVER_RACK] = NewHardwareChoiceCategory(CATEGORY_SERVER_RACK,hc)
	hc.categories[CATEGORY_SERVER_RACK].Move(0,75)
	hc.AddChild(hc.categories[CATEGORY_SERVER_RACK])
	hc.categories[CATEGORY_RACK] = NewHardwareChoiceCategory(CATEGORY_RACK,hc)
	hc.categories[CATEGORY_RACK].Move(0,150)
	hc.AddChild(hc.categories[CATEGORY_RACK])
	hc.categories[CATEGORY_AC] = NewHardwareChoiceCategory(CATEGORY_AC,hc)
	hc.categories[CATEGORY_AC].Move(0,225)
	hc.AddChild(hc.categories[CATEGORY_AC])
	hc.categories[CATEGORY_GENERATOR] = NewHardwareChoiceCategory(CATEGORY_GENERATOR,hc)
	hc.categories[CATEGORY_GENERATOR].Move(0,300)
	hc.AddChild(hc.categories[CATEGORY_GENERATOR])

	inventory.AddSubscriber(hc)
	
	return hc
}
