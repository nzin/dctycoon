package dctycoon

import (
	"fmt"
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
	rootwindow        *sws.RootWidget
	mainwidget        *sws.MainWidget
	tabwidget         *sws.TabWidget
	opponentswidget   *OpponentStatWidget
	playerwidget      *PlayerStatWidget
	demandswidget     *DemandStatWidget
	globalpowerwidget *PowerStatWidget
	game              *Game
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
		rootwindow:        root,
		mainwidget:        mainwidget,
		tabwidget:         sws.NewTabWidget(200, 200),
		game:              g,
		opponentswidget:   NewOpponentStatWidget(200, 200),
		playerwidget:      NewPlayerStatWidget(200, 200, g),
		demandswidget:     NewDemandStatWidget(200, 200, g),
		globalpowerwidget: NewPowerStatWidget(200, 200, g),
	}

	widget.mainwidget.SetInnerWidget(widget.tabwidget)

	widget.mainwidget.SetCloseCallback(func() {
		widget.Hide()
	})

	widget.tabwidget.AddTab("Competitors", widget.opponentswidget)
	widget.tabwidget.AddTab("You", widget.playerwidget)
	widget.tabwidget.AddTab("Customer demands", widget.demandswidget)
	widget.tabwidget.AddTab("Global Power", widget.globalpowerwidget)
	if g.GetDebug() {
		widget.tabwidget.AddTab("Opponent debug", NewDebugOpponentWidget(200, 200, g))
	}
	return widget
}

func (self *MainStatsWidget) SetGame() {
	self.opponentswidget.SetGame(self.game)
	self.playerwidget.SetGame(self.game.GetGameStats())
	self.demandswidget.SetGame(self.game.timer.CurrentTime, self.game.GetGameStats())
	self.globalpowerwidget.SetGame(self.game.timer.CurrentTime, self.game.GetGameStats())
}

//
// OpponentStatWidgetLine show 1 opponent summary info
// see OpponentStatWidget
type OpponentStatWidgetLine struct {
	sws.CoreWidget
	opponent   *NPDatacenter
	picture    *sws.LabelWidget
	name       *sws.LabelWidget
	location   *sws.LabelWidget
	nbservers  *sws.LabelWidget
	reputation *sws.LabelWidget
}

func NewOpponentStatWidgetLine(opponent *NPDatacenter) *OpponentStatWidgetLine {
	corewidget := sws.NewCoreWidget(400, 100)
	line := &OpponentStatWidgetLine{
		CoreWidget: *corewidget,
		opponent:   opponent,
		picture:    sws.NewLabelWidget(100, 100, ""),
		name:       sws.NewLabelWidget(300, 25, opponent.GetName()),
		location:   sws.NewLabelWidget(300, 25, opponent.GetLocation().Name),
		nbservers:  sws.NewLabelWidget(300, 25, fmt.Sprintf("%d", len(opponent.GetInventory().Items))),
		reputation: sws.NewLabelWidget(300, 25, fmt.Sprintf("%.2f", opponent.GetReputationScore())),
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

	labelReputation := sws.NewLabelWidget(100, 25, "Reputation: ")
	labelReputation.Move(100, 80)
	line.AddChild(labelReputation)

	line.reputation.Move(200, 80)
	line.AddChild(line.reputation)

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
	corewidget := sws.NewCoreWidget(525+15, 0)

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
	line.Resize(525, line.yoffset)
	return line
}

//
// DemandStatWidget is a customer demand stat widget
type DemandStatWidget struct {
	sws.CoreWidget
	demandstats *ui.TableWithDetails
	barchart    *ui.StackedBarChartWidget
}

func NewDemandStatWidget(w, h int32, g *Game) *DemandStatWidget {
	corewidget := sws.NewCoreWidget(w, h)
	widget := &DemandStatWidget{
		CoreWidget:  *corewidget,
		demandstats: ui.NewTableWithDetails(525+15, 200),
		barchart:    ui.NewStackedBarChartWidget(18, 525, 200),
	}
	widget.barchart.AddCategory("you", COLOR_SALE_YOU)
	widget.barchart.AddCategory("unassigned", COLOR_SALE_UNASSIGNED)
	widget.barchart.AddCategory("competitor", COLOR_SALE_COMPETITOR)
	g.AddGameTimerSubscriber(widget.barchart)
	widget.AddChild(widget.barchart)

	widget.demandstats.Move(0, 220)
	widget.demandstats.AddHeader("Date", 100, ui.TableWithDetailsRowByYearMonthDay)
	widget.demandstats.AddHeader("Price", 100, ui.TableWithDetailsRowByPriceDollar)
	widget.demandstats.AddHeader("Nb servers", 100, ui.TableWithDetailsRowByInteger)
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
	self.barchart.ClearData(t)

	for _, ds := range gamestats.demandsstats {
		self.NewDemandStat(ds)
	}
}

//
// PowerStatWidget is a about global power consumption / generation stat widget
type PowerStatWidget struct {
	sws.CoreWidget
	consumption      *ui.BarChartWidget
	generation       *ui.BarChartWidget
	provided         *ui.BarChartWidget
	labelConsumption *sws.LabelWidget
}

func NewPowerStatWidget(w, h int32, g *Game) *PowerStatWidget {
	corewidget := sws.NewCoreWidget(w, h)
	widget := &PowerStatWidget{
		CoreWidget:       *corewidget,
		consumption:      ui.NewBarChartWidget(525, 150),
		generation:       ui.NewBarChartWidget(525, 150),
		provided:         ui.NewBarChartWidget(525, 150),
		labelConsumption: sws.NewLabelWidget(250, 25, "Current consumption: 0 kwh"),
	}
	g.AddGameTimerSubscriber(widget.consumption)
	widget.consumption.SetChartColor(COLOR_SALE_YOU)
	widget.AddChild(widget.consumption)

	widget.labelConsumption.Move(525, 60)
	widget.AddChild(widget.labelConsumption)

	widget.generation.Move(0, 160)
	widget.generation.SetChartColor(COLOR_SALE_YOU)
	g.AddGameTimerSubscriber(widget.generation)
	widget.AddChild(widget.generation)

	labelGeneration := sws.NewLabelWidget(250, 25, "Generators capacity")
	labelGeneration.Move(525, 220)
	widget.AddChild(labelGeneration)

	widget.provided.Move(0, 320)
	widget.provided.SetChartColor(COLOR_SALE_YOU)
	g.AddGameTimerSubscriber(widget.provided)
	widget.AddChild(widget.provided)

	labelProvided := sws.NewLabelWidget(250, 25, "Utility transmission")
	labelProvided.Move(525, 380)
	widget.AddChild(labelProvided)

	g.GetGameStats().AddPowerStatSubscriber(widget)
	return widget
}

func (self *PowerStatWidget) NewPowerStat(ps *PowerStat) {
	self.consumption.SetPoint(ps.date, int32(ps.consumption))
	self.generation.SetPoint(ps.date, int32(ps.generation))
	self.provided.SetPoint(ps.date, int32(ps.provided))

	self.labelConsumption.SetText(fmt.Sprintf("Current consumption: %d kwh", int(ps.consumption)/1000))
}

func (self *PowerStatWidget) SetGame(t time.Time, gamestats *GameStats) {
	log.Debug("PowerStatWidget::SetGame(", gamestats, ")")
	self.consumption.ClearData(t)
	self.generation.ClearData(t)
	self.provided.ClearData(t)

	for _, ps := range gamestats.powerstats {
		self.NewPowerStat(ps)
	}
}

//
// PlayerStatWidget is a about player info (reputation)
type PlayerStatWidget struct {
	sws.CoreWidget
	reputation      *ui.BarChartWidget
	labelReputation *sws.LabelWidget
}

func (self *PlayerStatWidget) NewReputationStat(rs *ReputationStat) {
	log.Debug("PlayerStatWidget::NewReputationStat(", rs, ")")
	self.reputation.SetPoint(rs.date, int32(rs.reputation*100))
	self.labelReputation.SetText(fmt.Sprintf("Current reputation: %d%%", int(rs.reputation*100)))
}

func (self *PlayerStatWidget) SetGame(gamestats *GameStats) {
	gamestats.AddReputationStatSubscriber(self)

	for _, rs := range gamestats.reputationstats {
		self.NewReputationStat(rs)
	}
}

func NewPlayerStatWidget(w, h int32, g *Game) *PlayerStatWidget {
	corewidget := sws.NewCoreWidget(w, h)
	widget := &PlayerStatWidget{
		CoreWidget:      *corewidget,
		reputation:      ui.NewBarChartWidget(525, 150),
		labelReputation: sws.NewLabelWidget(200, 25, "Current reputation: 0%"),
	}

	g.AddGameTimerSubscriber(widget.reputation)
	widget.reputation.SetChartColor(COLOR_SALE_YOU)
	widget.AddChild(widget.reputation)

	widget.labelReputation.Move(525, 60)
	widget.AddChild(widget.labelReputation)

	return widget
}
