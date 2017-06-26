package dctycoon

import(
	"time"
	"fmt"
)

type GameTimer struct {
	CurrentTime time.Time
	TimerClock  func()
}

var GlobalGameTimer *GameTimer

func GameTimerLoad(game map[string]interface{}) *GameTimer {
	var year, month, day int
	fmt.Sscanf(game["timer"].(string), "%d-%d-%d", &year, &month, &day)
	timer :=&GameTimer{
		CurrentTime: time.Date(year, time.Month(month), day, 0, 0, 0, 0,   time.UTC),
	}
	timer.TimerClock=func() { 
		timer.CurrentTime=timer.CurrentTime.Add(24*time.Hour)
		// test of GlobalEventPublisher.Publish
		//GlobalEventPublisher.Publish(fmt.Sprintf("%d-%d-%d",timer.CurrentTime.Year(),timer.CurrentTime.Month(),timer.CurrentTime.Day()),"long title")
	}
	return timer
}

func (self *GameTimer) Save() string {
	return fmt.Sprintf(`{"timer": "%d-%d-%d"}`,self.CurrentTime.Year(),self.CurrentTime.Month(),self.CurrentTime.Day())
}

