package timer

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGameTimer(t *testing.T) {
	passfebruary := false

	gt := NewGameTimer()
	gt.AddEvent(time.Date(1997, 2, 1, 1, 1, 1, 0, time.UTC), func() {
		passfebruary = true
	})
	gt.CurrentTime = time.Date(1997, 1, 30, 1, 1, 1, 0, time.UTC)
	// go to 31/1, nothing happens
	gt.TimerClock()
	assert.Equal(t, false, passfebruary, "31/1, nothing happens yet")
	// go to 1/2, we are un february
	gt.TimerClock()
	assert.Equal(t, true, passfebruary, "new month")

	// skip some days...
	passmarch := false
	gt.AddEvent(time.Date(1997, 3, 1, 1, 1, 1, 0, time.UTC), func() {
		passmarch = true
	})
	gt.CurrentTime = time.Date(1997, 3, 3, 1, 1, 1, 0, time.UTC)
	gt.TimerClock()
	assert.Equal(t, true, passmarch, "new month")
}

func TestCronGameTimer(t *testing.T) {
	month := 1

	gt := NewGameTimer()
	evt := gt.AddCron(1, -1, -1, func() {
		month++
	})
	gt.CurrentTime = time.Date(1990, 1, 30, 1, 1, 1, 0, time.UTC)
	// go to 31/1, nothing happens
	gt.TimerClock()
	assert.Equal(t, 1, month, "31/1, nothing happens yet")
	// go to 1/2, we are un february
	for i := 1; i <= 28; i++ {
		gt.TimerClock()
		assert.Equal(t, 2, month, fmt.Sprintf("new month (%d/2)", i))
	}
	gt.TimerClock()
	assert.Equal(t, 3, month, "1/3, beginning of march")

	// skip some days... we dont trigger cron, it is on purpose currently
	// because we normaly dont skips days :-)
	gt.CurrentTime = time.Date(1990, 4, 3, 1, 1, 1, 0, time.UTC)
	gt.TimerClock()
	assert.NotEqual(t, 4, month, "pass april")

	// unregister cron
	assert.Equal(t, 1, len(gt.cron), "before unregistering cron")
	gt.RemoveCron(evt)
	assert.Equal(t, 0, len(gt.cron), "unregister cron")
}
