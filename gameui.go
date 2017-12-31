package dctycoon

import (
	"github.com/nzin/dctycoon/accounting"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/dctycoon/timer"
	"github.com/nzin/sws"
)

type GameUI struct {
	rootwindow       *sws.RootWidget
	dc               *DcWidget
	supplierwidget   *MainSupplierWidget
	inventorywidget  *MainInventoryWidget
	accountingwidget *accounting.MainAccountingWidget
	dock             *DockWidget
}

func NewGameUI(quit *bool, root *sws.RootWidget) *GameUI {
	gameui := &GameUI{
		rootwindow:       root,
		dc:               NewDcWidget(root.Width(), root.Height(), root),
		supplierwidget:   NewMainSupplierWidget(root),
		inventorywidget:  NewMainInventoryWidget(root),
		accountingwidget: accounting.NewMainAccountingWidget(root),
		dock:             NewDockWidget(root),
	}

	gameui.dock.SetShopCallback(func() {
		gameui.supplierwidget.Show()
	})

	gameui.dock.SetQuitCallback(func() {
		*quit = true
	})

	gameui.dock.SetLedgerCallback(func() {
		gameui.accountingwidget.Show()
	})

	gameui.dock.SetInventoryCallback(func() {
		gameui.inventorywidget.Show()
	})

	return gameui
}

//
// SetGame is used when creating a new game, or loadind a backup gamed
func (self *GameUI) SetGame(globaltimer *timer.GameTimer, inventory *supplier.Inventory, ledger *accounting.Ledger, trends *supplier.Trend) {
	self.dc.SetGame(inventory)
	self.supplierwidget.SetGame(globaltimer, inventory, ledger, trends)
	self.inventorywidget.SetGame(inventory)
	self.accountingwidget.SetGame(globaltimer, ledger)
	self.dock.SetGame(globaltimer, ledger)

	self.supplierwidget.Hide()
	self.inventorywidget.Hide()
	self.accountingwidget.Hide()
}

func (self *GameUI) LoadGame(v map[string]interface{}) {
	gamemap := v["map"].(map[string]interface{})
	self.dc.LoadMap(gamemap)
}

func (self *GameUI) ShowDC() {
	self.rootwindow.AddChild(self.dc)
	self.rootwindow.AddChild(self.dock)

	self.rootwindow.SetFocus(self.dc)
}
