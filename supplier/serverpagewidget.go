package supplier

import(
	"github.com/nzin/sws"
	"github.com/veandco/go-sdl2/sdl"
)

type ServerPageWidget struct {
	sws.SWS_CoreWidget
        towerbutton *sws.SWS_FlatButtonWidget
        rackbutton  *sws.SWS_FlatButtonWidget
        bladebutton *sws.SWS_FlatButtonWidget
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

func CreateServerPageWidget(width,height int32) *ServerPageWidget {
	serverpage:=&ServerPageWidget{
		SWS_CoreWidget: *sws.CreateCoreWidget(width,height),
	}

        servertitle:=sws.CreateLabel(150,40,"DEAL Servers")
        servertitle.SetFont(sws.LatoRegular20)
        servertitle.SetCentered(false)
        servertitle.Move(40,0)
        serverpage.AddChild(servertitle)

        towerservers:=sws.CreateFlatButtonWidget(120,20,"> Tower Servers")
        towerservers.SetTextColor(sdl.Color{0x06,0x84,0xdc,0xff})
        towerservers.SetCentered(false)
        towerservers.Move(0,40)
	serverpage.towerbutton=towerservers
        serverpage.AddChild(towerservers)

        rackservers:=sws.CreateFlatButtonWidget(120,20,"> Rack Servers")
        rackservers.SetTextColor(sdl.Color{0x06,0x84,0xdc,0xff})
        rackservers.SetCentered(false)
        rackservers.Move(0,60)
	serverpage.rackbutton=rackservers
        serverpage.AddChild(rackservers)

        bladeservers:=sws.CreateFlatButtonWidget(120,20,"> Blade Servers")
        bladeservers.SetTextColor(sdl.Color{0x06,0x84,0xdc,0xff})
        bladeservers.SetCentered(false)
        bladeservers.Move(0,80)
	serverpage.bladebutton=bladeservers
        serverpage.AddChild(bladeservers)
	
	return serverpage
}

