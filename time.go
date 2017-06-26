package dctycoon

import(
	"time"
	"fmt"
)

type Timer struct {
	CurrentTime time.Time
	TimerClock  func()
}

var GlobalTimer *Timer

func TimerLoad(game map[string]interface{}) *Timer {
	var year, month, day int
	fmt.Sscanf(game["timer"].(string), "%d-%d-%d", &year, &month, &day)
	timer :=&Timer{
		CurrentTime: time.Date(year, time.Month(month), day, 0, 0, 0, 0,   time.UTC),
	}
	timer.TimerClock=func() { 
		timer.CurrentTime=timer.CurrentTime.Add(24*time.Hour)
	}
	return timer
}

func (self *Timer) Save() string {
	return fmt.Sprintf(`{"timer": "%d-%d-%d"}`,self.CurrentTime.Year(),self.CurrentTime.Month(),self.CurrentTime.Day())
}

