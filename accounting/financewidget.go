package accounting

import(
	"github.com/nzin/sws"
	//"github.com/veandco/go-sdl2/sdl"
)

type FinanceLine struct {
	Title *sws.SWS_Label
	N     *sws.SWS_Label
	N1    *sws.SWS_Label
}

//
// base widget for all finance sheet widgets
//
type FinanceWidget struct {
	sws.SWS_CoreWidget
	lines map[string]FinanceLine
	y     int32
	yearN  *sws.SWS_Label
	yearN1 *sws.SWS_Label
}

func (self *FinanceWidget) addLine(name, title string) {
	line:=FinanceLine{
		Title: sws.CreateLabel(190,25,title),
		N: sws.CreateLabel(90,25,"0 $"),
		N1: sws.CreateLabel(90,25,"0 $"),
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
	hr:=sws.CreateHr(390)
	self.AddChild(hr)
	hr.Move(10,self.y)
	self.y+=5
	self.Resize(400,self.y)
}

func (self *FinanceWidget) addCategory(category string) {
	label:=sws.CreateLabel(190,25,category)
	label.Move(5,self.y)
	self.AddChild(label)
	self.y+=25
	self.Resize(400,self.y)
}

func CreateFinanceWidget() *FinanceWidget {
	widget:=&FinanceWidget{
		SWS_CoreWidget: *sws.CreateCoreWidget(400,50),
		lines: make(map[string]FinanceLine),
		y:     25,
		yearN: sws.CreateLabel(100,25,"Year N"),
		yearN1: sws.CreateLabel(100,25,"Year N-1"),
	}
	widget.yearN.Move(200,0)
	widget.AddChild(widget.yearN)
	widget.yearN1.Move(300,0)
	widget.AddChild(widget.yearN1)
	
	return widget
}
