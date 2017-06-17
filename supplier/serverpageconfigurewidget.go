package supplier

import(
	"github.com/nzin/sws"
)

//
// Page Shop>>Explore>>xxx servers>>Configure
//
type ServerPageConfigureWidget struct {
	sws.SWS_CoreWidget
	title         *sws.SWS_Label
	buybutton     *sws.SWS_ButtonWidget
	configureicon *sws.SWS_Label
	conftype      *ServerConfType
}

func (self *ServerPageConfigureWidget) SetBuyCallback(callback func()) {
	self.buybutton.SetClicked(callback)
}

func (self *ServerPageConfigureWidget) SetConfType(conftypename string) {
	for _,c := range(AvailableConfs) {
		if c.ServerName==conftypename {
			self.conftype=&c
		}
	}
	if self.conftype==nil { return }
	
	self.configureicon.SetImage(self.conftype.ServerSprite+"0.png")
	// todo
	sws.PostUpdate()
}

func CreateServerPageConfigureWidget(width,height int32) *ServerPageConfigureWidget {
	serverpageconfigure:=&ServerPageConfigureWidget{
		SWS_CoreWidget: *sws.CreateCoreWidget(width,height),
	}
	serverpageconfigure.SetColor(0xffffffff)
	
        title:=sws.CreateLabel(200,20,"Configure Server")
        title.SetFont(sws.LatoRegular20)
	title.SetColor(0xffffffff)
        title.Move(20,0)
        title.SetCentered(false)
        serverpageconfigure.AddChild(title)
        serverpageconfigure.title=title

	configureIcon:=sws.CreateLabel(150,100,"")
	configureIcon.SetColor(0xffffffff)
        configureIcon.SetCentered(true)
	configureIcon.Move(0,20)
        serverpageconfigure.AddChild(configureIcon)
        serverpageconfigure.configureicon=configureIcon



	buyButton:=sws.CreateButtonWidget(100,25,"Buy >")
	buyButton.SetColor(0xffffffff)
	buyButton.Move(0,320)
	serverpageconfigure.AddChild(buyButton)
	serverpageconfigure.buybutton=buyButton


	return serverpageconfigure
}

