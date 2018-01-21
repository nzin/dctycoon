package supplier

import (
	"github.com/nzin/dctycoon/global"
	"github.com/nzin/sws"
	"github.com/veandco/go-sdl2/sdl"
)

//
// Page Shop>>Explore
//
type ServerPageExploreWidget struct {
	sws.CoreWidget
	towerbutton      *sws.ButtonWidget
	towerflat        *sws.FlatButtonWidget
	rackserverbutton *sws.ButtonWidget
	rackserverflat   *sws.FlatButtonWidget
	bladebutton      *sws.ButtonWidget
	bladeflat        *sws.FlatButtonWidget
	acbutton         *sws.ButtonWidget
	acflat           *sws.FlatButtonWidget
	rackbutton       *sws.ButtonWidget
	rackflat         *sws.FlatButtonWidget
	generatorbutton  *sws.ButtonWidget
	generatorflat    *sws.FlatButtonWidget
}

func (self *ServerPageExploreWidget) SetTowerCallback(callback func()) {
	self.towerbutton.SetClicked(callback)
	self.towerflat.SetClicked(callback)
}

func (self *ServerPageExploreWidget) SetRackServerCallback(callback func()) {
	self.rackserverbutton.SetClicked(callback)
	self.rackserverflat.SetClicked(callback)
}

func (self *ServerPageExploreWidget) SetBladeCallback(callback func()) {
	self.bladebutton.SetClicked(callback)
	self.bladeflat.SetClicked(callback)
}

func (self *ServerPageExploreWidget) SetAcCallback(callback func()) {
	self.acbutton.SetClicked(callback)
	self.acflat.SetClicked(callback)
}

func (self *ServerPageExploreWidget) SetRackCallback(callback func()) {
	self.rackbutton.SetClicked(callback)
	self.rackflat.SetClicked(callback)
}

func (self *ServerPageExploreWidget) SetGeneratorCallback(callback func()) {
	self.generatorbutton.SetClicked(callback)
	self.generatorflat.SetClicked(callback)
}

func NewServerPageExploreWidget(width, height int32) *ServerPageExploreWidget {
	serverpageexplore := &ServerPageExploreWidget{
		CoreWidget: *sws.NewCoreWidget(width, height),
	}
	serverpageexplore.SetColor(0xffeeeeee)

	title := sws.NewLabelWidget(200, 20, "Explore DEAL Servers")
	title.SetFont(sws.LatoRegular20)
	title.SetColor(0xffeeeeee)
	title.Move(20, 0)
	title.SetCentered(false)
	serverpageexplore.AddChild(title)

	towerIcon := sws.NewFlatButtonWidget(150, 100, "")
	if img, err := global.LoadImageAsset("assets/ui/tower0.png"); err == nil {
		towerIcon.SetImageSurface(img)
	}
	towerIcon.SetColor(0xffeeeeee)
	towerIcon.SetCentered(true)
	towerIcon.Move(0, 20)
	serverpageexplore.AddChild(towerIcon)
	serverpageexplore.towerflat = towerIcon

	towerTitle := sws.NewLabelWidget(150, 20, "Tower servers")
	towerTitle.SetColor(0xffeeeeee)
	towerTitle.SetTextColor(sdl.Color{0x06, 0x84, 0xdc, 0xff})
	towerTitle.Move(0, 120)
	serverpageexplore.AddChild(towerTitle)

	towerDesc := sws.NewTextAreaWidget(150, 160, "Our professional workstation with up to 2 processors, is the ideal powerhouse machine you need to tackle your engineering problem")
	towerDesc.SetDisabled(true)
	towerDesc.SetFont(sws.LatoRegular14)
	towerDesc.SetColor(0xffeeeeee)
	towerDesc.Move(0, 160)
	serverpageexplore.AddChild(towerDesc)

	towerButton := sws.NewButtonWidget(100, 25, "Know more >")
	towerButton.SetColor(0xffeeeeee)
	towerButton.Move(0, 320)
	serverpageexplore.towerbutton = towerButton
	serverpageexplore.AddChild(towerButton)

	rackserverIcon := sws.NewFlatButtonWidget(150, 100, "")
	if img, err := global.LoadImageAsset("assets/ui/server.2u0.png"); err == nil {
		rackserverIcon.SetImageSurface(img)
	}
	rackserverIcon.SetColor(0xffeeeeee)
	rackserverIcon.SetCentered(true)
	rackserverIcon.Move(150, 20)
	serverpageexplore.AddChild(rackserverIcon)
	serverpageexplore.rackserverflat = rackserverIcon

	rackserverTitle := sws.NewLabelWidget(150, 20, "Rack servers")
	rackserverTitle.SetColor(0xffeeeeee)
	rackserverTitle.SetTextColor(sdl.Color{0x06, 0x84, 0xdc, 0xff})
	rackserverTitle.Move(150, 120)
	serverpageexplore.AddChild(rackserverTitle)

	rackserverDesc := sws.NewTextAreaWidget(150, 160, "Discover our large choice of rackserver server, from 1U to 4U, to tackle all your datacenter needs")
	rackserverDesc.SetDisabled(true)
	rackserverDesc.SetFont(sws.LatoRegular14)
	rackserverDesc.SetColor(0xffeeeeee)
	rackserverDesc.Move(150, 160)
	serverpageexplore.AddChild(rackserverDesc)

	rackserverButton := sws.NewButtonWidget(100, 25, "Know more >")
	rackserverButton.SetColor(0xffeeeeee)
	rackserverButton.Move(150, 320)
	serverpageexplore.rackserverbutton = rackserverButton
	serverpageexplore.AddChild(rackserverButton)

	bladeIcon := sws.NewFlatButtonWidget(150, 100, "")
	if img, err := global.LoadImageAsset("assets/ui/server.blade.8u0.png"); err == nil {
		bladeIcon.SetImageSurface(img)
	}
	bladeIcon.SetColor(0xffeeeeee)
	bladeIcon.SetCentered(true)
	bladeIcon.Move(300, 20)
	serverpageexplore.AddChild(bladeIcon)
	serverpageexplore.bladeflat = bladeIcon

	bladeTitle := sws.NewLabelWidget(150, 20, "Blade servers")
	bladeTitle.SetColor(0xffeeeeee)
	bladeTitle.SetTextColor(sdl.Color{0x06, 0x84, 0xdc, 0xff})
	bladeTitle.Move(300, 120)
	serverpageexplore.AddChild(bladeTitle)

	bladeDesc := sws.NewTextAreaWidget(150, 160, "For maximum rackserver density we propose our best in the class 8U blade server offers, with 8 blades (max 2 CPU per blade)")
	bladeDesc.SetDisabled(true)
	bladeDesc.SetFont(sws.LatoRegular14)
	bladeDesc.SetColor(0xffeeeeee)
	bladeDesc.Move(300, 160)
	serverpageexplore.AddChild(bladeDesc)

	bladeButton := sws.NewButtonWidget(100, 25, "Know more >")
	bladeButton.SetColor(0xffeeeeee)
	bladeButton.Move(300, 320)
	serverpageexplore.bladebutton = bladeButton
	serverpageexplore.AddChild(bladeButton)

	// next stage + 340

	acIcon := sws.NewFlatButtonWidget(150, 100, "")
	if img, err := global.LoadImageAsset("assets/ui/ac0.100.png"); err == nil {
		acIcon.SetImageSurface(img)
	}
	acIcon.SetColor(0xffeeeeee)
	acIcon.SetCentered(true)
	acIcon.Move(0, 360)
	serverpageexplore.AddChild(acIcon)
	serverpageexplore.acflat = acIcon

	acTitle := sws.NewLabelWidget(150, 20, "Air Climatiser")
	acTitle.SetColor(0xffeeeeee)
	acTitle.SetTextColor(sdl.Color{0x06, 0x84, 0xdc, 0xff})
	acTitle.Move(0, 460)
	serverpageexplore.AddChild(acTitle)

	acDesc := sws.NewTextAreaWidget(150, 160, "Our next generation efficient Data Center Air Climatiser")
	acDesc.SetDisabled(true)
	acDesc.SetFont(sws.LatoRegular14)
	acDesc.SetColor(0xffeeeeee)
	acDesc.Move(0, 500)
	serverpageexplore.AddChild(acDesc)

	acButton := sws.NewButtonWidget(100, 25, "Buy now >")
	acButton.SetColor(0xffeeeeee)
	acButton.Move(0, 660)
	serverpageexplore.acbutton = acButton
	serverpageexplore.AddChild(acButton)

	rackIcon := sws.NewFlatButtonWidget(150, 100, "")
	if img, err := global.LoadImageAsset("assets/ui/rack0.100.png"); err == nil {
		rackIcon.SetImageSurface(img)
	}
	rackIcon.SetColor(0xffeeeeee)
	rackIcon.SetCentered(true)
	rackIcon.Move(150, 360)
	serverpageexplore.AddChild(rackIcon)
	serverpageexplore.rackflat = rackIcon

	rackTitle := sws.NewLabelWidget(150, 20, "Rack")
	rackTitle.SetColor(0xffeeeeee)
	rackTitle.SetTextColor(sdl.Color{0x06, 0x84, 0xdc, 0xff})
	rackTitle.Move(150, 460)
	serverpageexplore.AddChild(rackTitle)

	rackDesc := sws.NewTextAreaWidget(150, 160, "Classic 42U Rack chassis. Up to 64A")
	rackDesc.SetDisabled(true)
	rackDesc.SetFont(sws.LatoRegular14)
	rackDesc.SetColor(0xffeeeeee)
	rackDesc.Move(150, 500)
	serverpageexplore.AddChild(rackDesc)

	rackButton := sws.NewButtonWidget(100, 25, "Buy now >")
	rackButton.SetColor(0xffeeeeee)
	rackButton.Move(150, 660)
	serverpageexplore.rackbutton = rackButton
	serverpageexplore.AddChild(rackButton)

	generatorIcon := sws.NewFlatButtonWidget(150, 100, "")
	if img, err := global.LoadImageAsset("assets/ui/generator0.100.png"); err == nil {
		generatorIcon.SetImageSurface(img)
	}
	generatorIcon.SetColor(0xffeeeeee)
	generatorIcon.SetCentered(true)
	generatorIcon.Move(300, 360)
	serverpageexplore.AddChild(generatorIcon)
	serverpageexplore.generatorflat = generatorIcon

	generatorTitle := sws.NewLabelWidget(150, 20, "Generator")
	generatorTitle.SetColor(0xffeeeeee)
	generatorTitle.SetTextColor(sdl.Color{0x06, 0x84, 0xdc, 0xff})
	generatorTitle.Move(300, 460)
	serverpageexplore.AddChild(generatorTitle)

	generatorDesc := sws.NewTextAreaWidget(150, 160, "Prevent electriciy outage with this powerfull yet compact diesel generator. Generates 50 kwh")
	generatorDesc.SetDisabled(true)
	generatorDesc.SetFont(sws.LatoRegular14)
	generatorDesc.SetColor(0xffeeeeee)
	generatorDesc.Move(300, 500)
	serverpageexplore.AddChild(generatorDesc)

	generatorButton := sws.NewButtonWidget(100, 25, "Buy now >")
	generatorButton.SetColor(0xffeeeeee)
	generatorButton.Move(300, 660)
	serverpageexplore.generatorbutton = generatorButton
	serverpageexplore.AddChild(generatorButton)

	return serverpageexplore
}
