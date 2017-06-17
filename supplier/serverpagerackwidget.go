package supplier

import(
	"github.com/nzin/sws"
)

//
// Page Shop>>Explore>>Rack servers
//
type ServerPageRackWidget struct {
	sws.SWS_CoreWidget
	configurerack1 *sws.SWS_ButtonWidget
	configurerack2 *sws.SWS_ButtonWidget
	configurerack4 *sws.SWS_ButtonWidget
	configurerack6 *sws.SWS_ButtonWidget
}

func (self *ServerPageRackWidget) SetConfigureRack1Callback(callback func()) {
	self.configurerack1.SetClicked(callback)
}

func (self *ServerPageRackWidget) SetConfigureRack2Callback(callback func()) {
	self.configurerack2.SetClicked(callback)
}

func (self *ServerPageRackWidget) SetConfigureRack4Callback(callback func()) {
	self.configurerack4.SetClicked(callback)
}

func (self *ServerPageRackWidget) SetConfigureRack6Callback(callback func()) {
	self.configurerack6.SetClicked(callback)
}

func CreateServerPageRackWidget(width,height int32) *ServerPageRackWidget {
	serverpagerack:=&ServerPageRackWidget{
		SWS_CoreWidget: *sws.CreateCoreWidget(width,height),
	}
	serverpagerack.SetColor(0xffffffff)
	
        title:=sws.CreateLabel(200,20,"DEAL Rack Servers")
        title.SetFont(sws.LatoRegular20)
	title.SetColor(0xffffffff)
        title.Move(20,0)
        title.SetCentered(false)
        serverpagerack.AddChild(title)

	rack1Icon:=sws.CreateLabel(150,100,"")
	rack1Icon.SetImage("resources/server.1u0.png")
	rack1Icon.SetColor(0xffffffff)
        rack1Icon.SetCentered(true)
	rack1Icon.Move(0,20)
        serverpagerack.AddChild(rack1Icon)

	rack1Title:=sws.CreateLabel(150,20,"Rack R100 server")
	rack1Title.SetColor(0xffffffff)
	rack1Title.Move(0,120)
        serverpagerack.AddChild(rack1Title)

	rack1Button:=sws.CreateButtonWidget(100,25,"Configure >")
	rack1Button.SetColor(0xffffffff)
	rack1Button.Move(0,320)
	serverpagerack.AddChild(rack1Button)
	serverpagerack.configurerack1=rack1Button


	rack2Icon:=sws.CreateLabel(150,100,"")
	rack2Icon.SetImage("resources/server.2u0.png")
	rack2Icon.SetColor(0xffffffff)
        rack2Icon.SetCentered(true)
	rack2Icon.Move(150,20)
        serverpagerack.AddChild(rack2Icon)

	rack2Title:=sws.CreateLabel(150,20,"Rack R200 server")
	rack2Title.SetColor(0xffffffff)
	rack2Title.Move(150,120)
        serverpagerack.AddChild(rack2Title)

	rack2Button:=sws.CreateButtonWidget(100,25,"Configure >")
	rack2Button.SetColor(0xffffffff)
	rack2Button.Move(150,320)
	serverpagerack.AddChild(rack2Button)
	serverpagerack.configurerack2=rack2Button


	rack4Icon:=sws.CreateLabel(150,100,"")
	rack4Icon.SetImage("resources/server.4u0.png")
	rack4Icon.SetColor(0xffffffff)
        rack4Icon.SetCentered(true)
	rack4Icon.Move(300,20)
        serverpagerack.AddChild(rack4Icon)

	rack4Title:=sws.CreateLabel(150,20,"Rack R400 server")
	rack4Title.SetColor(0xffffffff)
	rack4Title.Move(300,120)
	serverpagerack.AddChild(rack4Title)

	rack4Button:=sws.CreateButtonWidget(100,25,"Configure >")
	rack4Button.SetColor(0xffffffff)
	rack4Button.Move(300,320)
	serverpagerack.AddChild(rack4Button)
	serverpagerack.configurerack4=rack4Button


	rack6Icon:=sws.CreateLabel(150,100,"")
	rack6Icon.SetImage("resources/server.4u0.png")
	rack6Icon.SetColor(0xffffffff)
        rack6Icon.SetCentered(true)
	rack6Icon.Move(0,360)
        serverpagerack.AddChild(rack6Icon)

	rack6Title:=sws.CreateLabel(150,20,"Rack R600 server")
	rack6Title.SetColor(0xffffffff)
	rack6Title.Move(0,460)
        serverpagerack.AddChild(rack6Title)

	rack6Button:=sws.CreateButtonWidget(100,25,"Configure >")
	rack6Button.SetColor(0xffffffff)
	rack6Button.Move(0,660)
	serverpagerack.AddChild(rack6Button)
	serverpagerack.configurerack6=rack6Button


	return serverpagerack
}

