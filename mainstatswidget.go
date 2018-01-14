package dctycoon

import (
	"fmt"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/sws"

	log "github.com/sirupsen/logrus"
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
	mainwidget := sws.NewMainWidget(850, 400, " Graph and Statistics ", true, true)
	mainwidget.Center(root)

	widget := &MainStatsWidget{
		rootwindow:      root,
		mainwidget:      mainwidget,
		tabwidget:       sws.NewTabWidget(200, 200),
		game:            g,
		opponentswidget: NewOpponentStatWidget(200, 200),
		demandswidget:   NewDemandStatWidget(200, 200, g.GetGameStats()),
	}

	widget.mainwidget.SetInnerWidget(widget.tabwidget)

	widget.mainwidget.SetCloseCallback(func() {
		widget.Hide()
	})

	widget.tabwidget.AddTab("Opponents", widget.opponentswidget)
	widget.tabwidget.AddTab("Customer demands", widget.demandswidget)

	return widget
}

func (self *MainStatsWidget) SetGame() {
	self.opponentswidget.SetGame(self.game)
	self.demandswidget.SetGame(self.game.GetGameStats())
}

func (self *MainStatsWidget) LoadGame() {
	self.opponentswidget.SetGame(self.game)
	self.demandswidget.SetGame(self.game.GetGameStats())
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

type DemandStatWidgetLine struct {
	sws.CoreWidget
	details     *sws.FlatButtonWidget
	date        *sws.LabelWidget
	price       *sws.LabelWidget
	nbservers   *sws.LabelWidget
	buyer       *sws.LabelWidget
	showdetails bool
	yoffset     int32
}

func NewDemandStatWidgetLine(container *sws.VBoxWidget, stat *DemandStat) *DemandStatWidgetLine {
	corewidget := sws.NewCoreWidget(525, 25)
	nbservers := int32(0)
	for _, s := range stat.serverdemands {
		nbservers += s.nb
	}
	price := "-"
	bgcolor := uint32(0xffaaaaaa)
	if stat.buyer != "" {
		price = fmt.Sprintf("%.2f $", stat.price)
		bgcolor = 0xffdddddd
	}

	line := &DemandStatWidgetLine{
		CoreWidget:  *corewidget,
		details:     sws.NewFlatButtonWidget(25, 25, ""),
		date:        sws.NewLabelWidget(100, 25, fmt.Sprintf("%d-%d-%d", stat.date.Day(), stat.date.Month(), stat.date.Year())),
		price:       sws.NewLabelWidget(100, 25, price),
		nbservers:   sws.NewLabelWidget(100, 25, fmt.Sprintf("%d", nbservers)),
		buyer:       sws.NewLabelWidget(200, 25, stat.buyer),
		showdetails: false,
		yoffset:     25,
	}
	line.SetColor(bgcolor)
	line.details.SetColor(bgcolor)
	if surface, err := global.LoadImageAsset("assets/ui/icon-arrowhead-pointing-to-the-right.png"); err == nil {
		line.details.SetImageSurface(surface)
	}
	line.details.SetClicked(func() {
		if line.showdetails == false {
			if surface, err := global.LoadImageAsset("assets/ui/icon-sort-down.png"); err == nil {
				line.details.SetImageSurface(surface)
				line.Resize(525, line.yoffset)
				container.Rebox()
			}
		} else {
			if surface, err := global.LoadImageAsset("assets/ui/icon-arrowhead-pointing-to-the-right.png"); err == nil {
				line.details.SetImageSurface(surface)
				line.Resize(525, 25)
				container.Rebox()
			}
		}
		line.showdetails = !line.showdetails
	})
	line.AddChild(line.details)

	line.date.Move(25, 0)
	line.date.SetColor(bgcolor)
	line.AddChild(line.date)

	line.price.Move(125, 0)
	line.price.SetColor(bgcolor)
	line.AddChild(line.price)

	line.nbservers.Move(225, 0)
	line.nbservers.SetColor(bgcolor)
	line.AddChild(line.nbservers)

	line.buyer.Move(325, 0)
	line.buyer.SetColor(bgcolor)
	line.AddChild(line.buyer)

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
// DemandStatWidget is an customer demand stat view
type DemandStatWidget struct {
	sws.CoreWidget
	vbox      *sws.VBoxWidget
	scroll    *sws.ScrollWidget
	date      *sws.LabelWidget
	price     *sws.LabelWidget
	nbservers *sws.LabelWidget
	buyer     *sws.LabelWidget
}

func NewDemandStatWidget(w, h int32, gamestats *GameStats) *DemandStatWidget {
	corewidget := sws.NewCoreWidget(w, h)
	widget := &DemandStatWidget{
		CoreWidget: *corewidget,
		vbox:       sws.NewVBoxWidget(525, 0),
		scroll:     sws.NewScrollWidget(w, h-25),
		date:       sws.NewLabelWidget(100, 25, "Date"),
		price:      sws.NewLabelWidget(100, 25, "$/month"),
		nbservers:  sws.NewLabelWidget(100, 25, "Nb servers"),
		buyer:      sws.NewLabelWidget(100, 25, "Buyer"),
	}

	widget.date.Move(25, 0)
	widget.AddChild(widget.date)

	widget.price.Move(125, 0)
	widget.AddChild(widget.price)

	widget.nbservers.Move(225, 0)
	widget.AddChild(widget.nbservers)

	widget.buyer.Move(325, 0)
	widget.AddChild(widget.buyer)

	widget.scroll.Move(0, 25)
	widget.scroll.SetInnerWidget(widget.vbox)
	widget.scroll.ShowHorizontalScrollbar(false)
	widget.AddChild(widget.scroll)

	gamestats.AddDemandStatSubscriber(widget)
	return widget
}

// NewDemandStat is part of DemandStatSubscriber interface
func (self *DemandStatWidget) NewDemandStat(ds *DemandStat) {
	log.Debug("DemandStatWidget::NewDemandStat(", ds, ")")
	line := NewDemandStatWidgetLine(self.vbox, ds)
	self.vbox.AddChildTop(line)
	self.scroll.Resize(self.Width(), self.Height()-25)
}

func (self *DemandStatWidget) Resize(w, h int32) {
	self.CoreWidget.Resize(w, h)
	self.scroll.Resize(w, h-25)
}

func (self *DemandStatWidget) SetGame(gamestats *GameStats) {
	log.Debug("DemandStatWidget::SetGame(", gamestats, ")")
	self.vbox.RemoveAllChildren()

	for _, ds := range gamestats.demandsstats {
		line := NewDemandStatWidgetLine(self.vbox, ds)
		self.vbox.AddChildTop(line)
	}
	self.scroll.Resize(self.Width(), self.Height()-25)
}
