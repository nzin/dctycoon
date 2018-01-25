package dctycoon

import (
	"fmt"
	"time"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/sws"
)

type MainElectricityWidget struct {
	rootwindow  *sws.RootWidget
	mainwidget  *sws.MainWidget
	inventory   *supplier.Inventory
	location    *supplier.LocationType
	powerline1  *sws.DropdownWidget
	powerline2  *sws.DropdownWidget
	powerline3  *sws.DropdownWidget
	montlyprice *sws.LabelWidget
	usageprice  *sws.LabelWidget
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
	mainwidget := sws.NewMainWidget(650, 300, "Utility Management", true, true)
	mainwidget.Center(root)

	widget := &MainElectricityWidget{
		rootwindow:  root,
		mainwidget:  mainwidget,
		powerline1:  sws.NewDropdownWidget(100, 25, []string{"none", "10kW", "50kW", "200kW", "1MW"}),
		powerline2:  sws.NewDropdownWidget(100, 25, []string{"none", "10kW", "50kW", "200kW", "1MW"}),
		powerline3:  sws.NewDropdownWidget(100, 25, []string{"none", "10kW", "50kW", "200kW", "1MW"}),
		montlyprice: sws.NewLabelWidget(100, 25, "10 $"),
		usageprice:  sws.NewLabelWidget(100, 25, "0 $"),
	}

	pilon := sws.NewLabelWidget(193, 257, "")
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

	pricelabel := sws.NewLabelWidget(200, 25, "Fix montly fee:")
	pricelabel.Move(200, 125)
	widget.mainwidget.AddChild(pricelabel)
	widget.montlyprice.Move(400, 125)
	widget.mainwidget.AddChild(widget.montlyprice)

	usagelabel := sws.NewLabelWidget(200, 25, "Electricity usage (est.):")
	usagelabel.Move(200, 150)
	widget.mainwidget.AddChild(usagelabel)
	widget.usageprice.Move(400, 150)
	widget.mainwidget.AddChild(widget.usageprice)

	return widget
}

func (self *MainElectricityWidget) PowerChange(time time.Time, consumed, generated, delivered, cooler float64) {
	powerlines := self.inventory.GetPowerlines()

	self.powerline1.SetActiveChoice(powerlines[0])
	self.powerline2.SetActiveChoice(powerlines[1])
	self.powerline3.SetActiveChoice(powerlines[2])

	self.montlyprice.SetText(fmt.Sprintf("%.0f $", self.inventory.GetMonthlyPowerlinesPrice()))
	self.usageprice.SetText(fmt.Sprintf("%.0f $", consumed*24*30*self.location.Electricitycost/1000))
}

func (self *MainElectricityWidget) SetGame(inventory *supplier.Inventory, location *supplier.LocationType) {
	self.inventory = inventory
	self.location = location
	powerlines := inventory.GetPowerlines()
	self.powerline1.SetActiveChoice(powerlines[0])
	self.powerline2.SetActiveChoice(powerlines[1])
	self.powerline3.SetActiveChoice(powerlines[2])

	self.montlyprice.SetText(fmt.Sprintf("%.0f $", self.inventory.GetMonthlyPowerlinesPrice()))
	consumption, _, _ := inventory.GetGlobalPower()
	self.usageprice.SetText(fmt.Sprintf("%.0f $", consumption*24*30*location.Electricitycost/1000))

	inventory.AddPowerStatSubscriber(self)
}
