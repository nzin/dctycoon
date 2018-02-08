package dctycoon

import (
	"github.com/nzin/dctycoon/accounting"
	"github.com/nzin/dctycoon/firewall"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/dctycoon/timer"
	"github.com/nzin/sws"
	log "github.com/sirupsen/logrus"
)

type GameUI struct {
	rootwindow        *sws.RootWidget
	dc                *DcWidget
	opening           *MainOpening
	supplierwidget    *MainSupplierWidget
	inventorywidget   *MainInventoryWidget
	statswidget       *MainStatsWidget
	accountingwidget  *accounting.MainAccountingWidget
	electricitywidget *MainElectricityWidget
	firewallwidget    *MainFirewallWidget
	dock              *DockWidget
	eventpublisher    *timer.EventPublisher
}

func NewGameUI(quit *bool, root *sws.RootWidget, game *Game) *GameUI {
	gamemenu := NewMainGameMenu(game, root, quit)
	gameui := &GameUI{
		rootwindow:        root,
		dc:                NewDcWidget(root.Width(), root.Height(), root),
		opening:           NewMainOpening(root.Width(), root.Height(), root, gamemenu),
		supplierwidget:    NewMainSupplierWidget(root),
		inventorywidget:   NewMainInventoryWidget(root),
		statswidget:       NewMainStatsWidget(root, game),
		accountingwidget:  accounting.NewMainAccountingWidget(root),
		electricitywidget: NewMainElectricityWidget(root),
		dock:              NewDockWidget(root, game, gamemenu),
		eventpublisher:    timer.NewEventPublisher(root),
		firewallwidget:    NewMainFirewallWidget(root),
	}

	gameui.dock.SetShopCallback(func() {
		gameui.supplierwidget.Show()
	})

	gameui.dock.SetStatsCallback(func() {
		gameui.statswidget.Show()
	})

	gameui.dock.SetFirewallCallback(func() {
		gameui.firewallwidget.Show()
	})

	gameui.dock.SetLedgerCallback(func() {
		gameui.accountingwidget.Show()
	})

	gameui.dock.SetInventoryCallback(func() {
		gameui.inventorywidget.Show()
	})

	gameui.dock.SetElectricityCallback(func() {
		gameui.electricitywidget.Show()
	})

	gameui.dc.SetInventoryManagementCallback(func() {
		gameui.inventorywidget.Show()
	})

	return gameui
}

//
// SetGame is used when creating a new game, or loading a game
// This re-init all ledger / inventory UI.
// Therefore you have to populate the ledger and inventory AFTER calling this method
func (self *GameUI) SetGame(globaltimer *timer.GameTimer, inventory *supplier.Inventory, ledger *accounting.Ledger, trends *supplier.Trend, location *supplier.LocationType, dcmap *DatacenterMap, firewall *firewall.Firewall) {
	log.Debug("GameUI::InitGame()")
	self.dc.SetGame(inventory, globaltimer.CurrentTime, dcmap)
	self.supplierwidget.SetGame(globaltimer, inventory, ledger, trends)
	self.inventorywidget.SetGame(inventory, globaltimer.CurrentTime)
	self.accountingwidget.SetGame(globaltimer, ledger)
	self.dock.SetGame(globaltimer, ledger)
	self.statswidget.SetGame()
	self.firewallwidget.SetGame(firewall)
	self.electricitywidget.SetGame(inventory, location)

	self.supplierwidget.Hide()
	self.inventorywidget.Hide()
	self.accountingwidget.Hide()
	self.statswidget.Hide()
	self.firewallwidget.Hide()
}

func (self *GameUI) ShowDC() {
	self.rootwindow.RemoveChild(self.opening)

	self.rootwindow.AddChild(self.dc)
	self.rootwindow.AddChild(self.dock)
	self.rootwindow.SetFocus(self.dc)

	self.supplierwidget.Hide()
	self.inventorywidget.Hide()
	self.accountingwidget.Hide()
	self.statswidget.Hide()
	self.firewallwidget.Hide()
	self.dc.HideUpgrade()
	self.dc.PostUpdate()
}

func (self *GameUI) ShowOpening() {
	self.rootwindow.RemoveChild(self.dc)
	self.rootwindow.RemoveChild(self.dock)

	self.rootwindow.AddChild(self.opening)
	self.rootwindow.SetFocus(self.opening)
}

func (self *GameUI) ShowUpgrade(game *Game, nextmap string) {
	self.dc.ShowUpgrade()
	self.dc.SetUpgradeCallback(func() {
		game.MigrateMap(nextmap)
	})
}
