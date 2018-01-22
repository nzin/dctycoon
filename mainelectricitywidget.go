package dctycoon

import (
	"time"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/sws"
)

type MainElectricityWidget struct {
	rootwindow *sws.RootWidget
	mainwidget *sws.MainWidget
	inventory  *supplier.Inventory
	powerline1 *sws.DropdownWidget
	powerline2 *sws.DropdownWidget
	powerline3 *sws.DropdownWidget
}

func (self *MainElectricityWidget) Show() {
	self.rootwindow.AddChild(self.mainwidget)
	self.rootwindow.SetFocus(self.mainwidget)
}

func (self *MainElectricityWidget) Hide() {
	self.rootwindow.RemoveChild(self.mainwidget)
	children := self.rootwindow.GetChildren()
	if len(children) > 0 {
		self.rootwindow.SetFocus(children[0])
	}
}

// NewMainInventoryWidget presents the pool and offer management window
func NewMainElectricityWidget(root *sws.RootWidget) *MainElectricityWidget {
	mainwidget := sws.NewMainWidget(850, 400, " Power Utility ", true, true)
	mainwidget.Center(root)

	widget := &MainElectricityWidget{
		rootwindow: root,
		mainwidget: mainwidget,
		powerline1: sws.NewDropdownWidget(100, 25, []string{"none", "10kW", "100kW", "1MW", "10MW"}),
		powerline2: sws.NewDropdownWidget(100, 25, []string{"none", "10kW", "100kW", "1MW", "10MW"}),
		powerline3: sws.NewDropdownWidget(100, 25, []string{"none", "10kW", "100kW", "1MW", "10MW"}),
	}

	pilon := sws.NewLabelWidget(193, 213, "")
	if icon, err := global.LoadImageAsset("assets/ui/pilon2.png"); err == nil {
		pilon.SetImageSurface(icon)
	}
	widget.mainwidget.AddChild(pilon)

	widget.mainwidget.SetCloseCallback(func() {
		widget.Hide()
	})

	power1 := sws.NewLabelWidget(100, 25, "Main line:")
	power1.Move(200, 25)
	widget.mainwidget.AddChild(power1)
	widget.powerline1.Move(300, 25)
	widget.mainwidget.AddChild(widget.powerline1)
	widget.powerline1.SetCallbackValueChanged(func() {
		widget.inventory.SetPowerline(0, widget.powerline1.ActiveChoice)
	})

	power2 := sws.NewLabelWidget(100, 25, "Second line:")
	power2.Move(200, 50)
	widget.mainwidget.AddChild(power2)
	widget.powerline2.Move(300, 50)
	widget.mainwidget.AddChild(widget.powerline2)
	widget.powerline2.SetCallbackValueChanged(func() {
		widget.inventory.SetPowerline(1, widget.powerline2.ActiveChoice)
	})

	power3 := sws.NewLabelWidget(100, 25, "Third line:")
	power3.Move(200, 75)
	widget.mainwidget.AddChild(power3)
	widget.powerline3.Move(300, 75)
	widget.mainwidget.AddChild(widget.powerline3)
	widget.powerline3.SetCallbackValueChanged(func() {
		widget.inventory.SetPowerline(2, widget.powerline3.ActiveChoice)
	})

	return widget
}

func (self *MainElectricityWidget) PowerChange(time time.Time, consumed, generated, delivered float64) {
	powerlines := self.inventory.GetPowerlines()

	self.powerline1.SetActiveChoice(powerlines[0])
	self.powerline2.SetActiveChoice(powerlines[1])
	self.powerline3.SetActiveChoice(powerlines[2])
}

func (self *MainElectricityWidget) SetGame(inventory *supplier.Inventory) {
	self.inventory = inventory
	powerlines := inventory.GetPowerlines()
	self.powerline1.SetActiveChoice(powerlines[0])
	self.powerline2.SetActiveChoice(powerlines[1])
	self.powerline3.SetActiveChoice(powerlines[2])

	inventory.AddPowerStatSubscriber(self)
}
