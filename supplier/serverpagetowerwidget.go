package supplier

import(
	"github.com/nzin/sws"
)

//
// Page Shop>>Explore>>Tower servers
//
type ServerPageTowerWidget struct {
	sws.SWS_CoreWidget
	configuretower1 *sws.SWS_ButtonWidget
}

func (self *ServerPageTowerWidget) SetConfigureTower1Callback(callback func()) {
	self.configuretower1.SetClicked(callback)
}


func CreateServerPageTowerWidget(width,height int32) *ServerPageTowerWidget {
	serverpagetower:=&ServerPageTowerWidget{
		SWS_CoreWidget: *sws.CreateCoreWidget(width,height),
	}
	serverpagetower.SetColor(0xffffffff)
	
        title:=sws.CreateLabel(200,20,"DEAL Tower Servers")
        title.SetFont(sws.LatoRegular20)
	title.SetColor(0xffffffff)
        title.Move(20,0)
        title.SetCentered(false)
        serverpagetower.AddChild(title)

	towerIcon:=sws.CreateLabel(150,100,"")
	towerIcon.SetImage("resources/tower0.png")
	towerIcon.SetColor(0xffffffff)
        towerIcon.SetCentered(true)
	towerIcon.Move(0,20)
        serverpagetower.AddChild(towerIcon)

	towerTitle:=sws.CreateLabel(150,20,"Tower T1000")
	towerTitle.SetColor(0xffffffff)
	towerTitle.Move(0,120)
        serverpagetower.AddChild(towerTitle)

	towerButton:=sws.CreateButtonWidget(100,25,"Configure >")
	towerButton.SetColor(0xffffffff)
	towerButton.Move(0,320)
	serverpagetower.AddChild(towerButton)
	serverpagetower.configuretower1=towerButton


	return serverpagetower
}

