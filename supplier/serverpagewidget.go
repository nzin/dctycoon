package supplier

import (
	"github.com/nzin/sws"
	"github.com/veandco/go-sdl2/sdl"
)

//
// Page Shop
//
type ServerPageWidget struct {
	sws.CoreWidget
	towerbutton      *sws.FlatButtonWidget
	rackserverbutton *sws.FlatButtonWidget
	bladebutton      *sws.FlatButtonWidget
	acbutton         *sws.FlatButtonWidget
	generatorbutton  *sws.FlatButtonWidget
	rackbutton       *sws.FlatButtonWidget
}

func (self *ServerPageWidget) SetTowerCallback(callback func()) {
	self.towerbutton.SetClicked(callback)
}

func (self *ServerPageWidget) SetRackServerCallback(callback func()) {
	self.rackserverbutton.SetClicked(callback)
}

func (self *ServerPageWidget) SetBladeCallback(callback func()) {
	self.bladebutton.SetClicked(callback)
}

func (self *ServerPageWidget) SetAcCallback(callback func()) {
	self.acbutton.SetClicked(callback)
}

func (self *ServerPageWidget) SetRackCallback(callback func()) {
	self.rackbutton.SetClicked(callback)
}

func (self *ServerPageWidget) SetGeneratorCallback(callback func()) {
	self.generatorbutton.SetClicked(callback)
}

func NewServerPageWidget(width, height int32) *ServerPageWidget {
	serverpage := &ServerPageWidget{
		CoreWidget: *sws.NewCoreWidget(width, height),
	}
	serverpage.SetColor(0xffffffff)

	servertitle := sws.NewLabelWidget(150, 40, "DEAL Servers")
	servertitle.SetFont(sws.LatoRegular20)
	servertitle.SetCentered(false)
	servertitle.SetColor(0xffffffff)
	servertitle.Move(40, 0)
	serverpage.AddChild(servertitle)

	towerservers := sws.NewFlatButtonWidget(120, 20, "> Tower Servers")
	towerservers.SetTextColor(sdl.Color{0x06, 0x84, 0xdc, 0xff})
	towerservers.SetCentered(false)
	towerservers.SetColor(0xffffffff)
	towerservers.Move(0, 40)
	serverpage.towerbutton = towerservers
	serverpage.AddChild(towerservers)

	rackserverservers := sws.NewFlatButtonWidget(120, 20, "> Rack Servers")
	rackserverservers.SetTextColor(sdl.Color{0x06, 0x84, 0xdc, 0xff})
	rackserverservers.SetCentered(false)
	rackserverservers.SetColor(0xffffffff)
	rackserverservers.Move(0, 60)
	serverpage.rackserverbutton = rackserverservers
	serverpage.AddChild(rackserverservers)

	bladeservers := sws.NewFlatButtonWidget(120, 20, "> Blade Servers")
	bladeservers.SetTextColor(sdl.Color{0x06, 0x84, 0xdc, 0xff})
	bladeservers.SetCentered(false)
	bladeservers.SetColor(0xffffffff)
	bladeservers.Move(0, 80)
	serverpage.bladebutton = bladeservers
	serverpage.AddChild(bladeservers)

	ac := sws.NewFlatButtonWidget(120, 20, "> Air climatiser")
	ac.SetTextColor(sdl.Color{0x06, 0x84, 0xdc, 0xff})
	ac.SetCentered(false)
	ac.SetColor(0xffffffff)
	ac.Move(0, 100)
	serverpage.acbutton = ac
	serverpage.AddChild(ac)

	generator := sws.NewFlatButtonWidget(120, 20, "> Generator")
	generator.SetTextColor(sdl.Color{0x06, 0x84, 0xdc, 0xff})
	generator.SetCentered(false)
	generator.SetColor(0xffffffff)
	generator.Move(0, 120)
	serverpage.generatorbutton = generator
	serverpage.AddChild(generator)

	rack := sws.NewFlatButtonWidget(120, 20, "> Rack")
	rack.SetTextColor(sdl.Color{0x06, 0x84, 0xdc, 0xff})
	rack.SetCentered(false)
	rack.SetColor(0xffffffff)
	rack.Move(0, 140)
	serverpage.rackbutton = rack
	serverpage.AddChild(rack)

	return serverpage
}
