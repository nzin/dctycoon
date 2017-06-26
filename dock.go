package dctycoon

import(
	"github.com/nzin/sws"
	"fmt"
)

type DockWidget struct {
	sws.SWS_CoreWidget
	currentDay    *sws.SWS_Label
	timer         *Timer
	shop          *sws.SWS_FlatButtonWidget
	quit          *sws.SWS_FlatButtonWidget
}

func (self *DockWidget) SetQuitCallback(callback func()) {
	self.quit.SetClicked(callback)
}

func (self *DockWidget) SetShopCallback(callback func()) {
	self.shop.SetClicked(callback)
}

func CreateDockWidget(timer *Timer) *DockWidget {
	corewidget := sws.CreateCoreWidget(100, 100)
	today:=fmt.Sprintf("%d %s %d",timer.CurrentTime.Day(),timer.CurrentTime.Month().String(),timer.CurrentTime.Year())
	widget := &DockWidget { 
		SWS_CoreWidget: *corewidget,
		currentDay: sws.CreateLabel(100,25,today),
		timer: timer,
	}
	title:=sws.CreateLabel(100,25,"DC Tycoon")
	title.SetCentered(true)
	widget.AddChild(title)
	
	widget.currentDay.Move(0,25)
	widget.AddChild(widget.currentDay)
	
	pause:=sws.CreateFlatButtonWidget(25,25,"")
	pause.Move(0,50)
	pause.SetImage("resources/icon-pause-symbol.png")
	widget.AddChild(pause)
	
	play:=sws.CreateFlatButtonWidget(25,25,"")
	play.Move(25,50)
	play.SetImage("resources/icon-arrowhead-pointing-to-the-right.png")
	widget.AddChild(play)
	
	forward:=sws.CreateFlatButtonWidget(25,25,"")
	forward.Move(50,50)
	forward.SetImage("resources/icon-forward-button.png")
	widget.AddChild(forward)
	
	widget.shop=sws.CreateFlatButtonWidget(25,25,"")
	widget.shop.Move(0,75)
	widget.shop.SetImage("resources/icon-shopping-cart-black-shape.png")
	widget.AddChild(widget.shop)
	
	inventory:=sws.CreateFlatButtonWidget(25,25,"")
	inventory.Move(25,75)
	inventory.SetImage("resources/icon-delivery-truck-silhouette.png")
	widget.AddChild(inventory)
	
	save:=sws.CreateFlatButtonWidget(25,25,"")
	save.Move(50,75)
	save.SetImage("resources/icon-blank-file.png")
	widget.AddChild(save)
	
	widget.quit=sws.CreateFlatButtonWidget(25,25,"")
	widget.quit.Move(75,75)
	widget.quit.SetImage("resources/icon-power-button-off.png")
	widget.AddChild(widget.quit)
	
	return widget
}
