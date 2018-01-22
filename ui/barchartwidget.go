package ui

import (
	"fmt"
	"time"

	"github.com/nzin/dctycoon/timer"
	"github.com/nzin/sws"
	"github.com/veandco/go-sdl2/sdl"
)

// BarChartWidget is a bar chart widget
// streamlining nb months of graph in a sliding window style
type BarChartWidget struct {
	sws.CoreWidget
	lastrefresh time.Time
	data        []int32
	chartcolor  uint32
}

// NewBarChartWidget create a simple timeline stacked barchart graph widget with
// - legend on the right
// - time point on the bottom
func NewBarChartWidget(w, h int32) *BarChartWidget {
	corewidget := sws.NewCoreWidget(w, h)
	widget := &BarChartWidget{
		CoreWidget:  *corewidget,
		lastrefresh: time.Now(),
		data:        make([]int32, w, w),
	}
	return widget
}

func (self *BarChartWidget) SetChartColor(color uint32) {
	self.chartcolor = color
}

// ChangeSpeed is part of GameTimerSubscriber interface
func (self *BarChartWidget) ChangeSpeed(speed int) {
}

// NewDay is part of GameTimerSubscriber interface
func (self *BarChartWidget) NewDay(timer *timer.GameTimer) {
	// we switch all datas
	for i := self.Width() - 1; i >= 1; i-- {
		self.data[i] = self.data[i-1]
	}
	self.data[0] = -1
	self.PostUpdate()

	self.lastrefresh = timer.CurrentTime
}

// ClearData removes all data (but not categories)
func (self *BarChartWidget) ClearData(t time.Time) {
	self.lastrefresh = t
	for i := int32(0); i < self.Width(); i++ {
		self.data[i] = -1
	}
	self.data[self.Width()-1] = 0
}

// AddPoint is really about adding/appending a new data point into the currenttimeline barchart
func (self *BarChartWidget) SetPoint(t time.Time, value int32) {
	if t.After(self.lastrefresh) {
		if value > self.data[0] {
			self.data[0] = value
		}
	} else {
		lowPoint := self.lastrefresh.AddDate(0, int(-self.Width()), 0)
		if t.After(lowPoint) {
			durationSince := int32(self.lastrefresh.Sub(t).Hours() / 24)
			if value > self.data[durationSince] {
				self.data[durationSince] = value
			}
		}
	}
	self.PostUpdate()
}

func (self *BarChartWidget) Repaint() {
	self.CoreWidget.Repaint()
	max := int32(0)
	for i := int32(0); i < self.Width(); i++ {
		data := self.data[i]
		if max < data {
			max = data
		}
	}

	width := self.Width() - 50
	height := self.Height() - 25
	xoffset := int32(50)
	previousValue := int32(0)
	if max > 0 {
		for i := xoffset; i < self.Width(); i++ {
			value := self.data[self.Width()-i]
			if value == -1 {
				value = previousValue
			}
			self.SetDrawColorHex(self.chartcolor)
			self.DrawLine(i, (max-value)*height/max, i, height)
			previousValue = value

			// write month
			currentday := self.lastrefresh.AddDate(0, 0, int(-self.Width()+i))
			if currentday.Day() == 1 && currentday.Month()%4 == 1 {
				month := currentday.Month()
				year := currentday.Year()
				self.SetDrawColorHex(0xff000000)
				self.DrawLine(i, height, i, height+4)
				self.WriteText(i-10, height+2, fmt.Sprintf("%d/%d", month, year%100), sdl.Color{0, 0, 0, 0xff})
			}
		}
	}
	if max < 10000 {
		self.WriteText(0, 0, fmt.Sprintf("%d", max), sdl.Color{0, 0, 0, 0xff})
	} else {
		self.WriteText(0, 0, fmt.Sprintf("%dk", max/1000), sdl.Color{0, 0, 0, 0xff})
	}
	self.SetDrawColorHex(0xff000000)
	self.DrawLine(xoffset, 0, xoffset, height)
	self.DrawLine(xoffset, height, width+xoffset, height)

	self.SetDirtyFalse()
}
