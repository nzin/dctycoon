package timer

import (
	"github.com/nzin/sws"
)

type EventMessageWidget struct {
	sws.CoreWidget
}

func NewEventMessageWidget(root *sws.RootWidget, longdesc string) *sws.MainWidget {
	widget := &EventMessageWidget{
		CoreWidget: *sws.NewCoreWidget(400, 100),
	}

	label := sws.NewLabelWidget(50, 50, "")
	label.SetImage("resources/icon-information-symbol.png")
	label.SetCentered(true)
	widget.AddChild(label)

	text := sws.NewTextAreaWidget(350, 60, longdesc)
	text.SetDisabled(true)
	text.Move(50, 10)
	widget.AddChild(text)

	closebutton := sws.NewButtonWidget(80, 25, "Close")
	closebutton.Move(300, 80)
	widget.AddChild(closebutton)

	mainwidget := sws.NewMainWidget(410, 150, "New event", false, false)
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
