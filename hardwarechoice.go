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

type HardwareChoiceCategory struct {
	sws.LabelWidget
	Category int32
	main *HardwareChoice
}

func NewHardwareChoiceCategory(category int32, main *HardwareChoice) * HardwareChoiceCategory{
	c := &HardwareChoiceCategory {
		LabelWidget: *sws.NewLabelWidget(50,75,"0x"),
		Category: category,
		main: main,
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
	return c
} 

type HardwareChoice struct {
	sws.CoreWidget
	categories [5]*HardwareChoiceCategory
	inventory *supplier.Inventory
}

func (self *HardwareChoice) addItem(item *supplier.InventoryItem) {
}

func (self *HardwareChoice) removeItem(item *supplier.InventoryItem) {
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

func (self *HardwareChoice) HasFocus(focus bool) {
	fmt.Println("HardwareChoice::HasFocus",focus)
}

func NewHardwareChoice(inventory *supplier.Inventory) *HardwareChoice{
	hc := &HardwareChoice {
		CoreWidget: *sws.NewCoreWidget(50,375),
		inventory: inventory,
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
