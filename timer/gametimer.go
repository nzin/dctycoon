package timer

import(
	"time"
	"fmt"
	"github.com/google/btree"
)

type GamerTimerEvent struct {
	Date    time.Time
	Trigger func()
}

func (self *GamerTimerEvent) Less(b btree.Item) bool {
	return self.Date.Before(b.(*GamerTimerEvent).Date)
}


type GameTimer struct {
	CurrentTime time.Time // current day
	TimerClock  func() // method to call when we switch to a new day
	events      *btree.BTree
}

var GlobalGameTimer *GameTimer

func GameTimerLoad(game map[string]interface{}) *GameTimer {
	var year, month, day int
	fmt.Sscanf(game["timer"].(string), "%d-%d-%d", &year, &month, &day)
	timer :=&GameTimer{
		CurrentTime: time.Date(year, time.Month(month), day, 0, 0, 0, 0,   time.UTC),
		events: btree.New(10),
	}
	timer.TimerClock=func() { 
		timer.CurrentTime=timer.CurrentTime.Add(24*time.Hour)
		// trigger all events that are <= timer.CurrentTime
		for ev:=timer.events.Min(); ev!=nil; ev=timer.events.Min() {
			e:=ev.(*GamerTimerEvent)
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

func (self *GameTimer) AddEvent(evdate time.Time, callback func()) {
	if evdate.Before(self.CurrentTime) {
		return
	}
	self.events.ReplaceOrInsert(&GamerTimerEvent{
		Date: evdate,
		Trigger: callback,
	})
}

func (self *GameTimer) Save() string {
	return fmt.Sprintf(`{"timer": "%d-%d-%d"}`,self.CurrentTime.Year(),self.CurrentTime.Month(),self.CurrentTime.Day())
}

