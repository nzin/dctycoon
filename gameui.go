package dctycoon

import (
	"fmt"

	"github.com/nzin/dctycoon/accounting"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/dctycoon/timer"
	"github.com/nzin/sws"
	log "github.com/sirupsen/logrus"
)

type GameUI struct {
	rootwindow       *sws.RootWidget
	dc               *DcWidget
	opening          *MainOpening
	supplierwidget   *MainSupplierWidget
	inventorywidget  *MainInventoryWidget
	statswidget      *MainStatsWidget
	accountingwidget *accounting.MainAccountingWidget
	dock             *DockWidget
	eventpublisher   *timer.EventPublisher
}

func NewGameUI(quit *bool, root *sws.RootWidget, game *Game) *GameUI {
	gamemenu := NewMainGameMenu(game, root, quit)
	gameui := &GameUI{
		rootwindow:       root,
		dc:               NewDcWidget(root.Width(), root.Height(), root),
		opening:          NewMainOpening(root.Width(), root.Height(), root, gamemenu),
		supplierwidget:   NewMainSupplierWidget(root),
		inventorywidget:  NewMainInventoryWidget(root),
		statswidget:      NewMainStatsWidget(root, game),
		accountingwidget: accounting.NewMainAccountingWidget(root),
		dock:             NewDockWidget(root, game, gamemenu),
		eventpublisher:   timer.NewEventPublisher(root),
	}

	gameui.dock.SetShopCallback(func() {
		gameui.supplierwidget.Show()
	})

	gameui.dock.SetStatsCallback(func() {
		gameui.statswidget.Show()
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
// InitGame is used when creating a new game
// This re-init all ledger / inventory UI.
// Therefore you have to populate the ledger and inventory AFTER calling this method
func (self *GameUI) InitGame(globaltimer *timer.GameTimer, inventory *supplier.Inventory, ledger *accounting.Ledger, trends *supplier.Trend) {
	log.Debug("GameUI::InitGame()")
	self.dc.SetGame(inventory, globaltimer.CurrentTime)
	self.supplierwidget.SetGame(globaltimer, inventory, ledger, trends)
	self.inventorywidget.SetGame(inventory, globaltimer.CurrentTime)
	self.accountingwidget.SetGame(globaltimer, ledger)
	self.dock.SetGame(globaltimer, ledger)
	self.statswidget.SetGame()

	self.supplierwidget.Hide()
	self.inventorywidget.Hide()
	self.accountingwidget.Hide()
	self.dc.InitMap("24_24_standard.json")
}

//
// LoadGame is used when loading a new game
// This re-init all ledger / inventory UI.
// Therefore you have to populate the ledger and inventory AFTER calling this method
func (self *GameUI) LoadGame(v map[string]interface{}, globaltimer *timer.GameTimer, inventory *supplier.Inventory, ledger *accounting.Ledger, trends *supplier.Trend) {
	log.Debug("GameUI::LoadGame()")
	self.dc.SetGame(inventory, globaltimer.CurrentTime)
	self.supplierwidget.SetGame(globaltimer, inventory, ledger, trends)
	self.inventorywidget.SetGame(inventory, globaltimer.CurrentTime)
	self.accountingwidget.SetGame(globaltimer, ledger)
	self.dock.SetGame(globaltimer, ledger)
	self.statswidget.LoadGame()

	self.supplierwidget.Hide()
	self.inventorywidget.Hide()
	self.accountingwidget.Hide()
	gamemap := v["map"].(map[string]interface{})
	self.dc.LoadMap(gamemap)
}

func (self *GameUI) SaveGame() string {
	return fmt.Sprintf(`"map": %s`, self.dc.SaveMap())
}

func (self *GameUI) ShowDC() {
	self.rootwindow.RemoveChild(self.opening)

	self.rootwindow.AddChild(self.dc)
	self.rootwindow.AddChild(self.dock)
	self.rootwindow.SetFocus(self.dc)
}

func (self *GameUI) ShowOpening() {
	self.rootwindow.RemoveChild(self.dc)
	self.rootwindow.RemoveChild(self.dock)

	self.rootwindow.AddChild(self.opening)
	self.rootwindow.SetFocus(self.opening)
}

func (self *GameUI) ShowHeatmap() {

}
