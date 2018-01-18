package ui

import (
	"time"

	"github.com/nzin/dctycoon/timer"
	"github.com/nzin/sws"
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
	widget.SetColor(0xffffffff)
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
	for i := int32(0); i < self.nbmonths; i++ {
		xFrom := (self.Width() * i) / self.nbmonths
		xTo := (self.Width() * (i + 1)) / self.nbmonths

		data := self.data[i]
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
				self.FillRect(xFrom, self.Height()-(nbTo*self.Height()/total), xTo-xFrom, ((nbTo - nbFrom) * self.Height() / total), color)

				nbFrom = nbTo
			}
		}
	}
	self.SetDirtyFalse()
}
