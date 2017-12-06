package supplier

/*
 * Offer management widget allow to create and check physical/vps offers
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
)

type OfferManagementLineWidget struct {
	sws.CoreWidget
	Checkbox *sws.CheckboxWidget
	desc     *sws.LabelWidget
	Vps      *sws.LabelWidget
	Nbcores  *sws.LabelWidget
	Ramsize  *sws.LabelWidget
	Disksize *sws.LabelWidget
	Vt       *sws.LabelWidget
	Price    *sws.LabelWidget

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

func NewOfferManagementLineWidget(offer *ServerOffer) *OfferManagementLineWidget {
	vps := "yes"
	if offer.Vps == false {
		vps = "no"
	}
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

	vt := "yes"
	if offer.Vt == false {
		vt = "no"
	}

	line := &OfferManagementLineWidget{
		CoreWidget: *sws.NewCoreWidget(550, 25),
		Checkbox:   sws.NewCheckboxWidget(),
		desc:       sws.NewLabelWidget(200, 25, offer.Name),
		Vps:        sws.NewLabelWidget(50, 25, vps),
		Nbcores:    sws.NewLabelWidget(50, 25, fmt.Sprint(offer.Nbcores)),
		Ramsize:    sws.NewLabelWidget(50, 25, ramSizeText),
		Disksize:   sws.NewLabelWidget(50, 25, diskText),
		Vt:         sws.NewLabelWidget(50, 25, vt),
		Price:      sws.NewLabelWidget(75, 25, fmt.Sprintf("%.0f $", offer.Price)),
	}
	line.AddChild(line.Checkbox)

	line.desc.Move(25, 0)
	line.AddChild(line.desc)

	line.Vps.Move(225, 0)
	line.AddChild(line.Vps)

	line.Nbcores.Move(275, 0)
	line.AddChild(line.Nbcores)

	line.Ramsize.Move(325, 0)
	line.AddChild(line.Ramsize)

	line.Disksize.Move(375, 0)
	line.AddChild(line.Disksize)

	line.Vt.Move(425, 0)
	line.AddChild(line.Vt)

	line.Price.Move(475, 0)
	line.AddChild(line.Price)

	return line
}

type OfferManagementWidget struct {
	sws.CoreWidget
	inventory       *Inventory
	root            *sws.RootWidget
	addoffer        *sws.ButtonWidget
	removeoffer     *sws.ButtonWidget
	vbox            *sws.VBoxWidget
	scrolllisting   *sws.ScrollWidget
	selectallButton *sws.CheckboxWidget
	newofferwindow  *OfferManagementNewOfferWidget
}

func (self *OfferManagementWidget) Resize(width, height int32) {
	self.CoreWidget.Resize(width, height)
	if height > 75 {
		if width > 550 {
			width = 550
		}
		if width < 20 {
			width = 20
		}
		self.scrolllisting.Resize(width, height-75)
	}
}

func NewOfferManagementWidget(root *sws.RootWidget, inventory *Inventory) *OfferManagementWidget {
	corewidget := sws.NewCoreWidget(800, 400)
	widget := &OfferManagementWidget{
		CoreWidget:      *corewidget,
		inventory:       inventory,
		root:            root,
		addoffer:        sws.NewButtonWidget(200, 25, "Create offer"),
		removeoffer:     sws.NewButtonWidget(200, 25, "Remove offer(s)"),
		vbox:            sws.NewVBoxWidget(550, 25),
		scrolllisting:   sws.NewScrollWidget(550, 25),
		selectallButton: sws.NewCheckboxWidget(),
	}
	widget.newofferwindow = NewOfferManagementNewOfferWidget(root, func(offer *ServerOffer) {
		offerline := NewOfferManagementLineWidget(offer)
		widget.vbox.AddChild(offerline)
		widget.scrolllisting.PostUpdate()
	})

	widget.scrolllisting.SetInnerWidget(widget.vbox)
	widget.scrolllisting.ShowHorizontalScrollbar(false)
	widget.scrolllisting.Move(0, 75)
	widget.AddChild(widget.scrolllisting)

	widget.addoffer.Move(5, 5)
	widget.addoffer.SetClicked(func() {
		widget.newofferwindow.Show(nil)
	})
	widget.AddChild(widget.addoffer)

	// scroll menu
	widget.selectallButton.Move(0, 50)
	widget.AddChild(widget.selectallButton)

	desc := sws.NewLabelWidget(200, 25, "Offer")
	desc.Move(25, 50)
	widget.AddChild(desc)

	return widget
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
		self.Price.SetText("100")
	} else {
		self.offer = offer
		self.Name.SetText(offer.Name)
		self.Vps.SetSelected(offer.Vps)
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
		self.Price.SetText(fmt.Sprintf("%f", offer.Price))
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
}

func NewOfferManagementNewOfferWidget(root *sws.RootWidget, savecallback func(*ServerOffer)) *OfferManagementNewOfferWidget {
	mainwidget := sws.NewMainWidget(400, 300, "Offer settings", false, false)
	mainwidget.Move(100, 100)

	widget := &OfferManagementNewOfferWidget{
		offer:        nil,
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
	}
	mainwidget.SetCloseCallback(func() {
		widget.Hide()
	})

	name := sws.NewLabelWidget(150, 25, "Offer name:")
	name.Move(10, 25)
	mainwidget.AddChild(name)
	widget.Name.Move(150, 25)
	mainwidget.AddChild(widget.Name)

	vps := sws.NewLabelWidget(150, 25, "VPS offer?")
	vps.Move(10, 50)
	mainwidget.AddChild(vps)
	widget.Vps.Move(150, 50)
	mainwidget.AddChild(widget.Vps)
	widget.Vps.SetClicked(func() {
		if widget.Vps.Selected {
			widget.Vt.SetSelected(false)
			mainwidget.RemoveChild(widget.Vt)
			mainwidget.RemoveChild(widget.vtdesc)
		} else {
			mainwidget.AddChild(widget.Vt)
			mainwidget.AddChild(widget.vtdesc)
		}
	})

	nbcores := sws.NewLabelWidget(150, 25, "Nb cores:")
	nbcores.Move(10, 75)
	mainwidget.AddChild(nbcores)
	widget.Nbcores.Move(150, 75)
	mainwidget.AddChild(widget.Nbcores)

	ramsize := sws.NewLabelWidget(150, 25, "Ram size (M,G):")
	ramsize.Move(10, 100)
	mainwidget.AddChild(ramsize)
	widget.Ramsize.Move(150, 100)
	mainwidget.AddChild(widget.Ramsize)

	disksize := sws.NewLabelWidget(150, 25, "Disk size (M,G,T):")
	disksize.Move(10, 120)
	mainwidget.AddChild(disksize)
	widget.Disksize.Move(150, 120)
	mainwidget.AddChild(widget.Disksize)

	widget.vtdesc.Move(10, 150)
	mainwidget.AddChild(widget.vtdesc)
	widget.Vt.Move(150, 150)
	mainwidget.AddChild(widget.Vt)

	price := sws.NewLabelWidget(150, 25, "Price/month:")
	price.Move(10, 175)
	mainwidget.AddChild(price)
	widget.Price.Move(150, 175)
	mainwidget.AddChild(widget.Price)

	widget.Save.Move(150, 225)
	mainwidget.AddChild(widget.Save)
	widget.Save.SetClicked(func() {
		widget.Hide()
		widget.save()
	})

	widget.Cancel.Move(250, 225)
	mainwidget.AddChild(widget.Cancel)
	widget.Cancel.SetClicked(func() {
		widget.Hide()
	})

	return widget
}

func (self *OfferManagementNewOfferWidget) save() {
	offer := self.offer
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
	self.savecallback(offer)
}
