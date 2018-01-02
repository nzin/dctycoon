package dctycoon

import (
	"fmt"
	"time"

	"github.com/nzin/dctycoon/accounting"
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

func NewDockWidget(root *sws.RootWidget, gamemenu *MainGameMenu) *DockWidget {
	corewidget := sws.NewCoreWidget(150, 125)
	widget := &DockWidget{
		CoreWidget: *corewidget,
		currentDay: sws.NewLabelWidget(150, 25, "1 1 1990"),
		timer:      nil,
		ledger:     nil,
		timerevent: nil,
	}
	title := sws.NewLabelWidget(150, 25, "DC Tycoon")
	title.SetCentered(true)
	widget.AddChild(title)

	widget.currentDay.Move(5, 25)
	widget.AddChild(widget.currentDay)

	widget.pause = sws.NewFlatButtonWidget(25, 25, "")
	widget.pause.SetColor(0xff8888ff)
	widget.pause.Move(25, 50)
	widget.pause.SetImage("resources/icon-pause-symbol.png")
	widget.pause.SetClicked(func() {
		widget.pause.SetColor(0xff8888ff)
		widget.play.SetColor(0xffdddddd)
		widget.forward.SetColor(0xffdddddd)
		if widget.timerevent != nil {
			widget.timerevent.StopRepeat()
		}
	})
	widget.AddChild(widget.pause)

	widget.play = sws.NewFlatButtonWidget(25, 25, "")
	widget.play.Move(50, 50)
	widget.play.SetImage("resources/icon-arrowhead-pointing-to-the-right.png")
	widget.play.SetClicked(func() {
		widget.pause.SetColor(0xffdddddd)
		widget.play.SetColor(0xff8888ff)
		widget.forward.SetColor(0xffdddddd)
		if widget.timerevent != nil {
			widget.timerevent.StopRepeat()
		}
		widget.timerevent = sws.TimerAddEvent(time.Now().Add(2*time.Second), 2*time.Second, func(evt *sws.TimerEvent) {
			widget.timer.TimerClock()
			today := fmt.Sprintf("%d %s %d", widget.timer.CurrentTime.Day(), widget.timer.CurrentTime.Month().String(), widget.timer.CurrentTime.Year())
			widget.currentDay.SetText(today)
		})
	})
	widget.AddChild(widget.play)

	widget.forward = sws.NewFlatButtonWidget(25, 25, "")
	widget.forward.Move(75, 50)
	widget.forward.SetImage("resources/icon-forward-button.png")
	widget.forward.SetClicked(func() {
		widget.pause.SetColor(0xffdddddd)
		widget.play.SetColor(0xffdddddd)
		widget.forward.SetColor(0xff8888ff)
		if widget.timerevent != nil {
			widget.timerevent.StopRepeat()
		}
		widget.timerevent = sws.TimerAddEvent(time.Now().Add(time.Second/2), time.Second/2, func(evt *sws.TimerEvent) {
			widget.timer.TimerClock()
			today := fmt.Sprintf("%d %s %d", widget.timer.CurrentTime.Day(), widget.timer.CurrentTime.Month().String(), widget.timer.CurrentTime.Year())
			widget.currentDay.SetText(today)
		})
	})
	widget.AddChild(widget.forward)

	widget.shop = sws.NewFlatButtonWidget(25, 25, "")
	widget.shop.Move(25, 75)
	widget.shop.SetImage("resources/icon-shopping-cart-black-shape.png")
	widget.AddChild(widget.shop)

	widget.inventory = sws.NewFlatButtonWidget(25, 25, "")
	widget.inventory.Move(50, 75)
	widget.inventory.SetImage("resources/icon-dropbox-logo.png")
	widget.AddChild(widget.inventory)

	save := sws.NewFlatButtonWidget(25, 25, "")
	save.Move(75, 75)
	save.SetImage("resources/icon-blank-file.png")
	widget.AddChild(save)
	save.SetClicked(func() {
		gamemenu.ShowSave()
		root.AddChild(gamemenu)
		gamemenu.SetCancelCallback(func() {
			root.RemoveChild(gamemenu)
		})
		if widget.timerevent != nil {
			widget.timerevent.StopRepeat()
		}
	})

	widget.quit = sws.NewFlatButtonWidget(25, 25, "")
	widget.quit.Move(100, 75)
	widget.quit.SetImage("resources/icon-power-button-off.png")
	widget.AddChild(widget.quit)

	widget.ledgerButton = sws.NewFlatButtonWidget(150, 25, "")
	widget.ledgerButton.Move(0, 100)
	widget.AddChild(widget.ledgerButton)

	widget.Move(root.Width()-widget.Width(), 0)

	return widget
}

func (self *DockWidget) SetGame(timer *timer.GameTimer, ledger *accounting.Ledger) {
	self.timer = timer
	today := fmt.Sprintf("%d %s %d", timer.CurrentTime.Day(), timer.CurrentTime.Month().String(), timer.CurrentTime.Year())
	self.currentDay.SetText(today)

	if self.ledger != nil {
		self.ledger.RemoveSubscriber(self)
	}
	self.ledger = ledger
	ledger.AddSubscriber(self)
}
