package dctycoon

import (
	"fmt"
	"strconv"
	"time"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/dctycoon/ui"
	"github.com/nzin/sws"

	log "github.com/sirupsen/logrus"
)

const (
	COLOR_SALE_UNASSIGNED = 0xffaaaaaa
	COLOR_SALE_YOU        = 0xff8888ff
	COLOR_SALE_COMPETITOR = 0xffff7b44
)

//
// MainStatsWidget is a all stat main widget:
// - opponent view
// - customer demand stats
// - ...
type MainStatsWidget struct {
	rootwindow      *sws.RootWidget
	mainwidget      *sws.MainWidget
	tabwidget       *sws.TabWidget
	opponentswidget *OpponentStatWidget
	demandswidget   *DemandStatWidget
	game            *Game
}

func (self *MainStatsWidget) Show() {
	self.rootwindow.AddChild(self.mainwidget)
	self.rootwindow.SetFocus(self.mainwidget)
	self.tabwidget.SelectTab(0)
}

func (self *MainStatsWidget) Hide() {
	self.rootwindow.RemoveChild(self.mainwidget)
	children := self.rootwindow.GetChildren()
	if len(children) > 0 {
		self.rootwindow.SetFocus(children[0])
	}
}

// NewMainStatsWidget presents different stats and graphs
func NewMainStatsWidget(root *sws.RootWidget, g *Game) *MainStatsWidget {
	mainwidget := sws.NewMainWidget(850, 600, " Graph and Statistics ", true, true)
	mainwidget.Center(root)

	widget := &MainStatsWidget{
		rootwindow:      root,
		mainwidget:      mainwidget,
		tabwidget:       sws.NewTabWidget(200, 200),
		game:            g,
		opponentswidget: NewOpponentStatWidget(200, 200),
		demandswidget:   NewDemandStatWidget(200, 200, g),
	}

	widget.mainwidget.SetInnerWidget(widget.tabwidget)

	widget.mainwidget.SetCloseCallback(func() {
		widget.Hide()
	})

	widget.tabwidget.AddTab("Competitors", widget.opponentswidget)
	widget.tabwidget.AddTab("Customer demands", widget.demandswidget)

	return widget
}

func (self *MainStatsWidget) SetGame() {
	self.opponentswidget.SetGame(self.game)
	self.demandswidget.SetGame(self.game.timer.CurrentTime, self.game.GetGameStats())
}

func (self *MainStatsWidget) LoadGame() {
	self.opponentswidget.SetGame(self.game)
	self.demandswidget.SetGame(self.game.timer.CurrentTime, self.game.GetGameStats())
}

//
// OpponentStatWidgetLine show 1 opponent summary info
// see OpponentStatWidget
type OpponentStatWidgetLine struct {
	sws.CoreWidget
	opponent  *NPDatacenter
	picture   *sws.LabelWidget
	name      *sws.LabelWidget
	location  *sws.LabelWidget
	nbservers *sws.LabelWidget
}

func NewOpponentStatWidgetLine(opponent *NPDatacenter) *OpponentStatWidgetLine {
	corewidget := sws.NewCoreWidget(400, 100)
	line := &OpponentStatWidgetLine{
		CoreWidget: *corewidget,
		opponent:   opponent,
		picture:    sws.NewLabelWidget(100, 100, ""),
		name:       sws.NewLabelWidget(300, 25, opponent.GetName()),
		location:   sws.NewLabelWidget(300, 25, opponent.location.Name),
		nbservers:  sws.NewLabelWidget(300, 25, fmt.Sprintf("%d", len(opponent.GetInventory().Items))),
	}
	if opponent.GetPicture() != "" {
		if surface, err := global.LoadImageAsset("assets/faces/" + opponent.GetPicture()); err == nil {
			if adjusted, err := global.AdjustImage(surface, 100, 100); err == nil {
				line.picture.SetImageSurface(adjusted)
			}
		}
	}
	line.AddChild(line.picture)

	labelname := sws.NewLabelWidget(100, 25, "Name: ")
	labelname.Move(100, 5)
	line.AddChild(labelname)

	line.name.Move(200, 5)
	line.AddChild(line.name)

	labellocation := sws.NewLabelWidget(100, 25, "Location: ")
	labellocation.Move(100, 30)
	line.AddChild(labellocation)

	line.location.Move(200, 30)
	line.AddChild(line.location)

	labelNbServers := sws.NewLabelWidget(100, 25, "Nb servers: ")
	labelNbServers.Move(100, 55)
	line.AddChild(labelNbServers)

	line.nbservers.Move(200, 55)
	line.AddChild(line.nbservers)

	opponent.GetInventory().AddInventorySubscriber(line)

	return line
}

func (self *OpponentStatWidgetLine) ItemInTransit(*supplier.InventoryItem) {}
func (self *OpponentStatWidgetLine) ItemInStock(*supplier.InventoryItem) {
	self.nbservers.SetText(fmt.Sprintf("%d", len(self.opponent.GetInventory().Items)))
}
func (self *OpponentStatWidgetLine) ItemRemoveFromStock(*supplier.InventoryItem) {}
func (self *OpponentStatWidgetLine) ItemInstalled(*supplier.InventoryItem)       {}
func (self *OpponentStatWidgetLine) ItemUninstalled(*supplier.InventoryItem)     {}
func (self *OpponentStatWidgetLine) ItemChangedPool(*supplier.InventoryItem)     {}

//
// OpponentStatWidget is an opponent stat view
type OpponentStatWidget struct {
	sws.CoreWidget
	vbox   *sws.VBoxWidget
	scroll *sws.ScrollWidget
}

func NewOpponentStatWidget(w, h int32) *OpponentStatWidget {
	corewidget := sws.NewCoreWidget(w, h)
	widget := &OpponentStatWidget{
		CoreWidget: *corewidget,
		vbox:       sws.NewVBoxWidget(400, 0),
		scroll:     sws.NewScrollWidget(w, h),
	}
	widget.scroll.SetInnerWidget(widget.vbox)
	widget.AddChild(widget.scroll)

	return widget
}

func (self *OpponentStatWidget) Resize(w, h int32) {
	self.CoreWidget.Resize(w, h)
	self.scroll.Resize(w, h)
}

func (self *OpponentStatWidget) SetGame(game *Game) {
	log.Debug("OpponentStatWidget::SetGame(", game, ")")
	self.vbox.RemoveAllChildren()
	for _, o := range game.GetNPActors() {
		line := NewOpponentStatWidgetLine(o)
		self.vbox.AddChild(line)
	}
	self.scroll.Resize(self.Width(), self.Height())
}

// DemandStatDetailsWidget is used to give some insight/details on
// on particular demand stat
type DemandStatDetailsWidget struct {
	sws.CoreWidget
	yoffset int32
}

func NewDemandStatDetailsWidget(bgcolor uint32, stat *DemandStat) *DemandStatDetailsWidget {
	corewidget := sws.NewCoreWidget(525, 25)

	line := &DemandStatDetailsWidget{
		CoreWidget: *corewidget,
	}
	line.SetColor(bgcolor)

	for _, s := range stat.serverdemands {
		nb := sws.NewLabelWidget(50, 25, fmt.Sprintf("%d x", s.nb))
		nb.Move(25, line.yoffset)
		nb.SetColor(bgcolor)
		line.AddChild(nb)

		ram := sws.NewLabelWidget(150, 25, fmt.Sprintf("min ram = "+global.AdjustMega(s.ramsize)))
		ram.Move(75, line.yoffset)
		ram.SetColor(bgcolor)
		line.AddChild(ram)

		disk := sws.NewLabelWidget(150, 25, fmt.Sprintf("min disk = "+global.AdjustMega(s.disksize)))
		disk.Move(225, line.yoffset)
		disk.SetColor(bgcolor)
		line.AddChild(disk)

		cpu := sws.NewLabelWidget(150, 25, fmt.Sprintf("min nb cpus = %d", s.nbcores))
		cpu.Move(375, line.yoffset)
		cpu.SetColor(bgcolor)
		line.AddChild(cpu)

		line.yoffset += 25
	}
	return line
}

//
// DemandStatWidget is a customer demand stat widget
type DemandStatWidget struct {
	sws.CoreWidget
	demandstats *ui.TableWithDetails
	barchart    *ui.BarChartWidget
}

func NewDemandStatWidget(w, h int32, g *Game) *DemandStatWidget {
	corewidget := sws.NewCoreWidget(w, h)
	widget := &DemandStatWidget{
		CoreWidget:  *corewidget,
		demandstats: ui.NewTableWithDetails(525+15, 200),
		barchart:    ui.NewBarChartWidget(18, 525, 200),
	}
	widget.barchart.AddCategory("you", COLOR_SALE_YOU)
	widget.barchart.AddCategory("unassigned", COLOR_SALE_UNASSIGNED)
	widget.barchart.AddCategory("competitor", COLOR_SALE_COMPETITOR)
	g.AddGameTimerSubscriber(widget.barchart)
	widget.AddChild(widget.barchart)

	widget.demandstats.Move(0, 220)
	widget.demandstats.AddHeader("Date", 100, func(l1, l2 string) bool {
		var d1, d2, m1, m2, y1, y2 int
		fmt.Sscanf(l1, "%d-%d-%d", &d1, &m1, &y1)
		fmt.Sscanf(l2, "%d-%d-%d", &d2, &m2, &y2)
		if y1 != y2 {
			return y1 < y2
		}
		if m1 != m2 {
			return m1 < m2
		}
		return d1 < d2
	})
	widget.demandstats.AddHeader("Price", 100, func(l1, l2 string) bool {
		var f1, f2 float64
		fmt.Sscanf(l1, "%f $", &f1)
		fmt.Sscanf(l2, "%f $", &f2)
		return f1 < f2
	})
	widget.demandstats.AddHeader("Nb servers", 100, func(l1, l2 string) bool {
		nb1, _ := strconv.Atoi(l1)
		nb2, _ := strconv.Atoi(l2)
		return nb1 < nb2
	})
	widget.demandstats.AddHeader("Buyer", 200, func(l1, l2 string) bool { return l1 < l2 })
	widget.AddChild(widget.demandstats)

	g.GetGameStats().AddDemandStatSubscriber(widget)
	return widget
}

// NewDemandStat is part of DemandStatSubscriber interface
func (self *DemandStatWidget) NewDemandStat(ds *DemandStat) {
	log.Debug("DemandStatWidget::NewDemandStat(", ds, ")")
	nbservers := int32(0)
	for _, s := range ds.serverdemands {
		nbservers += s.nb
	}
	price := "-"
	bgcolor := uint32(COLOR_SALE_UNASSIGNED)
	if ds.buyer != "" {
		price = fmt.Sprintf("%.2f $", ds.price)
		bgcolor = COLOR_SALE_COMPETITOR
		if ds.buyer == "you" {
			bgcolor = COLOR_SALE_YOU
		}
	}

	labels := []string{
		fmt.Sprintf("%d-%d-%d", ds.date.Day(), ds.date.Month(), ds.date.Year()),
		price,
		fmt.Sprintf("%d", nbservers),
		ds.buyer,
	}
	line := ui.NewTableWithDetailsRow(bgcolor, labels, NewDemandStatDetailsWidget(bgcolor, ds))
	self.demandstats.AddRowTop(line)

	categoryname := "competitor"
	switch ds.buyer {
	case "":
		categoryname = "unassigned"
	case "you": // defined in player.go
		categoryname = "you"
	}
	self.barchart.AddPoint(ds.date, categoryname)
}

func (self *DemandStatWidget) Resize(w, h int32) {
	self.CoreWidget.Resize(w, h)
}

func (self *DemandStatWidget) SetGame(t time.Time, gamestats *GameStats) {
	log.Debug("DemandStatWidget::SetGame(", gamestats, ")")
	self.demandstats.ClearRows()
	self.barchart.Clear(t)

	for _, ds := range gamestats.demandsstats {
		self.NewDemandStat(ds)
	}
}
