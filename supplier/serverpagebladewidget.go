package supplier

import(
	"github.com/nzin/sws"
)

type ServerPageBladeWidget struct {
	sws.SWS_CoreWidget
}

func CreateServerPageBladeWidget(width,height int32) *ServerPageBladeWidget {
	serverpageblade:=&ServerPageBladeWidget{
		SWS_CoreWidget: *sws.CreateCoreWidget(width,height),
	}
	serverpageblade.SetColor(0xffffffff)
	
        title:=sws.CreateLabel(200,20,"DEAL Blade Servers")
        title.SetFont(sws.LatoRegular20)
	title.SetColor(0xffffffff)
        title.Move(20,0)
        title.SetCentered(false)
        serverpageblade.AddChild(title)

	blade1Icon:=sws.CreateLabel(150,100,"")
	blade1Icon.SetImage("resources/server.blade.8u0.png")
	blade1Icon.SetColor(0xffffffff)
        blade1Icon.SetCentered(true)
	blade1Icon.Move(0,20)
        serverpageblade.AddChild(blade1Icon)

	blade1Title:=sws.CreateLabel(150,20,"Blade B100 server")
	blade1Title.SetColor(0xffffffff)
	blade1Title.Move(0,120)
        serverpageblade.AddChild(blade1Title)

	blade1Button:=sws.CreateButtonWidget(100,25,"Configure >")
	blade1Button.SetColor(0xffffffff)
	blade1Button.Move(0,320)
	serverpageblade.AddChild(blade1Button)


	blade2Icon:=sws.CreateLabel(150,100,"")
	blade2Icon.SetImage("resources/server.blade.8u0.png")
	blade2Icon.SetColor(0xffffffff)
        blade2Icon.SetCentered(true)
	blade2Icon.Move(150,20)
        serverpageblade.AddChild(blade2Icon)

	blade2Title:=sws.CreateLabel(150,20,"Blade B200 server")
	blade2Title.SetColor(0xffffffff)
	blade2Title.Move(150,120)
        serverpageblade.AddChild(blade2Title)

	blade2Button:=sws.CreateButtonWidget(100,25,"Configure >")
	blade2Button.SetColor(0xffffffff)
	blade2Button.Move(150,320)
	serverpageblade.AddChild(blade2Button)


	return serverpageblade
}

