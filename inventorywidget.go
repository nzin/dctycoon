package dctycoon

import (
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/sws"
	//	"github.com/veandco/go-sdl2/sdl"
)

//
// We have to:
// - list all categories (ac, generator, rack, servers)
// - servers: list non attributed servers (or be able to filter? attribute, type, subtype, ...)
//    -> a la gmail? (faire des checkboxwidget)
// - see/build pools (Hardware / VPS)
// - see/build offers
// - see/build contract?
//
// tabwidget?
// upper: title, + buttons
//
type InventoryWidget struct {
	rootwindow *sws.RootWidget
	mainwidget *sws.MainWidget
	tabwidget  *sws.TabWidget
	servers    *supplier.ServerWidget
}

func (self *InventoryWidget) Show() {
	self.rootwindow.AddChild(self.mainwidget)
	self.rootwindow.SetFocus(self.mainwidget)
	self.tabwidget.SelectTab(0)
}

func (self *InventoryWidget) Hide() {
	self.rootwindow.RemoveChild(self.mainwidget)
	children := self.rootwindow.GetChildren()
	if len(children) > 0 {
		self.rootwindow.SetFocus(children[0])
	}
}

func NewInventoryWidget(root *sws.RootWidget) *InventoryWidget {
	mainwidget := sws.NewMainWidget(850, 400, " Inventory Management ", true, true)
	widget := &InventoryWidget{
		rootwindow: root,
		mainwidget: mainwidget,
		tabwidget:  sws.NewTabWidget(200, 200),
		servers:    supplier.NewServerWidget(root, supplier.GlobalInventory),
	}
	widget.tabwidget.AddTab("servers", widget.servers)

	widget.mainwidget.SetInnerWidget(widget.tabwidget)

	widget.mainwidget.SetCloseCallback(func() {
		widget.Hide()
	})

	return widget
}
