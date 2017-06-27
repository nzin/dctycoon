package timer

import(
	"github.com/nzin/sws"
)

type EventMessageWidget struct {
	sws.SWS_CoreWidget
}

func CreateEventMessageWidget(root *sws.SWS_RootWidget, longdesc string) *sws.SWS_MainWidget { 
	widget:=&EventMessageWidget{
		SWS_CoreWidget: *sws.CreateCoreWidget(400,100),
	}
	
	label:=sws.CreateLabel(50,50,"")
	label.SetImage("resources/icon-information-symbol.png")
	label.SetCentered(true)
	widget.AddChild(label)
	
	text:=sws.CreateTextAreaWidget(350,60,longdesc)
	text.SetReadonly(true)
	text.Move(50,10)
	widget.AddChild(text)
	
	closebutton:=sws.CreateButtonWidget(80,25,"Close")
	closebutton.Move(300,80)
	widget.AddChild(closebutton)
	
	mainwidget := sws.CreateMainWidget(410,150,"New event",false,false)
	mainwidget.SetCloseCallback(func() {
                root.RemoveChild(mainwidget)
        })
	
	closebutton.SetClicked(func() {
		root.RemoveChild(mainwidget)
	})
	
	mainwidget.SetInnerWidget(widget)
	
	root.AddChild(mainwidget)
	return mainwidget
}


