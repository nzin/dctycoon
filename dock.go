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
	quit         *sws.FlatButtonWidget
	ledgerButton *sws.FlatButtonWidget
	timerevent   *sws.TimerEvent
	ledger       *accounting.Ledger
}

func (self *DockWidget) SetQuitCallback(callback func()) {
	self.quit.SetClicked(callback)
}

func (self *DockWidget) SetShopCallback(callback func()) {
	self.shop.SetClicked(callback)
}

func (self *DockWidget) SetLedgerCallback(callback func()) {
	self.ledgerButton.SetClicked(callback)
}

func (self *DockWidget) SetInventoryCallback(callback func()) {
	self.inventory.SetClicked(callback)
}

func (self *DockWidget) LedgerChange() {
	accounts := self.ledger.GetYearAccount(self.timer.CurrentTime.Year())
	self.ledgerButton.SetText(fmt.Sprintf("%.2f $", accounts["51"]))
}

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

func (self *DockWidget) NewDay(timer *timer.GameTimer) {
	today := fmt.Sprintf("%d %s %d", timer.CurrentTime.Day(), timer.CurrentTime.Month().String(), timer.CurrentTime.Year())
	self.currentDay.SetText(today)
}

func NewDockWidget(root *sws.RootWidget, game *Game, gamemenu *MainGameMenu) *DockWidget {
	corewidget := sws.NewCoreWidget(150, 125)
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

	widget.pause = sws.NewFlatButtonWidget(25, 25, "")
	widget.pause.SetColor(0xff8888ff)
	widget.pause.Move(25, 50)
	if icon, err := global.LoadImageAsset("assets/ui/icon-pause-symbol.png"); err == nil {
		widget.pause.SetImageSurface(icon)
	}
	widget.pause.SetClicked(func() {
		game.ChangeGameSpeed(SPEED_STOP)
	})
	widget.AddChild(widget.pause)

	widget.play = sws.NewFlatButtonWidget(25, 25, "")
	widget.play.Move(50, 50)
	if icon, err := global.LoadImageAsset("assets/ui/icon-arrowhead-pointing-to-the-right.png"); err == nil {
		widget.play.SetImageSurface(icon)
	}
	widget.play.SetClicked(func() {
		game.ChangeGameSpeed(SPEED_FORWARD)
	})
	widget.AddChild(widget.play)

	widget.forward = sws.NewFlatButtonWidget(25, 25, "")
	widget.forward.Move(75, 50)
	if icon, err := global.LoadImageAsset("assets/ui/icon-forward-button.png"); err == nil {
		widget.forward.SetImageSurface(icon)
	}
	widget.forward.SetClicked(func() {
		game.ChangeGameSpeed(SPEED_FASTFORWARD)
	})
	widget.AddChild(widget.forward)

	widget.shop = sws.NewFlatButtonWidget(25, 25, "")
	widget.shop.Move(25, 75)
	if icon, err := global.LoadImageAsset("assets/ui/icon-shopping-cart-black-shape.png"); err == nil {
		widget.shop.SetImageSurface(icon)
	}
	widget.AddChild(widget.shop)

	widget.inventory = sws.NewFlatButtonWidget(25, 25, "")
	widget.inventory.Move(50, 75)
	if icon, err := global.LoadImageAsset("assets/ui/icon-dropbox-logo.png"); err == nil {
		widget.inventory.SetImageSurface(icon)
	}
	widget.AddChild(widget.inventory)

	save := sws.NewFlatButtonWidget(25, 25, "")
	save.Move(75, 75)
	if icon, err := global.LoadImageAsset("assets/ui/icon-blank-file.png"); err == nil {
		save.SetImageSurface(icon)
	}
	widget.AddChild(save)
	save.SetClicked(func() {
		gamemenu.ShowSave()
	})

	widget.quit = sws.NewFlatButtonWidget(25, 25, "")
	widget.quit.Move(100, 75)
	if icon, err := global.LoadImageAsset("assets/ui/icon-power-button-off.png"); err == nil {
		widget.quit.SetImageSurface(icon)
	}
	widget.AddChild(widget.quit)

	widget.ledgerButton = sws.NewFlatButtonWidget(150, 25, "")
	widget.ledgerButton.Move(0, 100)
	widget.AddChild(widget.ledgerButton)

	widget.Move(root.Width()-widget.Width(), 0)

	return widget
}

func (self *DockWidget) SetGame(timer *timer.GameTimer, ledger *accounting.Ledger) {
	self.timer = timer
	//	today := fmt.Sprintf("%d %s %d", timer.CurrentTime.Day(), timer.CurrentTime.Month().String(), timer.CurrentTime.Year())
	//	self.currentDay.SetText(today)

	if self.ledger != nil {
		self.ledger.RemoveSubscriber(self)
	}
	self.ledger = ledger
	ledger.AddSubscriber(self)
}
