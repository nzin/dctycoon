package supplier

/*
 * Offer management widget allow to create and check physical/vps OFFERs
 * see MainInventoryWiget
 *
 * - offer list
 * - offer details: nb servers, nb customers...
 */

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/sws"
	log "github.com/sirupsen/logrus"
)

type OfferManagementLineWidget struct {
	sws.CoreWidget
	Checkbox *sws.CheckboxWidget
	desc     *sws.LabelWidget
	Vps      *sws.CheckboxWidget
	Nbcores  *sws.LabelWidget
	Ramsize  *sws.LabelWidget
	Disksize *sws.LabelWidget
	Vt       *sws.CheckboxWidget
	Price    *sws.LabelWidget
	Fitting  *sws.LabelWidget

	offer *ServerOffer
}

func (self *OfferManagementLineWidget) AddChild(child sws.Widget) {
	self.CoreWidget.AddChild(child)
	child.SetParent(self)
}

func (self *OfferManagementLineWidget) MousePressDown(x, y int32, button uint8) {
	self.Checkbox.MousePressDown(1, 1, button)
}

func (self *OfferManagementLineWidget) MousePressUp(x, y int32, button uint8) {
	self.Checkbox.MousePressUp(1, 1, button)
}

// refresh the UI with the content of self.offer
func (self *OfferManagementLineWidget) Update() {
	self.desc.SetText(self.offer.Name)

	self.Vps.SetSelected(self.offer.Vps)

	self.Nbcores.SetText(fmt.Sprint(self.offer.Nbcores))

	ramSizeText := fmt.Sprintf("%d Mo", self.offer.Ramsize)
	if self.offer.Ramsize >= 2048 {
		ramSizeText = fmt.Sprintf("%d Go", self.offer.Ramsize/1024)
	}
	self.Ramsize.SetText(ramSizeText)

	diskText := fmt.Sprintf("%d Mo", self.offer.Disksize)
	if self.offer.Disksize > 4096 {
		diskText = fmt.Sprintf("%d Go", self.offer.Disksize/1024)
	}
	if self.offer.Disksize > 4*1024*1024 {
		diskText = fmt.Sprintf("%d To", self.offer.Disksize/(1024*1024))
	}
	self.Disksize.SetText(diskText)

	self.Vt.SetSelected(self.offer.Vt)

	self.Price.SetText(fmt.Sprintf("%.0f $", self.offer.Price))

	self.Fitting.SetText(fmt.Sprintf("%d", self.offer.Pool.HowManyFit(self.offer.Nbcores, self.offer.Ramsize, self.offer.Disksize, self.offer.Vt)))
}

func (self *OfferManagementLineWidget) InventoryItemAdd(*InventoryItem) {
	self.Update()
}
func (self *OfferManagementLineWidget) InventoryItemRemove(*InventoryItem) {
	self.Update()
}
func (self *OfferManagementLineWidget) InventoryItemAllocate(*InventoryItem) {
	self.Update()
}
func (self *OfferManagementLineWidget) InventoryItemRelease(*InventoryItem) {
	self.Update()
}

func NewOfferManagementLineWidget(offer *ServerOffer) *OfferManagementLineWidget {
	ramSizeText := fmt.Sprintf("%d Mo", offer.Ramsize)
	if offer.Ramsize >= 2048 {
		ramSizeText = fmt.Sprintf("%d Go", offer.Ramsize/1024)
	}

	diskText := fmt.Sprintf("%d Mo", offer.Disksize)
	if offer.Disksize > 4096 {
		diskText = fmt.Sprintf("%d Go", offer.Disksize/1024)
	}
	if offer.Disksize > 4*1024*1024 {
		diskText = fmt.Sprintf("%d To", offer.Disksize/(1024*1024))
	}

	line := &OfferManagementLineWidget{
		CoreWidget: *sws.NewCoreWidget(650, 25),
		Checkbox:   sws.NewCheckboxWidget(),
		desc:       sws.NewLabelWidget(200, 25, offer.Name),
		Vps:        sws.NewCheckboxWidget(),
		Nbcores:    sws.NewLabelWidget(50, 25, fmt.Sprint(offer.Nbcores)),
		Ramsize:    sws.NewLabelWidget(75, 25, ramSizeText),
		Disksize:   sws.NewLabelWidget(75, 25, diskText),
		Vt:         sws.NewCheckboxWidget(),
		Price:      sws.NewLabelWidget(75, 25, fmt.Sprintf("%.0f $", offer.Price)),
		Fitting:    sws.NewLabelWidget(50, 25, fmt.Sprintf("%d", offer.Pool.HowManyFit(offer.Nbcores, offer.Ramsize, offer.Disksize, offer.Vt))),
		offer:      offer,
	}
	line.Checkbox.SetColor(0)
	line.AddChild(line.Checkbox)

	line.desc.Move(25, 0)
	line.desc.SetColor(0)
	line.AddChild(line.desc)

	line.Vps.SetSelected(offer.Vps)
	line.Vps.SetDisabled(true)
	line.Vps.Move(225, 0)
	line.Vps.SetColor(0)
	line.AddChild(line.Vps)

	line.Nbcores.Move(275, 0)
	line.Nbcores.SetColor(0)
	line.AddChild(line.Nbcores)

	line.Ramsize.Move(325, 0)
	line.Ramsize.SetColor(0)
	line.AddChild(line.Ramsize)

	line.Disksize.Move(400, 0)
	line.Disksize.SetColor(0)
	line.AddChild(line.Disksize)

	line.Vt.SetSelected(offer.Vt)
	line.Vt.SetDisabled(true)
	line.Vt.Move(475, 0)
	line.Vt.SetColor(0)
	line.AddChild(line.Vt)

	line.Price.Move(525, 0)
	line.Price.SetColor(0)
	line.AddChild(line.Price)

	line.Fitting.Move(600, 0)
	line.Fitting.SetColor(0)
	line.AddChild(line.Fitting)

	line.SetColor(0xffffffff)

	return line
}

type OfferManagementWidget struct {
	sws.CoreWidget
	inventory      *Inventory
	root           *sws.RootWidget
	addoffer       *sws.ButtonWidget
	updateoffer    *sws.ButtonWidget
	removeoffer    *sws.ButtonWidget
	vbox           *sws.VBoxWidget
	scrolllisting  *sws.ScrollWidget
	highlight      *OfferManagementLineWidget
	newofferwindow *OfferManagementNewOfferWidget
}

func (self *OfferManagementWidget) Resize(width, height int32) {
	self.CoreWidget.Resize(width, height)
	if height > 75 {
		if width > 650 {
			width = 650
		}
		if width < 20 {
			width = 20
		}
		self.scrolllisting.Resize(width, height-75)
	}
}

//
// when we click on a offer line
func (self *OfferManagementWidget) HighlightLine(line *OfferManagementLineWidget, highlight bool) {
	if self.highlight != nil {
		self.highlight.Checkbox.SetSelected(false)
	}
	if highlight {
		line.Checkbox.SetSelected(true)
		self.highlight = line
		self.AddChild(self.updateoffer)
		self.AddChild(self.removeoffer)
	} else {
		self.RemoveChild(self.updateoffer)
		self.RemoveChild(self.removeoffer)
	}
}

func NewOfferManagementWidget(root *sws.RootWidget) *OfferManagementWidget {
	corewidget := sws.NewCoreWidget(800, 400)
	widget := &OfferManagementWidget{
		CoreWidget:    *corewidget,
		inventory:     nil,
		root:          root,
		addoffer:      sws.NewButtonWidget(150, 25, "Create offer"),
		updateoffer:   sws.NewButtonWidget(150, 25, "Update offer"),
		removeoffer:   sws.NewButtonWidget(150, 25, "Remove offer"),
		vbox:          sws.NewVBoxWidget(650, 0),
		scrolllisting: sws.NewScrollWidget(650, 0),
		highlight:     nil,
	}
	widget.newofferwindow = NewOfferManagementNewOfferWidget(root, func(offer *ServerOffer) {
		if widget.highlight == nil {
			offerline := NewOfferManagementLineWidget(offer)
			offerline.Checkbox.SetClicked(func() {
				widget.HighlightLine(offerline, offerline.Checkbox.Selected)
			})
			widget.vbox.AddChild(offerline)
			widget.AddChild(widget.scrolllisting)
			widget.scrolllisting.PostUpdate()
			// Inventory add offer
			offer.Pool.AddPoolSubscriber(offerline)
			widget.inventory.AddOffer(offer)
		} else {
			// Inventory update offer
			widget.inventory.UpdateOffer(offer)
			offer.Pool.AddPoolSubscriber(widget.highlight)
			widget.highlight.Update()
		}
	})

	na := global.NewNothingWidget(650, 25)
	na.Move(0, 75)
	widget.AddChild(na)

	widget.scrolllisting.SetInnerWidget(widget.vbox)
	widget.scrolllisting.ShowHorizontalScrollbar(false)
	widget.scrolllisting.Move(0, 75)
	//	widget.AddChild(widget.scrolllisting)

	widget.addoffer.Move(5, 5)
	widget.addoffer.SetClicked(func() {
		widget.HighlightLine(widget.highlight, false)
		widget.highlight = nil
		widget.newofferwindow.Show(nil)
	})
	widget.AddChild(widget.addoffer)

	widget.removeoffer.Move(160, 5)
	widget.removeoffer.SetClicked(func() {
		sws.ShowModalYesNo(root, "Remove Offer", "resources/icon-triangular-big.png", "Are you sure you want to remove this offer?", func() {
			widget.HighlightLine(widget.highlight, false)
			widget.vbox.RemoveChild(widget.highlight)
			//unsubscribe pool
			widget.highlight.offer.Pool.RemovePoolSubscriber(widget.highlight)
			// Inventory remove offer
			widget.inventory.RemoveOffer(widget.highlight.offer)
			if len(widget.vbox.GetChildren()) == 0 {
				widget.RemoveChild(widget.scrolllisting)
			}
		}, nil)
	})

	widget.updateoffer.Move(315, 5)
	widget.updateoffer.SetClicked(func() {
		if widget.highlight != nil {
			widget.highlight.offer.Pool.RemovePoolSubscriber(widget.highlight)
			widget.newofferwindow.Show(widget.highlight.offer)
		}
	})

	// scroll menu
	desc := sws.NewLabelWidget(200, 25, "Offer")
	desc.Move(25, 50)
	widget.AddChild(desc)

	vps := sws.NewLabelWidget(50, 25, "Vps")
	vps.Move(225, 50)
	widget.AddChild(vps)

	nbcores := sws.NewLabelWidget(50, 25, "#cores")
	nbcores.Move(275, 50)
	widget.AddChild(nbcores)

	ramsize := sws.NewLabelWidget(75, 25, "RAM")
	ramsize.Move(325, 50)
	widget.AddChild(ramsize)

	disksize := sws.NewLabelWidget(75, 25, "Disk")
	disksize.Move(400, 50)
	widget.AddChild(disksize)

	vt := sws.NewLabelWidget(25, 25, "VT")
	vt.Move(475, 50)
	widget.AddChild(vt)

	price := sws.NewLabelWidget(75, 25, "Price")
	price.Move(525, 50)
	widget.AddChild(price)

	fitting := sws.NewLabelWidget(50, 25, "#offers")
	fitting.Move(600, 50)
	widget.AddChild(fitting)

	return widget
}

func (self *OfferManagementWidget) SetGame(inventory *Inventory) {
	log.Debug("OfferManagementWidget::SetGame(", inventory, ")")
	self.newofferwindow.SetGame(inventory)
	if self.inventory != nil {
		for _, o := range self.vbox.GetChildren() {
			line := o.(*OfferManagementLineWidget)

			self.vbox.RemoveChild(line)
			//unsubscribe pool
			line.offer.Pool.RemovePoolSubscriber(line)
			// Inventory remove offer
		}
		self.RemoveChild(self.scrolllisting)
		self.highlight = nil
	}
	self.inventory = inventory

	for _, o := range self.inventory.offers {
		offerline := NewOfferManagementLineWidget(o)
		offerline.Checkbox.SetClicked(func() {
			self.HighlightLine(offerline, offerline.Checkbox.Selected)
		})
		self.vbox.AddChild(offerline)
		self.AddChild(self.scrolllisting)
		self.scrolllisting.PostUpdate()
		// Inventory add offer
		o.Pool.AddPoolSubscriber(offerline)
	}
}

func (self *OfferManagementNewOfferWidget) SetGame(inventory *Inventory) {
	self.inventory = inventory
}

func (self *OfferManagementNewOfferWidget) Show(offer *ServerOffer) {
	self.rootwindow.AddChild(self.mainwidget)
	self.rootwindow.SetFocus(self.mainwidget)
	if offer == nil {
		self.offer = &ServerOffer{}
		self.Name.SetText("")
		self.Vps.SetSelected(false)
		self.Nbcores.SetActiveChoice(0)
		self.Ramsize.SetText("32M")
		self.Disksize.SetText("150M")
		self.Vt.SetSelected(false)
		self.Vt.SetDisabled(false)
		self.Price.SetText("100")
	} else {
		self.offer = offer
		self.Name.SetText(offer.Name)
		self.Vps.SetSelected(offer.Vps)
		if offer.Vps {
			self.Vt.SetDisabled(true)
		} else {
			self.Vt.SetDisabled(false)
		}
		nbcores := offer.Nbcores
		if nbcores > 12 {
			nbcores = 12
		}
		self.Nbcores.SetActiveChoice(nbcores - 1)
		ramsize := fmt.Sprintf("%d M", offer.Ramsize)
		if offer.Ramsize > 4*1024 {
			ramsize = fmt.Sprintf("%d G", offer.Ramsize/1024)
		}
		self.Ramsize.SetText(ramsize)
		disksize := fmt.Sprintf("%d M", offer.Disksize)
		if offer.Disksize > 4*1024 {
			disksize = fmt.Sprintf("%d G", offer.Disksize/1024)
		}
		if offer.Disksize > 4*1024*1024 {
			disksize = fmt.Sprintf("%d T", offer.Disksize/(1024*1024))
		}
		self.Disksize.SetText(disksize)
		self.Vt.SetSelected(offer.Vt)
		self.Price.SetText(fmt.Sprintf("%.0f", offer.Price))
	}
	self.rootwindow.SetFocus(self.Name)
}

func (self *OfferManagementNewOfferWidget) Hide() {
	self.rootwindow.RemoveChild(self.mainwidget)
	children := self.rootwindow.GetChildren()
	if len(children) > 0 {
		self.rootwindow.SetFocus(children[len(children)-1])
	}
}

type OfferManagementNewOfferWidget struct {
	offer        *ServerOffer
	inventory    *Inventory
	rootwindow   *sws.RootWidget
	mainwidget   *sws.MainWidget
	Name         *sws.InputWidget
	Vps          *sws.CheckboxWidget
	Nbcores      *sws.DropdownWidget
	Ramsize      *sws.InputWidget
	Disksize     *sws.InputWidget
	Vt           *sws.CheckboxWidget
	vtdesc       *sws.LabelWidget
	Price        *sws.InputWidget
	Save         *sws.ButtonWidget
	Cancel       *sws.ButtonWidget
	savecallback func(*ServerOffer)
	HowManyFit   *sws.LabelWidget
}

func NewOfferManagementNewOfferWidget(root *sws.RootWidget, savecallback func(*ServerOffer)) *OfferManagementNewOfferWidget {
	mainwidget := sws.NewMainWidget(400, 350, "Offer settings", false, false)
	mainwidget.Move(100, 100)

	widget := &OfferManagementNewOfferWidget{
		offer: nil,
		//		inventory:    inventory,
		rootwindow:   root,
		mainwidget:   mainwidget,
		Name:         sws.NewInputWidget(200, 25, ""),
		Vps:          sws.NewCheckboxWidget(),
		Nbcores:      sws.NewDropdownWidget(100, 25, []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}),
		Ramsize:      sws.NewInputWidget(100, 25, "32M"),
		Disksize:     sws.NewInputWidget(100, 25, ""),
		vtdesc:       sws.NewLabelWidget(150, 25, "VT offer (non VPS):"),
		Vt:           sws.NewCheckboxWidget(),
		Price:        sws.NewInputWidget(100, 25, ""),
		Save:         sws.NewButtonWidget(75, 25, "Save"),
		Cancel:       sws.NewButtonWidget(75, 25, "Cancel"),
		savecallback: savecallback,
		HowManyFit:   sws.NewLabelWidget(100, 25, "0"),
	}
	mainwidget.SetCloseCallback(func() {
		widget.Hide()
	})

	name := sws.NewLabelWidget(150, 25, "Offer name:")
	name.Move(10, 25)
	mainwidget.AddChild(name)
	widget.Name.Move(160, 25)
	mainwidget.AddChild(widget.Name)

	vps := sws.NewLabelWidget(150, 25, "VPS offer?")
	vps.Move(10, 50)
	mainwidget.AddChild(vps)
	widget.Vps.Move(160, 50)
	mainwidget.AddChild(widget.Vps)
	widget.Vps.SetClicked(func() {
		if widget.Vps.Selected {
			widget.Vt.SetSelected(false)
			widget.Vt.SetDisabled(true)
		} else {
			widget.Vt.SetDisabled(false)
		}
	})

	nbcores := sws.NewLabelWidget(150, 25, "Nb cores:")
	nbcores.Move(10, 75)
	mainwidget.AddChild(nbcores)
	widget.Nbcores.Move(160, 75)
	mainwidget.AddChild(widget.Nbcores)

	ramsize := sws.NewLabelWidget(150, 25, "Ram size (M,G):")
	ramsize.Move(10, 100)
	mainwidget.AddChild(ramsize)
	widget.Ramsize.Move(160, 100)
	mainwidget.AddChild(widget.Ramsize)

	disksize := sws.NewLabelWidget(150, 25, "Disk size (M,G,T):")
	disksize.Move(10, 125)
	mainwidget.AddChild(disksize)
	widget.Disksize.Move(160, 125)
	mainwidget.AddChild(widget.Disksize)

	widget.vtdesc.Move(10, 150)
	mainwidget.AddChild(widget.vtdesc)
	widget.Vt.Move(160, 150)
	mainwidget.AddChild(widget.Vt)

	price := sws.NewLabelWidget(150, 25, "Price/month:")
	price.Move(10, 175)
	mainwidget.AddChild(price)
	widget.Price.Move(160, 175)
	mainwidget.AddChild(widget.Price)

	widget.Save.Move(160, 275)
	mainwidget.AddChild(widget.Save)
	widget.Save.SetClicked(func() {
		widget.Hide()
		widget.save()
	})

	widget.Cancel.Move(260, 275)
	mainwidget.AddChild(widget.Cancel)
	widget.Cancel.SetClicked(func() {
		widget.Hide()
	})

	// how many fit
	estimation := sws.NewLabelWidget(150, 25, "Nb offers fitting:")
	estimation.Move(10, 225)
	mainwidget.AddChild(estimation)

	widget.HowManyFit.Move(160, 225)
	mainwidget.AddChild(widget.HowManyFit)
	widget.Vps.SetCallbackValueChanged(func() {
		widget.UpdateHowManyFit()
	})
	widget.Nbcores.SetCallbackValueChanged(func() {
		widget.UpdateHowManyFit()
	})
	widget.Ramsize.SetCallbackValueChanged(func() {
		widget.UpdateHowManyFit()
	})
	widget.Disksize.SetCallbackValueChanged(func() {
		widget.UpdateHowManyFit()
	})
	widget.Vt.SetCallbackValueChanged(func() {
		widget.UpdateHowManyFit()
	})
	return widget
}

// will update the HowManyFit label with an estimation of how
// many offers we can fullfill
func (self *OfferManagementNewOfferWidget) UpdateHowManyFit() {
	vps := self.Vps.Selected
	nbcores := self.Nbcores.ActiveChoice + 1
	ramsize := global.ParseMega(self.Ramsize.GetText())
	disksize := global.ParseMega(self.Disksize.GetText())
	vt := self.Vt.Selected
	for _, pool := range self.inventory.GetPools() {
		if pool.IsVps() == vps {
			howmany := pool.HowManyFit(nbcores, ramsize, disksize, vt)
			self.HowManyFit.SetText(fmt.Sprintf("%d", howmany))
			return
		}
	}
}

// save to the self.offer with the correct values
func (self *OfferManagementNewOfferWidget) save() {
	offer := self.offer
	offer.Active = true
	offer.Name = self.Name.GetText()
	offer.Vps = self.Vps.Selected
	offer.Nbcores = self.Nbcores.ActiveChoice + 1
	offer.Ramsize = global.ParseMega(self.Ramsize.GetText())
	offer.Disksize = global.ParseMega(self.Disksize.GetText())
	offer.Vt = self.Vt.Selected
	re := regexp.MustCompile("([0-9]+)")
	match := re.FindStringSubmatch(self.Price.GetText())
	if price, err := strconv.ParseFloat(match[1], 64); err == nil {
		offer.Price = price
	} else {
		offer.Price = 0
	}
	// to be extended when we will have several pools
	for _, pool := range self.inventory.GetPools() {
		if pool.IsVps() == offer.Vps {
			offer.Pool = pool
		}
	}
	self.savecallback(offer)
}
