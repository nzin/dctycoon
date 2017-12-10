package dctycoon

// Pool / offer / contract management page.
//
// We have to:
// - list all categories (ac, generator, rack, servers)
// - servers: list non attributed servers (or be able to filter? attribute, type, subtype, ...)
//    -> a la gmail? (with some checkboxwidget)
// - see/build pools (Hardware / VPS)
// - see/build offers
// - see/build contract?
//

import (
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/sws"
)

type MainInventoryWidget struct {
	rootwindow  *sws.RootWidget
	mainwidget  *sws.MainWidget
	tabwidget   *sws.TabWidget
	serverpools *supplier.PoolManagementWidget
	offers      *supplier.OfferManagementWidget
}

func (self *MainInventoryWidget) Show() {
	self.rootwindow.AddChild(self.mainwidget)
	self.rootwindow.SetFocus(self.mainwidget)
	self.tabwidget.SelectTab(0)
}

func (self *MainInventoryWidget) Hide() {
	self.rootwindow.RemoveChild(self.mainwidget)
	children := self.rootwindow.GetChildren()
	if len(children) > 0 {
		self.rootwindow.SetFocus(children[0])
	}
}

// NewMainInventoryWidget presents the pool and offer management window
func NewMainInventoryWidget(root *sws.RootWidget) *MainInventoryWidget {
	mainwidget := sws.NewMainWidget(850, 400, " Inventory Management ", true, true)
	widget := &MainInventoryWidget{
		rootwindow:  root,
		mainwidget:  mainwidget,
		tabwidget:   sws.NewTabWidget(200, 200),
		serverpools: supplier.NewPoolManagementWidget(root, supplier.GlobalInventory),
		offers:      supplier.NewOfferManagementWidget(root, supplier.GlobalInventory),
	}
	widget.tabwidget.AddTab("server pools", widget.serverpools)
	widget.tabwidget.AddTab("offers", widget.offers)

	widget.mainwidget.SetInnerWidget(widget.tabwidget)

	widget.mainwidget.SetCloseCallback(func() {
		widget.Hide()
	})

	return widget
}