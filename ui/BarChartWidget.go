package ui

import (
	"fmt"
	"time"

	"github.com/nzin/dctycoon/timer"
	"github.com/nzin/sws"
	"github.com/veandco/go-sdl2/sdl"
)

type BarChartCategory struct {
	name  string
	color uint32
}

// BarChartWidget is a bar chart widget
// streamlining nb months of graph in a sliding window style
type BarChartWidget struct {
	sws.CoreWidget
	nbmonths    int32
	lastrefresh time.Time
	data        []map[string]int32 // category -> nb
	categories  []*BarChartCategory
}

func NewBarChartWidget(nbmonths int32, w, h int32) *BarChartWidget {
	corewidget := sws.NewCoreWidget(w, h)
	widget := &BarChartWidget{
		CoreWidget:  *corewidget,
		nbmonths:    nbmonths,
		lastrefresh: time.Now(),
		data:        make([]map[string]int32, nbmonths, nbmonths),
		categories:  make([]*BarChartCategory, 0, 0),
	}
	for i := int32(0); i < nbmonths; i++ {
		widget.data[i] = make(map[string]int32)
	}
	return widget
}

// ChangeSpeed is part of GameTimerSubscriber interface
func (self *BarChartWidget) ChangeSpeed(speed int) {
}

// NewDay is part of GameTimerSubscriber interface
func (self *BarChartWidget) NewDay(timer *timer.GameTimer) {
	if self.lastrefresh.Year() != timer.CurrentTime.Year() ||
		self.lastrefresh.Month() != timer.CurrentTime.Month() {
		// we changed from month, we switch all datas
		for i := self.nbmonths - 1; i >= 1; i-- {
			self.data[i] = self.data[i-1]
		}
		self.data[0] = make(map[string]int32)
		self.PostUpdate()
	}
	self.lastrefresh = timer.CurrentTime
}

func (self *BarChartWidget) Clear(t time.Time) {
	self.lastrefresh = t
	for i := int32(0); i < self.nbmonths; i++ {
		self.data[i] = make(map[string]int32)
	}
}

func (self *BarChartWidget) AddCategory(name string, color uint32) {
	category := &BarChartCategory{
		name:  name,
		color: color,
	}
	self.categories = append(self.categories, category)
}

func (self *BarChartWidget) AddPoint(t time.Time, category string) {
	pointMonth := t.Year()*12 + int(t.Month())
	currentMonth := self.lastrefresh.Year()*12 + int(self.lastrefresh.Month())
	diff := int32(currentMonth - pointMonth)
	if diff >= 0 && diff < self.nbmonths {
		self.data[diff][category]++
	}
	self.PostUpdate()
}

func (self *BarChartWidget) Repaint() {
	self.CoreWidget.Repaint()
	max := int32(0)
	for i := int32(0); i < self.nbmonths; i++ {
		data := self.data[i]
		total := int32(0)
		for j := 0; j < len(self.categories); j++ {
			total += data[self.categories[j].name]
		}
		if max < total {
			max = total
		}
	}

	width := self.Width() - 25
	if self.Width() > 400 {
		width = self.Width() - 150 - 25
	}
	height := self.Height() - 25
	xoffset := int32(25)
	for i := int32(0); i < self.nbmonths; i++ {
		xFrom := (width*i)/self.nbmonths + xoffset
		xTo := (width*(i+1))/self.nbmonths + xoffset

		data := self.data[self.nbmonths-1-i]
		total := int32(0)
		for j := 0; j < len(self.categories); j++ {
			total += data[self.categories[j].name]
		}
		if total > 0 {
			nbFrom := int32(0)
			for j := 0; j < len(self.categories); j++ {
				color := self.categories[j].color
				nbTo := data[self.categories[j].name] + nbFrom
				//				fmt.Println(i, xFrom, self.Height()-(nbTo*self.Height()/total), xTo-xFrom, ((nbTo - nbFrom) * self.Height() / total), color)
				self.FillRect(xFrom+1, height-(nbTo*height/max), xTo-xFrom-2, ((nbTo - nbFrom) * height / max), color)

				nbFrom = nbTo
			}
		}
		// write month
		if i%4 == 0 {
			month := self.lastrefresh.Month()
			year := self.lastrefresh.Year()
			for j := self.nbmonths - 1; j > i; j-- {
				if month == 1 {
					year--
					month = 12
				} else {
					month--
				}
			}
			self.SetDrawColorHex(0xff000000)
			self.DrawLine((xFrom+xTo)/2, height, (xFrom+xTo)/2, height+4)
			self.WriteText(xFrom-10, height+2, fmt.Sprintf("%d/%d", month, year%100), sdl.Color{0, 0, 0, 0xff})
		}
	}
	self.WriteText(0, 0, fmt.Sprintf("%d", max), sdl.Color{0, 0, 0, 0xff})
	self.SetDrawColorHex(0xff000000)
	self.DrawLine(25, 0, 25, height)
	self.DrawLine(25, height, width+25, height)

	// show labels
	if self.Width() > 400 {
		xoffset = self.Width() - 150
		for j := 0; j < len(self.categories); j++ {
			self.FillRect(xoffset+5, int32(j*25)+5, 15, 15, self.categories[j].color)
			self.WriteText(xoffset+25, int32(j*25), self.categories[j].name, sdl.Color{0, 0, 0, 0xff})
		}
	}
	self.SetDirtyFalse()
}
