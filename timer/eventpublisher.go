package timer

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/nzin/sws"
)

type EventPublished struct {
	sws.CoreWidget
	eventpub      *EventPublisher
	messagewidget *sws.MainWidget
	shortdesc     string
	longdesc      string
	fadeintime    int32 // remaining second before appearing
	staytime      int32 // remaining second before fadeout
	fadeouttime   int32 // fadeout value before removing
	te            *sws.TimerEvent
}

func NewEventPublished(shortdesc string, longdesc string, eventpub *EventPublisher, pos int32) *EventPublished {
	log.Debug("NewEventPublished(", shortdesc, ",", longdesc, ",", eventpub, ",", pos, ")")
	corewidget := sws.NewCoreWidget(300, 30)
	widget := &EventPublished{
		CoreWidget:  *corewidget,
		eventpub:    eventpub,
		shortdesc:   shortdesc,
		longdesc:    longdesc,
		staytime:    4 * 40, // 4 seconds
		fadeouttime: 40,
		te:          nil,
	}
	flat := sws.NewFlatButtonWidget(300, 30, shortdesc)
	widget.AddChild(flat)

	flat.SetClicked(func() {
		if widget.messagewidget == nil {
			widget.messagewidget = NewEventMessageWidget(eventpub.root, longdesc)
		}
	})

	widget.Move(eventpub.root.Width()-300, eventpub.root.Height()-30-30*pos)
	widget.te = sws.TimerAddEvent(time.Now(), 25*time.Millisecond, func(evt *sws.TimerEvent) {
		if widget.fadeintime > 0 {
			widget.fadeintime--
			if widget.fadeintime == 0 {
			}
		} else if widget.staytime > 0 {
			widget.staytime--
		} else if widget.fadeouttime > 0 {
			widget.fadeouttime--
			widget.SetAlphaMod(uint8(255 * widget.fadeouttime / 40))
			widget.PostUpdate()
		} else {
			evt.StopRepeat()
			eventpub.remove(widget)
		}
	})

	return widget
}

type EventPublisherService interface {
	Publish(shortdesc string, longdesc string)
}

type EventPublisher struct {
	root   *sws.RootWidget
	events map[*EventPublished]int32
}

var GlobalEventPublisher *EventPublisher

func (self *EventPublisher) Publish(shortdesc string, longdesc string) {
	log.Debug("EventPublisher::Publish(", shortdesc, ",", longdesc, ")")
	var pos int32
	pos = 0
	for i := 0; i < len(self.events); i++ {
		for _, value := range self.events {
			if value == pos {
				pos++
			}
		}
	}

	ev := NewEventPublished(shortdesc, longdesc, self, pos)
	self.events[ev] = pos
	self.root.AddChild(ev)
}

func (self *EventPublisher) remove(event *EventPublished) {
	self.root.RemoveChild(event)
	delete(self.events, event)
}

func NewEventPublisher(root *sws.RootWidget) *EventPublisher {
	log.Debug("NewEventPublisher(", root, ")")
	return &EventPublisher{
		root:   root,
		events: make(map[*EventPublished]int32),
	}
}
