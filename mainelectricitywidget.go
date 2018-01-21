package dctycoon

import (
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/sws"
)

type MainElectricityWidget struct {
	rootwindow *sws.RootWidget
	mainwidget *sws.MainWidget
	inventory  *supplier.Inventory
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
	}

	widget.mainwidget.SetCloseCallback(func() {
		widget.Hide()
	})

	return widget
}

func (self *MainElectricityWidget) SetGame(inventory *supplier.Inventory) {
	self.inventory = inventory
}
