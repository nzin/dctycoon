package dctycoon

import (
	"fmt"
	"time"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/sws"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	CHOICE_WIDTH = 230

	CATEGORY_SERVER_TOWER = 0
	CATEGORY_SERVER_RACK  = 1
	CATEGORY_RACK         = 2
	CATEGORY_AC           = 3
	CATEGORY_GENERATOR    = 4
)

type ElementDragPayload struct {
	item        *supplier.InventoryItem
	imageheight int32
}

func (self *ElementDragPayload) GetType() int32 {
	return global.DRAG_ELEMENT_PAYLOAD
}

func (self *ElementDragPayload) PayloadAccepted(bool) {
}

type HardwareChoiceItem struct {
	sws.LabelWidget
	item *supplier.InventoryItem
}

func (self *HardwareChoiceItem) UpdateSprite() {
	if self.item.Pool != nil {
		color := uint32(global.VPS_COLOR)
		if self.item.Pool.IsVps() == false {
			color = global.PHYSICAL_COLOR
		}
		self.LabelWidget.SetImageSurface(global.GlowImage("assets/ui/"+self.item.GetSprite()+"half.png", color))
	} else {
		if icon, err := global.LoadImageAsset("assets/ui/" + self.item.GetSprite() + "half.png"); err == nil {
			self.LabelWidget.SetImageSurface(icon)
		}
	}
}

func NewHardwareChoiceItem(item *supplier.InventoryItem) *HardwareChoiceItem {
	i := &HardwareChoiceItem{
		LabelWidget: *sws.NewLabelWidget(250, 50, item.ShortDescription(true)),
		item:        item,
	}
	i.AlignImageLeft(true)
	i.UpdateSprite()
	return i
}

//
// if we are selecting from an hardware subcategory, i.e. rack or towers
//
func (self *HardwareChoiceItem) MousePressDown(x, y int32, button uint8) {
	var payload sws.DragPayload

	// if we are dealing with a rack server
	if self.item.Serverconf.ConfType.NbU > 0 {
		payload = &ServerDragPayload{
			item: self.item,
		}
	} else { // tower
		elementPayload := &ElementDragPayload{
			item: self.item,
		}
		if img, err := global.LoadImageAsset("assets/ui/" + self.item.Serverconf.ConfType.ServerSprite + "0.png"); err == nil {
			elementPayload.imageheight = img.H
		}
		payload = elementPayload
	}

	var parent sws.Widget
	parent = self
	for parent != nil {
		x += parent.X()
		y += parent.Y()
		parent = parent.Parent()
	}
	if self.item.Pool != nil {
		color := uint32(global.VPS_COLOR)
		if self.item.Pool.IsVps() == false {
			color = global.PHYSICAL_COLOR
		}
		sws.NewDragEventSprite(x, y, global.GlowImage("assets/ui/"+self.item.Serverconf.ConfType.ServerSprite+"0.png", color), payload)
	} else {
		if icon, err := global.LoadImageAsset("assets/ui/" + self.item.Serverconf.ConfType.ServerSprite + "0.png"); err == nil {
			sws.NewDragEventSprite(x, y, icon, payload)
		}
	}
}

type HardwareChoiceCategory struct {
	sws.LabelWidget
	Category int32
	main     *HardwareChoice
	subpanel *sws.ScrollWidget
	vbox     *sws.VBoxWidget
	items    map[int32]*supplier.InventoryItem
}

func NewHardwareChoiceCategory(category int32, main *HardwareChoice) *HardwareChoiceCategory {
	c := &HardwareChoiceCategory{
		LabelWidget: *sws.NewLabelWidget(50, 75, "0x"),
		Category:    category,
		main:        main,
		items:       make(map[int32]*supplier.InventoryItem),
		subpanel:    sws.NewScrollWidget(250, 0),
		vbox:        sws.NewVBoxWidget(250, 0),
	}
	c.SetTextColor(sdl.Color{0xff, 0xff, 0xff, 0xff})
	c.SetColor(0)
	c.SetCentered(true)
	switch category {
	case CATEGORY_SERVER_TOWER:
		if icon, err := global.LoadImageAsset("assets/ui/icon.tower.png"); err == nil {
			c.SetImageSurface(icon)
		}
	case CATEGORY_SERVER_RACK:
		if icon, err := global.LoadImageAsset("assets/ui/icon.rackserver.png"); err == nil {
			c.SetImageSurface(icon)
		}
	case CATEGORY_RACK:
		if icon, err := global.LoadImageAsset("assets/ui/icon.rack.png"); err == nil {
			c.SetImageSurface(icon)
		}
	case CATEGORY_AC:
		if icon, err := global.LoadImageAsset("assets/ui/icon.ac.png"); err == nil {
			c.SetImageSurface(icon)
		}
	case CATEGORY_GENERATOR:
		if icon, err := global.LoadImageAsset("assets/ui/icon.generator.png"); err == nil {
			c.SetImageSurface(icon)
		}
	}
	c.subpanel.SetInnerWidget(c.vbox)
	c.subpanel.ShowHorizontalScrollbar(false)

	return c
}

func (self *HardwareChoiceCategory) addItem(item *supplier.InventoryItem) {
	self.items[item.Id] = item
	self.SetText(fmt.Sprintf("%dx", len(self.items)))

	if self.Category == CATEGORY_SERVER_TOWER || self.Category == CATEGORY_SERVER_RACK {
		self.vbox.AddChild(NewHardwareChoiceItem(item))
		if len(self.items) <= 10 {
			self.subpanel.Resize(250, int32(50*len(self.items)))
		}
	}
}

func (self *HardwareChoiceCategory) removeItem(item *supplier.InventoryItem) {
	delete(self.items, item.Id)
	self.SetText(fmt.Sprintf("%dx", len(self.items)))

	if self.Category == CATEGORY_SERVER_TOWER || self.Category == CATEGORY_SERVER_RACK {
		for _, i := range self.vbox.GetChildren() {
			hci := i.(*HardwareChoiceItem)
			if hci.item == item {
				self.vbox.RemoveChild(hci)
			}
		}
		if len(self.items) == 0 {
			self.main.closeItemPanel()
		} else if len(self.items) <= 10 {
			self.subpanel.Resize(250, int32(50*len(self.items)))
		}
	}
}

func (self *HardwareChoiceCategory) MousePressDown(x, y int32, button uint8) {
	if (self.Category == CATEGORY_RACK || self.Category == CATEGORY_AC || self.Category == CATEGORY_GENERATOR) &&
		button == sdl.BUTTON_LEFT && len(self.items) > 0 {

		// ugly way to get one of the item
		var item *supplier.InventoryItem
		for _, v := range self.items {
			item = v
			break
		}
		payload := &ElementDragPayload{
			item: item,
		}
		if img, err := global.LoadImageAsset("assets/ui/" + item.GetSprite() + "0.png"); err == nil {
			payload.imageheight = img.H
		}
		var parent sws.Widget
		parent = self
		for parent != nil {
			x += parent.X()
			y += parent.Y()
			parent = parent.Parent()
		}
		if img, err := global.LoadImageAsset("assets/ui/" + item.GetSprite() + "0.png"); err == nil {
			sws.NewDragEventSprite(x, y, img, payload)
		}
	}
}

func (self *HardwareChoiceCategory) MousePressUp(x, y int32, button uint8) {
	switch self.Category {
	case CATEGORY_SERVER_TOWER:
		if button == sdl.BUTTON_LEFT && len(self.items) > 0 {
			self.main.switchItemPanel(CATEGORY_SERVER_TOWER, self.subpanel)
		}
	case CATEGORY_SERVER_RACK:
		if button == sdl.BUTTON_LEFT && len(self.items) > 0 {
			self.main.switchItemPanel(CATEGORY_SERVER_RACK, self.subpanel)
		}
	}
}

type HardwareChoice struct {
	sws.CoreWidget
	categories           [5]*HardwareChoiceCategory
	inventory            *supplier.Inventory
	currentPanel         sws.Widget // extensions with servers list
	currentPanelCategory int32
}

func (self *HardwareChoice) addItem(item *supplier.InventoryItem) {
	var category int
	switch item.Typeitem {
	case supplier.PRODUCT_RACK:
		category = CATEGORY_RACK
	case supplier.PRODUCT_AC:
		category = CATEGORY_AC
	case supplier.PRODUCT_GENERATOR:
		category = CATEGORY_GENERATOR
	case supplier.PRODUCT_SERVER:
		if item.Serverconf.ConfType.NbU < 0 {
			category = CATEGORY_SERVER_TOWER
		} else {
			category = CATEGORY_SERVER_RACK
		}
	}
	if len(self.categories[category].items) == 0 {
		self.AddChild(self.categories[category])
		self.categories[category].Move(0, self.Height())
		self.categories[category].SetAlphaMod(0)
		self.Resize(50, self.Height()+75)

		var fadein = 0
		sws.TimerAddEvent(time.Now(), 100*time.Millisecond, func(evt *sws.TimerEvent) {
			myfadein := &fadein
			mycategory := category
			*myfadein++
			self.categories[mycategory].SetAlphaMod(uint8(255 * (*myfadein) / 5))
			self.categories[mycategory].PostUpdate()
			if *myfadein == 5 {
				evt.StopRepeat()
			}
		})

	}
	self.categories[category].addItem(item)
}

func (self *HardwareChoice) removeItem(item *supplier.InventoryItem, fadeeffect bool) {
	var category int
	switch item.Typeitem {
	case supplier.PRODUCT_RACK:
		category = CATEGORY_RACK
	case supplier.PRODUCT_AC:
		category = CATEGORY_AC
	case supplier.PRODUCT_GENERATOR:
		category = CATEGORY_GENERATOR
	case supplier.PRODUCT_SERVER:
		if item.Serverconf.ConfType.NbU < 0 {
			category = CATEGORY_SERVER_TOWER
		} else {
			category = CATEGORY_SERVER_RACK
		}
	}
	self.categories[category].removeItem(item)

	if len(self.categories[category].items) == 0 {
		if fadeeffect {
			var fadein = 0
			sws.TimerAddEvent(time.Now(), 100*time.Millisecond, func(evt *sws.TimerEvent) {
				myfadein := &fadein
				mycategory := category
				*myfadein++
				self.categories[mycategory].SetAlphaMod(uint8(255 * (5 - *myfadein) / 5))
				self.categories[mycategory].PostUpdate()
				if *myfadein == 5 {
					evt.StopRepeat()
					self.RemoveChild(self.categories[mycategory])
					self.Resize(50, self.Height()-75)

					for i, c := range self.GetChildren() {
						c.Move(0, int32(i)*75)
					}
					if self.currentPanel != nil {
						self.currentPanel.Move(50, self.categories[self.currentPanelCategory].Y())
					}
				}
			})
		} else {
			self.RemoveChild(self.categories[category])
			self.Resize(50, self.Height()-75)

			for i, c := range self.GetChildren() {
				c.Move(0, int32(i)*75)
			}
			if self.currentPanel != nil {
				self.currentPanel.Move(50, self.categories[self.currentPanelCategory].Y())
			}
		}
	}
}

func (self *HardwareChoice) ItemInTransit(*supplier.InventoryItem) {
}

func (self *HardwareChoice) ItemInStock(item *supplier.InventoryItem) {
	self.addItem(item)
}

func (self *HardwareChoice) ItemRemoveFromStock(item *supplier.InventoryItem) {
	self.removeItem(item, true)
}

func (self *HardwareChoice) ItemInstalled(item *supplier.InventoryItem) {
}

func (self *HardwareChoice) ItemUninstalled(item *supplier.InventoryItem) {
}

func (self *HardwareChoice) ItemChangedPool(item *supplier.InventoryItem) {
	if item.Typeitem == supplier.PRODUCT_SERVER {
		category := CATEGORY_SERVER_RACK
		if item.Serverconf.ConfType.NbU < 0 {
			category = CATEGORY_SERVER_TOWER
		}

		for _, l := range self.categories[category].vbox.GetChildren() {
			line := l.(*HardwareChoiceItem)
			line.UpdateSprite()
		}
	}
}

func (self *HardwareChoice) switchItemPanel(category int32, widget sws.Widget) {
	if self.currentPanel != nil {
		self.Parent().RemoveChild(self.currentPanel)
		for _, w := range self.categories {
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
	self.currentPanel.Move(50, self.categories[category].Y())

	self.Parent().AddChild(self.currentPanel)
	self.categories[category].SetColor(0xdddddddd)

	self.PostUpdate()
}

func (self *HardwareChoice) closeItemPanel() {
	self.Parent().RemoveChild(self.currentPanel)
	for _, w := range self.categories {
		w.SetColor(0)
	}
	self.currentPanel = nil
	self.currentPanelCategory = -1
	self.PostUpdate()
}

func NewHardwareChoice() *HardwareChoice {
	hc := &HardwareChoice{
		CoreWidget:           *sws.NewCoreWidget(50, 0),
		inventory:            nil,
		currentPanelCategory: -1,
	}
	hc.SetColor(0)
	hc.categories[CATEGORY_SERVER_TOWER] = NewHardwareChoiceCategory(CATEGORY_SERVER_TOWER, hc)
	hc.categories[CATEGORY_SERVER_RACK] = NewHardwareChoiceCategory(CATEGORY_SERVER_RACK, hc)
	hc.categories[CATEGORY_RACK] = NewHardwareChoiceCategory(CATEGORY_RACK, hc)
	hc.categories[CATEGORY_AC] = NewHardwareChoiceCategory(CATEGORY_AC, hc)
	hc.categories[CATEGORY_GENERATOR] = NewHardwareChoiceCategory(CATEGORY_GENERATOR, hc)

	return hc
}

func (self *HardwareChoice) SetGame(inventory *supplier.Inventory, currenttime time.Time) {
	if self.inventory != nil {
		self.inventory.RemoveInventorySubscriber(self)
		for _, i := range self.inventory.Items {
			self.removeItem(i, false)
		}
	}
	self.inventory = inventory
	inventory.AddInventorySubscriber(self)
	// for material not placed but in stock
	for _, item := range self.inventory.Items {
		if !item.IsPlaced() && item.HasArrived(currenttime) {
			self.ItemInStock(item)
		}
	}
}
