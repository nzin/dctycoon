package timer

import (
	"fmt"
	"github.com/google/btree"
	"time"
)

type GamerTimerEvent struct {
	Id      int32
	Date    time.Time
	Trigger func()
}

func (self *GamerTimerEvent) Less(b btree.Item) bool {
	bevt := b.(*GamerTimerEvent)
	if self.Date.Equal(bevt.Date) {
		return self.Id < bevt.Id
	}
	return self.Date.Before(bevt.Date)
}

type GameTimer struct {
	autoinc     int32
	CurrentTime time.Time // current day
	TimerClock  func()    // method to call when we switch to a new day
	events      *btree.BTree
}

var GlobalGameTimer *GameTimer

func NewGameTimer() *GameTimer {
	timer := &GameTimer{
		autoinc:     0,
		CurrentTime: time.Date(1990, time.Month(01), 01, 0, 0, 0, 0, time.UTC),
		events:      btree.New(10),
	}
	timer.TimerClock = func() {
		timer.CurrentTime = timer.CurrentTime.Add(24 * time.Hour)
		// trigger all events that are <= timer.CurrentTime
		for ev := timer.events.Min(); ev != nil; ev = timer.events.Min() {
			e := ev.(*GamerTimerEvent)
			if e.Date.After(timer.CurrentTime) {
				break
			} else {
				e.Trigger()
				timer.events.DeleteMin()
			}
		}
		// test of GlobalEventPublisher.Publish
		//GlobalEventPublisher.Publish(fmt.Sprintf("%d-%d-%d",timer.CurrentTime.Year(),timer.CurrentTime.Month(),timer.CurrentTime.Day()),"long title")
	}
	return timer
}

func (self *GameTimer) Load(game map[string]interface{}) {
	var year, month, day int
	fmt.Sscanf(game["timer"].(string), "%d-%d-%d", &year, &month, &day)
	self.CurrentTime = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	self.events = btree.New(10)
}

func (self *GameTimer) AddEvent(evdate time.Time, callback func()) {
	if evdate.Before(self.CurrentTime) {
		return
	}
	self.events.ReplaceOrInsert(&GamerTimerEvent{
		Id:      self.autoinc,
		Date:    evdate,
		Trigger: callback,
	})
	self.autoinc++
}

func (self *GameTimer) Save() string {
	return fmt.Sprintf(`{"timer": "%d-%d-%d"}`, self.CurrentTime.Year(), self.CurrentTime.Month(), self.CurrentTime.Day())
}
