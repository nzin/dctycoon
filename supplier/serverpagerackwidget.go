package supplier

import(
	"github.com/veandco/go-sdl2/sdl"
	"github.com/nzin/sws"
)

//
// Page Shop>>Explore>>Rack servers
//
type ServerPageRackWidget struct {
	sws.SWS_CoreWidget
	configurerack1     *sws.SWS_ButtonWidget
	configurerackflat1 *sws.SWS_FlatButtonWidget
	configurerack2     *sws.SWS_ButtonWidget
	configurerackflat2 *sws.SWS_FlatButtonWidget
	configurerack4     *sws.SWS_ButtonWidget
	configurerackflat4 *sws.SWS_FlatButtonWidget
	configurerack6     *sws.SWS_ButtonWidget
	configurerackflat6 *sws.SWS_FlatButtonWidget
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

func CreateServerPageRackWidget(width,height int32) *ServerPageRackWidget {
	serverpagerack:=&ServerPageRackWidget{
		SWS_CoreWidget: *sws.CreateCoreWidget(width,height),
	}
	serverpagerack.SetColor(0xffeeeeee)
	
        title:=sws.CreateLabel(200,20,"DEAL Rack Servers")
        title.SetFont(sws.LatoRegular20)
	title.SetColor(0xffeeeeee)
        title.Move(20,0)
        title.SetCentered(false)
        serverpagerack.AddChild(title)

	rack1Icon:=sws.CreateFlatButtonWidget(150,100,"")
	rack1Icon.SetImage("resources/server.1u0.png")
	rack1Icon.SetColor(0xffeeeeee)
        rack1Icon.SetCentered(true)
	rack1Icon.Move(0,20)
        serverpagerack.AddChild(rack1Icon)
        serverpagerack.configurerackflat1=rack1Icon

	rack1Title:=sws.CreateLabel(150,20,"Rack R100 server")
	rack1Title.SetColor(0xffeeeeee)
	rack1Title.SetTextColor(sdl.Color{0x06,0x84,0xdc,0xff})
	rack1Title.Move(0,120)
        serverpagerack.AddChild(rack1Title)

	rack1Desc:=sws.CreateTextAreaWidget(150,160,"This thin yet powerfull 1U rack server is equipped from 1 to 2 CPU and up to 2 RAM slots for the best processing job you need")
	rack1Desc.SetReadonly(true)
	rack1Desc.SetFont(sws.LatoRegular14)
	rack1Desc.SetColor(0xffeeeeee)
	rack1Desc.Move(0,160)
	serverpagerack.AddChild(rack1Desc)

	rack1Button:=sws.CreateButtonWidget(100,25,"Configure >")
	rack1Button.SetColor(0xffeeeeee)
	rack1Button.Move(0,320)
	serverpagerack.AddChild(rack1Button)
	serverpagerack.configurerack1=rack1Button


	rack2Icon:=sws.CreateFlatButtonWidget(150,100,"")
	rack2Icon.SetImage("resources/server.2u0.png")
	rack2Icon.SetColor(0xffeeeeee)
        rack2Icon.SetCentered(true)
	rack2Icon.Move(150,20)
        serverpagerack.AddChild(rack2Icon)
        serverpagerack.configurerackflat2=rack2Icon

	rack2Title:=sws.CreateLabel(150,20,"Rack R200 server")
	rack2Title.SetColor(0xffeeeeee)
	rack2Title.SetTextColor(sdl.Color{0x06,0x84,0xdc,0xff})
	rack2Title.Move(150,120)
        serverpagerack.AddChild(rack2Title)

	rack2Desc:=sws.CreateTextAreaWidget(150,160,"Our best balanced 2U rack server, is equipped from 1 to 2 CPU and up to 4 RAM slots and up to 4 disks to fit all job you require")
	rack2Desc.SetReadonly(true)
	rack2Desc.SetFont(sws.LatoRegular14)
	rack2Desc.SetColor(0xffeeeeee)
	rack2Desc.Move(150,160)
	serverpagerack.AddChild(rack2Desc)

	rack2Button:=sws.CreateButtonWidget(100,25,"Configure >")
	rack2Button.SetColor(0xffeeeeee)
	rack2Button.Move(150,320)
	serverpagerack.AddChild(rack2Button)
	serverpagerack.configurerack2=rack2Button


	rack4Icon:=sws.CreateFlatButtonWidget(150,100,"")
	rack4Icon.SetImage("resources/server.4u0.png")
	rack4Icon.SetColor(0xffeeeeee)
        rack4Icon.SetCentered(true)
	rack4Icon.Move(300,20)
        serverpagerack.AddChild(rack4Icon)
        serverpagerack.configurerackflat4=rack4Icon

	rack4Title:=sws.CreateLabel(150,20,"Rack R400 server")
	rack4Title.SetColor(0xffeeeeee)
	rack4Title.SetTextColor(sdl.Color{0x06,0x84,0xdc,0xff})
	rack4Title.Move(300,120)
	serverpagerack.AddChild(rack4Title)

	rack4Desc:=sws.CreateTextAreaWidget(150,160,"This 4U storage rack server, with 1 to 2 CPU and up to 10 disks is specifically designed for your heavy storage fulfillment")
	rack4Desc.SetReadonly(true)
	rack4Desc.SetFont(sws.LatoRegular14)
	rack4Desc.SetColor(0xffeeeeee)
	rack4Desc.Move(300,160)
	serverpagerack.AddChild(rack4Desc)

	rack4Button:=sws.CreateButtonWidget(100,25,"Configure >")
	rack4Button.SetColor(0xffeeeeee)
	rack4Button.Move(300,320)
	serverpagerack.AddChild(rack4Button)
	serverpagerack.configurerack4=rack4Button


	rack6Icon:=sws.CreateFlatButtonWidget(150,100,"")
	rack6Icon.SetImage("resources/server.4u0.png")
	rack6Icon.SetColor(0xffeeeeee)
        rack6Icon.SetCentered(true)
	rack6Icon.Move(0,360)
        serverpagerack.AddChild(rack6Icon)
        serverpagerack.configurerackflat6=rack6Icon

	rack6Title:=sws.CreateLabel(150,20,"Rack R600 server")
	rack6Title.SetColor(0xffeeeeee)
	rack6Title.SetTextColor(sdl.Color{0x06,0x84,0xdc,0xff})
	rack6Title.Move(0,460)
        serverpagerack.AddChild(rack6Title)

	rack6Desc:=sws.CreateTextAreaWidget(150,160,"This 4U HPC rack server, with 2 to 4 CPU and up to 8 RAM slots will put in the dust any cpu-intensive job")
	rack6Desc.SetReadonly(true)
	rack6Desc.SetFont(sws.LatoRegular14)
	rack6Desc.SetColor(0xffeeeeee)
	rack6Desc.Move(0,500)
	serverpagerack.AddChild(rack6Desc)

	rack6Button:=sws.CreateButtonWidget(100,25,"Configure >")
	rack6Button.SetColor(0xffeeeeee)
	rack6Button.Move(0,660)
	serverpagerack.AddChild(rack6Button)
	serverpagerack.configurerack6=rack6Button


	return serverpagerack
}

