package supplier

import(
	"github.com/nzin/sws"
	"github.com/veandco/go-sdl2/sdl"
)

//
// Page Shop
//
type ServerPageWidget struct {
	sws.CoreWidget
        towerbutton *sws.FlatButtonWidget
        rackbutton  *sws.FlatButtonWidget
        bladebutton *sws.FlatButtonWidget
}

func (self *ServerPageWidget) SetTowerCallback(callback func()) {
        self.towerbutton.SetClicked(callback)
}

func (self *ServerPageWidget) SetRackCallback(callback func()) {
        self.rackbutton.SetClicked(callback)
}

func (self *ServerPageWidget) SetBladeCallback(callback func()) {
        self.bladebutton.SetClicked(callback)
}

func NewServerPageWidget(width,height int32) *ServerPageWidget {
	serverpage:=&ServerPageWidget{
		CoreWidget: *sws.NewCoreWidget(width,height),
	}
	serverpage.SetColor(0xffffffff)

        servertitle:=sws.NewLabelWidget(150,40,"DEAL Servers")
        servertitle.SetFont(sws.LatoRegular20)
        servertitle.SetCentered(false)
	servertitle.SetColor(0xffffffff)
        servertitle.Move(40,0)
        serverpage.AddChild(servertitle)

        towerservers:=sws.NewFlatButtonWidget(120,20,"> Tower Servers")
        towerservers.SetTextColor(sdl.Color{0x06,0x84,0xdc,0xff})
        towerservers.SetCentered(false)
	towerservers.SetColor(0xffffffff)
        towerservers.Move(0,40)
	serverpage.towerbutton=towerservers
        serverpage.AddChild(towerservers)

        rackservers:=sws.NewFlatButtonWidget(120,20,"> Rack Servers")
        rackservers.SetTextColor(sdl.Color{0x06,0x84,0xdc,0xff})
        rackservers.SetCentered(false)
	rackservers.SetColor(0xffffffff)
        rackservers.Move(0,60)
	serverpage.rackbutton=rackservers
        serverpage.AddChild(rackservers)

        bladeservers:=sws.NewFlatButtonWidget(120,20,"> Blade Servers")
        bladeservers.SetTextColor(sdl.Color{0x06,0x84,0xdc,0xff})
        bladeservers.SetCentered(false)
	bladeservers.SetColor(0xffffffff)
        bladeservers.Move(0,80)
	serverpage.bladebutton=bladeservers
        serverpage.AddChild(bladeservers)
	
	return serverpage
}

