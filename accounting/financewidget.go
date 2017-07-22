package accounting

import(
	"github.com/nzin/sws"
	//"github.com/veandco/go-sdl2/sdl"
)

type FinanceLine struct {
	Title *sws.LabelWidget
	N     *sws.LabelWidget
	N1    *sws.LabelWidget
}

//
// base widget for all finance sheet widgets
//
type FinanceWidget struct {
	sws.CoreWidget
	lines map[string]FinanceLine
	y     int32
	yearN  *sws.LabelWidget
	yearN1 *sws.LabelWidget
}

func (self *FinanceWidget) addLine(name, title string) {
	line:=FinanceLine{
		Title: sws.NewLabelWidget(190,25,title),
		N: sws.NewLabelWidget(90,25,"0 $"),
		N1: sws.NewLabelWidget(90,25,"0 $"),
	}
	line.Title.Move(10,self.y)
	self.AddChild(line.Title)
	line.N.Move(200,self.y)
	self.AddChild(line.N)
	line.N1.Move(300,self.y)
	self.AddChild(line.N1)
	self.lines[name]=line
	self.y+=25
	self.Resize(400,self.y)
}

func (self *FinanceWidget) addSeparator() {
	hr:=sws.NewHr(390)
	self.AddChild(hr)
	hr.Move(10,self.y)
	self.y+=5
	self.Resize(400,self.y)
}

func (self *FinanceWidget) addCategory(category string) {
	label:=sws.NewLabelWidget(190,25,category)
	label.Move(5,self.y)
	self.AddChild(label)
	self.y+=25
	self.Resize(400,self.y)
}

func NewFinanceWidget() *FinanceWidget {
	widget:=&FinanceWidget{
		CoreWidget: *sws.NewCoreWidget(400,50),
		lines: make(map[string]FinanceLine),
		y:     25,
		yearN: sws.NewLabelWidget(100,25,"Year N (forecast)"),
		yearN1: sws.NewLabelWidget(100,25,"Year N-1"),
	}
	widget.yearN.Move(200,0)
	widget.AddChild(widget.yearN)
	widget.yearN1.Move(300,0)
	widget.AddChild(widget.yearN1)
	
	return widget
}
