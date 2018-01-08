package timer

import (
	"github.com/nzin/dctycoon/global"
	"github.com/nzin/sws"
	log "github.com/sirupsen/logrus"
)

type EventMessageWidget struct {
	sws.CoreWidget
}

func NewEventMessageWidget(root *sws.RootWidget, longdesc string) *sws.MainWidget {
	log.Debug("NewEventMessageWidget(", root, ",", longdesc, ")")
	widget := &EventMessageWidget{
		CoreWidget: *sws.NewCoreWidget(400, 100),
	}

	label := sws.NewLabelWidget(50, 50, "")
	if img, err := global.LoadImageAsset("assets/ui/icon-information-symbol.png"); err == nil {
		label.SetImageSurface(img)
	}
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
