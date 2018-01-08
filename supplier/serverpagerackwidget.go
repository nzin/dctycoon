package supplier

import (
	"github.com/nzin/dctycoon/global"
	"github.com/nzin/sws"
	"github.com/veandco/go-sdl2/sdl"
)

//
// Page Shop>>Explore>>Rack servers
//
type ServerPageRackWidget struct {
	sws.CoreWidget
	configurerack1     *sws.ButtonWidget
	configurerackflat1 *sws.FlatButtonWidget
	configurerack2     *sws.ButtonWidget
	configurerackflat2 *sws.FlatButtonWidget
	configurerack4     *sws.ButtonWidget
	configurerackflat4 *sws.FlatButtonWidget
	configurerack6     *sws.ButtonWidget
	configurerackflat6 *sws.FlatButtonWidget
}

func (self *ServerPageRackWidget) SetConfigureRack1Callback(callback func()) {
	self.configurerack1.SetClicked(callback)
	self.configurerackflat1.SetClicked(callback)
}

func (self *ServerPageRackWidget) SetConfigureRack2Callback(callback func()) {
	self.configurerack2.SetClicked(callback)
	self.configurerackflat2.SetClicked(callback)
}

func (self *ServerPageRackWidget) SetConfigureRack4Callback(callback func()) {
	self.configurerack4.SetClicked(callback)
	self.configurerackflat4.SetClicked(callback)
}

func (self *ServerPageRackWidget) SetConfigureRack6Callback(callback func()) {
	self.configurerack6.SetClicked(callback)
	self.configurerackflat6.SetClicked(callback)
}

func NewServerPageRackWidget(width, height int32) *ServerPageRackWidget {
	serverpagerack := &ServerPageRackWidget{
		CoreWidget: *sws.NewCoreWidget(width, height),
	}
	serverpagerack.SetColor(0xffeeeeee)

	title := sws.NewLabelWidget(200, 20, "DEAL Rack Servers")
	title.SetFont(sws.LatoRegular20)
	title.SetColor(0xffeeeeee)
	title.Move(20, 0)
	title.SetCentered(false)
	serverpagerack.AddChild(title)

	rack1Icon := sws.NewFlatButtonWidget(150, 100, "")
	if img, err := global.LoadImageAsset("assets/ui/server.1u0.png"); err == nil {
		rack1Icon.SetImageSurface(img)
	}
	rack1Icon.SetColor(0xffeeeeee)
	rack1Icon.SetCentered(true)
	rack1Icon.Move(0, 20)
	serverpagerack.AddChild(rack1Icon)
	serverpagerack.configurerackflat1 = rack1Icon

	rack1Title := sws.NewLabelWidget(150, 20, "Rack R100 server")
	rack1Title.SetColor(0xffeeeeee)
	rack1Title.SetTextColor(sdl.Color{0x06, 0x84, 0xdc, 0xff})
	rack1Title.Move(0, 120)
	serverpagerack.AddChild(rack1Title)

	rack1Desc := sws.NewTextAreaWidget(150, 160, "This thin yet powerfull 1U rack server is equipped from 1 to 2 CPU and up to 2 RAM slots for the best processing job you need")
	rack1Desc.SetDisabled(true)
	rack1Desc.SetFont(sws.LatoRegular14)
	rack1Desc.SetColor(0xffeeeeee)
	rack1Desc.Move(0, 160)
	serverpagerack.AddChild(rack1Desc)

	rack1Button := sws.NewButtonWidget(100, 25, "Configure >")
	rack1Button.SetColor(0xffeeeeee)
	rack1Button.Move(0, 320)
	serverpagerack.AddChild(rack1Button)
	serverpagerack.configurerack1 = rack1Button

	rack2Icon := sws.NewFlatButtonWidget(150, 100, "")
	if img, err := global.LoadImageAsset("assets/ui/server.2u0.png"); err == nil {
		rack2Icon.SetImageSurface(img)
	}
	rack2Icon.SetColor(0xffeeeeee)
	rack2Icon.SetCentered(true)
	rack2Icon.Move(150, 20)
	serverpagerack.AddChild(rack2Icon)
	serverpagerack.configurerackflat2 = rack2Icon

	rack2Title := sws.NewLabelWidget(150, 20, "Rack R200 server")
	rack2Title.SetColor(0xffeeeeee)
	rack2Title.SetTextColor(sdl.Color{0x06, 0x84, 0xdc, 0xff})
	rack2Title.Move(150, 120)
	serverpagerack.AddChild(rack2Title)

	rack2Desc := sws.NewTextAreaWidget(150, 160, "Our best balanced 2U rack server, is equipped from 1 to 2 CPU and up to 4 RAM slots and up to 4 disks to fit all job you require")
	rack2Desc.SetDisabled(true)
	rack2Desc.SetFont(sws.LatoRegular14)
	rack2Desc.SetColor(0xffeeeeee)
	rack2Desc.Move(150, 160)
	serverpagerack.AddChild(rack2Desc)

	rack2Button := sws.NewButtonWidget(100, 25, "Configure >")
	rack2Button.SetColor(0xffeeeeee)
	rack2Button.Move(150, 320)
	serverpagerack.AddChild(rack2Button)
	serverpagerack.configurerack2 = rack2Button

	rack4Icon := sws.NewFlatButtonWidget(150, 100, "")
	if img, err := global.LoadImageAsset("assets/ui/server.4u0.png"); err == nil {
		rack4Icon.SetImageSurface(img)
	}
	rack4Icon.SetColor(0xffeeeeee)
	rack4Icon.SetCentered(true)
	rack4Icon.Move(300, 20)
	serverpagerack.AddChild(rack4Icon)
	serverpagerack.configurerackflat4 = rack4Icon

	rack4Title := sws.NewLabelWidget(150, 20, "Rack R400 server")
	rack4Title.SetColor(0xffeeeeee)
	rack4Title.SetTextColor(sdl.Color{0x06, 0x84, 0xdc, 0xff})
	rack4Title.Move(300, 120)
	serverpagerack.AddChild(rack4Title)

	rack4Desc := sws.NewTextAreaWidget(150, 160, "This 4U storage rack server, with 1 to 2 CPU and up to 10 disks is specifically designed for your heavy storage fulfillment")
	rack4Desc.SetDisabled(true)
	rack4Desc.SetFont(sws.LatoRegular14)
	rack4Desc.SetColor(0xffeeeeee)
	rack4Desc.Move(300, 160)
	serverpagerack.AddChild(rack4Desc)

	rack4Button := sws.NewButtonWidget(100, 25, "Configure >")
	rack4Button.SetColor(0xffeeeeee)
	rack4Button.Move(300, 320)
	serverpagerack.AddChild(rack4Button)
	serverpagerack.configurerack4 = rack4Button

	rack6Icon := sws.NewFlatButtonWidget(150, 100, "")
	if img, err := global.LoadImageAsset("assets/ui/server.4u0.png"); err == nil {
		rack6Icon.SetImageSurface(img)
	}
	rack6Icon.SetColor(0xffeeeeee)
	rack6Icon.SetCentered(true)
	rack6Icon.Move(0, 360)
	serverpagerack.AddChild(rack6Icon)
	serverpagerack.configurerackflat6 = rack6Icon

	rack6Title := sws.NewLabelWidget(150, 20, "Rack R600 server")
	rack6Title.SetColor(0xffeeeeee)
	rack6Title.SetTextColor(sdl.Color{0x06, 0x84, 0xdc, 0xff})
	rack6Title.Move(0, 460)
	serverpagerack.AddChild(rack6Title)

	rack6Desc := sws.NewTextAreaWidget(150, 160, "This 4U HPC rack server, with 2 to 4 CPU and up to 8 RAM slots will put in the dust any cpu-intensive job")
	rack6Desc.SetDisabled(true)
	rack6Desc.SetFont(sws.LatoRegular14)
	rack6Desc.SetColor(0xffeeeeee)
	rack6Desc.Move(0, 500)
	serverpagerack.AddChild(rack6Desc)

	rack6Button := sws.NewButtonWidget(100, 25, "Configure >")
	rack6Button.SetColor(0xffeeeeee)
	rack6Button.Move(0, 660)
	serverpagerack.AddChild(rack6Button)
	serverpagerack.configurerack6 = rack6Button

	return serverpagerack
}
