package dctycoon

import(
	"github.com/nzin/sws"
	"github.com/nzin/dctycoon/timer"
	"github.com/nzin/dctycoon/accounting"
	"time"
	"fmt"
)

type DockWidget struct {
	sws.SWS_CoreWidget
	currentDay    *sws.SWS_Label
	timer         *timer.GameTimer
	pause         *sws.SWS_FlatButtonWidget
	play          *sws.SWS_FlatButtonWidget
	forward       *sws.SWS_FlatButtonWidget
	shop          *sws.SWS_FlatButtonWidget
	quit          *sws.SWS_FlatButtonWidget
	ledger        *sws.SWS_FlatButtonWidget
	timerevent    *sws.TimerEvent
}

func (self *DockWidget) SetQuitCallback(callback func()) {
	self.quit.SetClicked(callback)
}

func (self *DockWidget) SetShopCallback(callback func()) {
	self.shop.SetClicked(callback)
}

func (self *DockWidget) SetLedgerCallback(callback func()) {
	self.ledger.SetClicked(callback)
}

func (self *DockWidget) LedgerChange(ledger *accounting.Ledger) {
	accounts:=ledger.GetYearAccount(self.timer.CurrentTime.Year())
	self.ledger.SetText(fmt.Sprintf("%.2f $",accounts["51"]))
}

func CreateDockWidget(timer *timer.GameTimer) *DockWidget {
	corewidget := sws.CreateCoreWidget(150, 125)
	today:=fmt.Sprintf("%d %s %d",timer.CurrentTime.Day(),timer.CurrentTime.Month().String(),timer.CurrentTime.Year())
	widget := &DockWidget { 
		SWS_CoreWidget: *corewidget,
		currentDay: sws.CreateLabel(150,25,today),
		timer: timer,
		timerevent: nil,
	}
	title:=sws.CreateLabel(150,25,"DC Tycoon")
	title.SetCentered(true)
	widget.AddChild(title)
	
	widget.currentDay.Move(5,25)
	widget.AddChild(widget.currentDay)
	
	widget.pause=sws.CreateFlatButtonWidget(25,25,"")
	widget.pause.SetColor(0xff8888ff)
	widget.pause.Move(25,50)
	widget.pause.SetImage("resources/icon-pause-symbol.png")
	widget.pause.SetClicked(func() {
		widget.pause.SetColor(0xff8888ff)
		widget.play.SetColor(0xffdddddd)
		widget.forward.SetColor(0xffdddddd)
		if widget.timerevent!=nil {
			widget.timerevent.StopRepeat()
		}
	})
	widget.AddChild(widget.pause)
	
	widget.play=sws.CreateFlatButtonWidget(25,25,"")
	widget.play.Move(50,50)
	widget.play.SetImage("resources/icon-arrowhead-pointing-to-the-right.png")
	widget.play.SetClicked(func() {
		widget.pause.SetColor(0xffdddddd)
		widget.play.SetColor(0xff8888ff)
		widget.forward.SetColor(0xffdddddd)
		if widget.timerevent!=nil {
			widget.timerevent.StopRepeat()
		}
		widget.timerevent=sws.TimerAddEvent(time.Now().Add(4*time.Second),4*time.Second,func() {
			timer.TimerClock()
			today:=fmt.Sprintf("%d %s %d",timer.CurrentTime.Day(),timer.CurrentTime.Month().String(),timer.CurrentTime.Year())
			widget.currentDay.SetText(today)
		})
	})
	widget.AddChild(widget.play)
	
	widget.forward=sws.CreateFlatButtonWidget(25,25,"")
	widget.forward.Move(75,50)
	widget.forward.SetImage("resources/icon-forward-button.png")
	widget.forward.SetClicked(func() {
		widget.pause.SetColor(0xffdddddd)
		widget.play.SetColor(0xffdddddd)
		widget.forward.SetColor(0xff8888ff)
		if widget.timerevent!=nil {
			widget.timerevent.StopRepeat()
		}
		widget.timerevent=sws.TimerAddEvent(time.Now().Add(time.Second),time.Second,func() {
			timer.TimerClock()
			today:=fmt.Sprintf("%d %s %d",timer.CurrentTime.Day(),timer.CurrentTime.Month().String(),timer.CurrentTime.Year())
			widget.currentDay.SetText(today)
		})
	})
	widget.AddChild(widget.forward)
	
	widget.shop=sws.CreateFlatButtonWidget(25,25,"")
	widget.shop.Move(25,75)
	widget.shop.SetImage("resources/icon-shopping-cart-black-shape.png")
	widget.AddChild(widget.shop)
	
	inventory:=sws.CreateFlatButtonWidget(25,25,"")
	inventory.Move(50,75)
	inventory.SetImage("resources/icon-delivery-truck-silhouette.png")
	widget.AddChild(inventory)
	
	save:=sws.CreateFlatButtonWidget(25,25,"")
	save.Move(75,75)
	save.SetImage("resources/icon-blank-file.png")
	widget.AddChild(save)
	
	widget.quit=sws.CreateFlatButtonWidget(25,25,"")
	widget.quit.Move(100,75)
	widget.quit.SetImage("resources/icon-power-button-off.png")
	widget.AddChild(widget.quit)
	
	widget.ledger=sws.CreateFlatButtonWidget(150,25,"")
	widget.ledger.Move(0,100)
	widget.AddChild(widget.ledger)
	accounting.GlobalLedger.AddSubscriber(widget)
	
	return widget
}
