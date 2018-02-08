package dctycoon

import (
	"fmt"

	"github.com/nzin/dctycoon/firewall"
	"github.com/nzin/dctycoon/global"
	"github.com/nzin/sws"
)

type MainFirewallWidget struct {
	rootwindow    *sws.RootWidget
	mainwidget    *sws.MainWidget
	tabwidget     *sws.TabWidget
	firewallrules *FirewallRulesWidget
}

func (self *MainFirewallWidget) Show() {
	self.rootwindow.AddChild(self.mainwidget)
	self.rootwindow.SetFocus(self.mainwidget)
	self.tabwidget.SelectTab(0)
}

func (self *MainFirewallWidget) Hide() {
	self.rootwindow.RemoveChild(self.mainwidget)
	children := self.rootwindow.GetChildren()
	if len(children) > 0 {
		self.rootwindow.SetFocus(children[0])
	}
}

func (self *MainFirewallWidget) SetGame(firewall *firewall.Firewall) {
	self.firewallrules.SetGame(firewall)
}

// NewMainInventoryWidget presents the pool and offer management window
func NewMainFirewallWidget(root *sws.RootWidget) *MainFirewallWidget {
	mainwidget := sws.NewMainWidget(900, 600, " Firewall Management ", true, true)
	mainwidget.Center(root)

	widget := &MainFirewallWidget{
		rootwindow:    root,
		mainwidget:    mainwidget,
		tabwidget:     sws.NewTabWidget(200, 200),
		firewallrules: NewFirewallRulesWidget(200, 200, root),
	}

	widget.tabwidget.AddTab("firewall rules", widget.firewallrules)

	widget.mainwidget.SetInnerWidget(widget.tabwidget)

	widget.mainwidget.SetCloseCallback(func() {
		widget.Hide()
	})

	return widget
}

type FirewallRulesWidget struct {
	sws.CoreWidget
	firewall *firewall.Firewall
	rules    *sws.TextAreaWidget
	apply    *sws.ButtonWidget
	reset    *sws.ButtonWidget
}

func (self *FirewallRulesWidget) SetGame(firewall *firewall.Firewall) {
	self.firewall = firewall
	self.rules.SetText(firewall.GetRules())
}

func (self *FirewallRulesWidget) Resize(w, h int32) {
	self.CoreWidget.Resize(w, h)
	if w < 40 {
		w = 40
	}
	if h < 65 {
		h = 55
	}
	self.rules.Resize(w-20, h-45)
	self.apply.Move(10, h-35)
	self.reset.Move(170, h-35)
}

func NewFirewallRulesWidget(w, h int32, root *sws.RootWidget) *FirewallRulesWidget {
	corewidget := sws.NewCoreWidget(w, h)
	widget := &FirewallRulesWidget{
		CoreWidget: *corewidget,
		firewall:   nil,
		rules:      sws.NewTextAreaWidget(w-20, h-45, ""),
		apply:      sws.NewButtonWidget(150, 25, "Apply"),
		reset:      sws.NewButtonWidget(150, 25, "Reset"),
	}

	widget.rules.Move(10, 10)
	widget.AddChild(widget.rules)

	widget.apply.Move(10, h-35)
	widget.apply.SetClicked(func() {
		rules := widget.firewall.GetRules()
		err := widget.firewall.SetRulesAndApply(widget.rules.GetText())
		if err != nil {
			widget.firewall.SetRulesAndApply(rules)
			iconsurface, _ := global.LoadImageAsset("assets/ui/icon-triangular-big.png")
			sws.ShowModalErrorSurfaceicon(root, "Syntax error", iconsurface, fmt.Sprintf("Error applying the rules: %v", err), nil)
		}
	})
	widget.AddChild(widget.apply)

	widget.reset.Move(170, h-35)
	widget.reset.SetClicked(func() {
		iconsurface, _ := global.LoadImageAsset("assets/ui/icon-triangular-big.png")
		sws.ShowModalYesNoSurfaceicon(root, "Reset all firewall rules", iconsurface, "Are you sure you want to reset all the rules?", func() {
			widget.firewall.ResetRules()
			widget.rules.SetText(widget.firewall.GetRules())
		}, nil)
	})
	widget.AddChild(widget.reset)

	return widget
}
