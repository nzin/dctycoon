package supplier

import(
	"github.com/nzin/sws"
)

//
// Page Shop>>Explore
//
type ServerPageExploreWidget struct {
	sws.SWS_CoreWidget
	towerbutton *sws.SWS_ButtonWidget
	rackbutton  *sws.SWS_ButtonWidget
	bladebutton *sws.SWS_ButtonWidget
}

func (self *ServerPageExploreWidget) SetTowerCallback(callback func()) {
	self.towerbutton.SetClicked(callback)
}

func (self *ServerPageExploreWidget) SetRackCallback(callback func()) {
	self.rackbutton.SetClicked(callback)
}

func (self *ServerPageExploreWidget) SetBladeCallback(callback func()) {
	self.bladebutton.SetClicked(callback)
}

func CreateServerPageExploreWidget(width,height int32) *ServerPageExploreWidget {
	serverpageexplore:=&ServerPageExploreWidget{
		SWS_CoreWidget: *sws.CreateCoreWidget(width,height),
	}
	serverpageexplore.SetColor(0xffffffff)
	
        title:=sws.CreateLabel(200,20,"Explore DEAL Servers")
        title.SetFont(sws.LatoRegular20)
	title.SetColor(0xffffffff)
        title.Move(20,0)
        title.SetCentered(false)
        serverpageexplore.AddChild(title)

	towerIcon:=sws.CreateLabel(150,100,"")
	towerIcon.SetImage("resources/tower0.png")
	towerIcon.SetColor(0xffffffff)
        towerIcon.SetCentered(true)
	towerIcon.Move(0,20)
        serverpageexplore.AddChild(towerIcon)

	towerTitle:=sws.CreateLabel(150,20,"Tower servers")
	towerTitle.SetColor(0xffffffff)
	towerTitle.Move(0,120)
        serverpageexplore.AddChild(towerTitle)

	towerButton:=sws.CreateButtonWidget(100,25,"Know more >")
	towerButton.SetColor(0xffffffff)
	towerButton.Move(0,320)
	serverpageexplore.towerbutton=towerButton
	serverpageexplore.AddChild(towerButton)


	rackIcon:=sws.CreateLabel(150,100,"")
	rackIcon.SetImage("resources/server.2u0.png")
	rackIcon.SetColor(0xffffffff)
        rackIcon.SetCentered(true)
	rackIcon.Move(150,20)
        serverpageexplore.AddChild(rackIcon)

	rackTitle:=sws.CreateLabel(150,20,"Rack servers")
	rackTitle.SetColor(0xffffffff)
	rackTitle.Move(150,120)
        serverpageexplore.AddChild(rackTitle)

	rackButton:=sws.CreateButtonWidget(100,25,"Know more >")
	rackButton.SetColor(0xffffffff)
	rackButton.Move(150,320)
	serverpageexplore.rackbutton=rackButton
	serverpageexplore.AddChild(rackButton)


	bladeIcon:=sws.CreateLabel(150,100,"")
	bladeIcon.SetImage("resources/server.blade.8u0.png")
	bladeIcon.SetColor(0xffffffff)
        bladeIcon.SetCentered(true)
	bladeIcon.Move(300,20)
        serverpageexplore.AddChild(bladeIcon)

	bladeTitle:=sws.CreateLabel(150,20,"Blade servers")
	bladeTitle.SetColor(0xffffffff)
	bladeTitle.Move(300,120)
	serverpageexplore.AddChild(bladeTitle)

	bladeButton:=sws.CreateButtonWidget(100,25,"Know more >")
	bladeButton.SetColor(0xffffffff)
	bladeButton.Move(300,320)
	serverpageexplore.bladebutton=bladeButton
	serverpageexplore.AddChild(bladeButton)

	return serverpageexplore
}

