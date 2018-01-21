package dctycoon

import (
	"fmt"

	"github.com/nzin/dctycoon/accounting"
	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/timer"
	"github.com/nzin/sws"
)

type DockWidget struct {
	sws.CoreWidget
	currentDay   *sws.LabelWidget
	timer        *timer.GameTimer
	pause        *sws.FlatButtonWidget
	play         *sws.FlatButtonWidget
	forward      *sws.FlatButtonWidget
	shop         *sws.FlatButtonWidget
	inventory    *sws.FlatButtonWidget
	electricity  *sws.FlatButtonWidget
	accounting   *sws.FlatButtonWidget
	stats        *sws.FlatButtonWidget
	ledgerButton *sws.FlatButtonWidget
	timerevent   *sws.TimerEvent
	ledger       *accounting.Ledger
}

func (self *DockWidget) SetStatsCallback(callback func()) {
	self.stats.SetClicked(callback)
}

func (self *DockWidget) SetShopCallback(callback func()) {
	self.shop.SetClicked(callback)
}

func (self *DockWidget) SetLedgerCallback(callback func()) {
	self.ledgerButton.SetClicked(callback)
	self.accounting.SetClicked(callback)
}

func (self *DockWidget) SetInventoryCallback(callback func()) {
	self.inventory.SetClicked(callback)
}

func (self *DockWidget) SetElectricityCallback(callback func()) {
	self.electricity.SetClicked(callback)
}

func (self *DockWidget) LedgerChange() {
	accounts := self.ledger.GetYearAccount(self.timer.CurrentTime.Year())
	self.ledgerButton.SetText(fmt.Sprintf("%.2f $", accounts["51"]))
}

// from GameTimerSubscriber interface
func (self *DockWidget) ChangeSpeed(speed int) {
	self.pause.SetColor(0xffdddddd)
	self.play.SetColor(0xffdddddd)
	self.forward.SetColor(0xffdddddd)
	switch speed {
	case SPEED_STOP:
		self.pause.SetColor(0xff8888ff)

	case SPEED_FORWARD:
		self.play.SetColor(0xff8888ff)

	case SPEED_FASTFORWARD:
		self.forward.SetColor(0xff8888ff)
	}
}

// from GameTimerSubscriber interface
func (self *DockWidget) NewDay(timer *timer.GameTimer) {
	today := fmt.Sprintf("%d %s %d", timer.CurrentTime.Day(), timer.CurrentTime.Month().String(), timer.CurrentTime.Year())
	self.currentDay.SetText(today)
}

func (self *DockWidget) helperAddButton(x, y int32, iconasset string) *sws.FlatButtonWidget {
	button := sws.NewFlatButtonWidget(25, 25, "")
	button.Move(x, y)
	if icon, err := global.LoadImageAsset(iconasset); err == nil {
		button.SetImageSurface(icon)
	}
	self.AddChild(button)
	return button
}

func NewDockWidget(root *sws.RootWidget, game *Game, gamemenu *MainGameMenu) *DockWidget {
	corewidget := sws.NewCoreWidget(150, 150)
	widget := &DockWidget{
		CoreWidget: *corewidget,
		currentDay: sws.NewLabelWidget(150, 25, "1 1 1990"),
		timer:      nil,
		ledger:     nil,
		timerevent: nil,
	}
	game.AddGameTimerSubscriber(widget)
	title := sws.NewLabelWidget(150, 25, "DC Tycoon")
	title.SetCentered(true)
	widget.AddChild(title)

	widget.currentDay.Move(5, 25)
	widget.AddChild(widget.currentDay)

	widget.pause = widget.helperAddButton(25, 50, "assets/ui/icon-pause-symbol.png")
	widget.pause.SetColor(0xff8888ff)
	widget.pause.SetClicked(func() {
		game.ChangeGameSpeed(SPEED_STOP)
	})

	widget.play = widget.helperAddButton(50, 50, "assets/ui/icon-arrowhead-pointing-to-the-right.png")
	widget.play.SetClicked(func() {
		game.ChangeGameSpeed(SPEED_FORWARD)
	})

	widget.forward = widget.helperAddButton(75, 50, "assets/ui/icon-forward-button.png")
	widget.forward.SetClicked(func() {
		game.ChangeGameSpeed(SPEED_FASTFORWARD)
	})

	widget.shop = widget.helperAddButton(25, 75, "assets/ui/icon-shopping-cart-black-shape.png")

	widget.inventory = widget.helperAddButton(50, 75, "assets/ui/icon-dropbox-logo.png")

	save := widget.helperAddButton(75, 75, "assets/ui/icon-blank-file.png")
	save.SetClicked(func() {
		gamemenu.ShowSave()
	})

	widget.stats = widget.helperAddButton(100, 75, "assets/ui/icon-graph.png")

	widget.electricity = widget.helperAddButton(25, 100, "assets/ui/icon-electricity.png")

	widget.accounting = widget.helperAddButton(50, 100, "assets/ui/icon-paper-bill.png")

	widget.ledgerButton = sws.NewFlatButtonWidget(150, 25, "")
	widget.ledgerButton.Move(0, 125)
	widget.AddChild(widget.ledgerButton)

	widget.Move(root.Width()-widget.Width(), 0)

	return widget
}

func (self *DockWidget) SetGame(timer *timer.GameTimer, ledger *accounting.Ledger) {
	self.timer = timer

	if self.ledger != nil {
		self.ledger.RemoveSubscriber(self)
	}
	self.ledger = ledger
	ledger.AddSubscriber(self)
}
