package timer

import (
	"fmt"
	"time"

	"github.com/google/btree"
	log "github.com/sirupsen/logrus"
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

type GameCronEvent struct {
	day     int32 // -1 = '*'
	month   int32 // -1 = '*'
	year    int32 // -1 = '*'
	Trigger func()
}

type GameTimer struct {
	autoinc     int32
	CurrentTime time.Time // current day
	TimerClock  func()    // method to call when we switch to a new day
	events      *btree.BTree
	cron        []*GameCronEvent
}

//
// NewGameTimer create a timer object starting at 1/1/1990
func NewGameTimer() *GameTimer {
	log.Debug("NewGameTimer()")
	timer := &GameTimer{
		autoinc:     0,
		CurrentTime: time.Date(1995, time.Month(01), 01, 0, 0, 0, 0, time.UTC),
		events:      btree.New(10),
		cron:        make([]*GameCronEvent, 0),
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
		// check for cron
		for _, c := range timer.cron {
			if c.day == -1 || timer.CurrentTime.Day() == int(c.day) {
				if c.month == -1 || timer.CurrentTime.Month() == time.Month(c.month) {
					if c.year == -1 || timer.CurrentTime.Year() == int(c.year) {
						c.Trigger()
					}
				}
			}
		}
	}
	return timer
}

func (self *GameTimer) Load(game map[string]interface{}) {
	log.Debug("GameTimer::Load(", game, ")")
	var year, month, day int
	fmt.Sscanf(game["timer"].(string), "%d-%d-%d", &year, &month, &day)
	self.CurrentTime = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	self.events = btree.New(10)
}

func (self *GameTimer) AddEvent(evdate time.Time, callback func()) {
	log.Debug("GameTimer::AddEvent(", evdate, ",callback)")
	if evdate.Before(self.CurrentTime) || evdate.Equal(self.CurrentTime) {
		return
	}
	self.events.ReplaceOrInsert(&GamerTimerEvent{
		Id:      self.autoinc,
		Date:    evdate,
		Trigger: callback,
	})
	self.autoinc++
}

func (self *GameTimer) AddCron(day, month, year int32, callback func()) *GameCronEvent {
	log.Debug("GameTimer::AddCron(", day, ",", month, ",", year, ",callback)")
	cronevent := &GameCronEvent{
		day:     day,
		month:   month,
		year:    year,
		Trigger: callback,
	}
	self.cron = append(self.cron, cronevent)
	return cronevent
}

func (self *GameTimer) RemoveCron(evt *GameCronEvent) {
	log.Debug("GameTimer::RemoveCron(", evt, ")")
	for i, c := range self.cron {
		if c == evt {
			self.cron = append(self.cron[:i], self.cron[i+1:]...)
			return
		}
	}
}

func (self *GameTimer) Save() string {
	log.Debug("GameTimer::Save()")
	return fmt.Sprintf(`{"timer": "%d-%d-%d"}`, self.CurrentTime.Year(), self.CurrentTime.Month(), self.CurrentTime.Day())
}
