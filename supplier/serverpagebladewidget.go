package supplier

import (
	"github.com/nzin/sws"
	"github.com/veandco/go-sdl2/sdl"
)

//
// Page Shop>>Explore>>Blade servers
//
type ServerPageBladeWidget struct {
	sws.CoreWidget
	configureblade2     *sws.ButtonWidget
	configurebladeflat2 *sws.FlatButtonWidget
	configureblade1     *sws.ButtonWidget
	configurebladeflat1 *sws.FlatButtonWidget
}

func (self *ServerPageBladeWidget) SetConfigureBlade1Callback(callback func()) {
	self.configureblade1.SetClicked(callback)
	self.configurebladeflat1.SetClicked(callback)
}

func (self *ServerPageBladeWidget) SetConfigureBlade2Callback(callback func()) {
	self.configureblade2.SetClicked(callback)
	self.configurebladeflat2.SetClicked(callback)
}

func NewServerPageBladeWidget(width, height int32) *ServerPageBladeWidget {
	serverpageblade := &ServerPageBladeWidget{
		CoreWidget: *sws.NewCoreWidget(width, height),
	}
	serverpageblade.SetColor(0xffeeeeee)

	title := sws.NewLabelWidget(200, 20, "DEAL Blade Servers")
	title.SetFont(sws.LatoRegular20)
	title.SetColor(0xffeeeeee)
	title.Move(20, 0)
	title.SetCentered(false)
	serverpageblade.AddChild(title)

	blade1Icon := sws.NewFlatButtonWidget(150, 100, "")
	blade1Icon.SetImage("resources/server.blade.8u0.png")
	blade1Icon.SetColor(0xffeeeeee)
	blade1Icon.SetCentered(true)
	blade1Icon.Move(0, 20)
	serverpageblade.AddChild(blade1Icon)
	serverpageblade.configurebladeflat1 = blade1Icon

	blade1Title := sws.NewLabelWidget(150, 20, "Blade B100 server")
	blade1Title.SetColor(0xffeeeeee)
	blade1Title.SetTextColor(sdl.Color{0x06, 0x84, 0xdc, 0xff})
	blade1Title.Move(0, 120)
	serverpageblade.AddChild(blade1Title)

	blade1Desc := sws.NewTextAreaWidget(150, 160, "Simple, efficient but powerfull blade server, this 8U server is pre-equipped with 8 blades, each with 4 slots of RAM")
	blade1Desc.SetReadonly(true)
	blade1Desc.SetFont(sws.LatoRegular14)
	blade1Desc.SetColor(0xffeeeeee)
	blade1Desc.Move(0, 160)
	serverpageblade.AddChild(blade1Desc)

	blade1Button := sws.NewButtonWidget(100, 25, "Configure >")
	blade1Button.SetColor(0xffeeeeee)
	blade1Button.Move(0, 320)
	serverpageblade.AddChild(blade1Button)
	serverpageblade.configureblade1 = blade1Button

	blade2Icon := sws.NewFlatButtonWidget(150, 100, "")
	blade2Icon.SetImage("resources/server.blade.8u0.png")
	blade2Icon.SetColor(0xffeeeeee)
	blade2Icon.SetCentered(true)
	blade2Icon.Move(150, 20)
	serverpageblade.AddChild(blade2Icon)
	serverpageblade.configurebladeflat2 = blade2Icon

	blade2Title := sws.NewLabelWidget(150, 20, "Blade B200 server")
	blade2Title.SetColor(0xffeeeeee)
	blade2Title.SetTextColor(sdl.Color{0x06, 0x84, 0xdc, 0xff})
	blade2Title.Move(150, 120)
	serverpageblade.AddChild(blade2Title)

	blade2Desc := sws.NewTextAreaWidget(150, 160, "The ultimate 8U solution to all your problems, this pre-equipped 8 blades server has 4 slots of RAM and 2 processors per blade")
	blade2Desc.SetReadonly(true)
	blade2Desc.SetFont(sws.LatoRegular14)
	blade2Desc.SetColor(0xffeeeeee)
	blade2Desc.Move(150, 160)
	serverpageblade.AddChild(blade2Desc)

	blade2Button := sws.NewButtonWidget(100, 25, "Configure >")
	blade2Button.SetColor(0xffeeeeee)
	blade2Button.Move(150, 320)
	serverpageblade.AddChild(blade2Button)
	serverpageblade.configureblade2 = blade2Button

	return serverpageblade
}
