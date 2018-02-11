package dctycoon

import (
	"fmt"

	"github.com/nzin/dctycoon/firewall"
	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/ui"
	"github.com/nzin/sws"
)

type MainFirewallWidget struct {
	rootwindow     *sws.RootWidget
	mainwidget     *sws.MainWidget
	tabwidget      *sws.TabWidget
	firewallrules  *FirewallRulesWidget
	firewallstatus *FirewallStatusWidget
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
	self.firewallstatus.SetGame(firewall)
}

// NewMainInventoryWidget presents the pool and offer management window
func NewMainFirewallWidget(root *sws.RootWidget) *MainFirewallWidget {
	mainwidget := sws.NewMainWidget(900, 600, " Firewall Management ", true, true)
	mainwidget.Center(root)

	widget := &MainFirewallWidget{
		rootwindow:     root,
		mainwidget:     mainwidget,
		tabwidget:      sws.NewTabWidget(200, 200),
		firewallrules:  NewFirewallRulesWidget(200, 200, root),
		firewallstatus: NewFirewallStatusWidget(200, 200),
	}

	widget.tabwidget.AddTab("firewall status", widget.firewallstatus)
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

type FirewallEventDetailIcmp struct {
	sws.CoreWidget
}

func NewFirewallEventDetailIcmp(bgColor uint32, packet *firewall.Packet) *FirewallEventDetailIcmp {
	corewidget := sws.NewCoreWidget(900, 75)
	corewidget.SetColor(bgColor)
	details := &FirewallEventDetailIcmp{
		CoreWidget: *corewidget,
	}
	labeltype := sws.NewLabelWidget(100, 25, "Type:")
	labeltype.SetColor(bgColor)
	labeltype.Move(0, 0)
	details.AddChild(labeltype)

	icmptype := sws.NewLabelWidget(100, 25, fmt.Sprintf("%d", packet.IcmpHeader[0]))
	icmptype.SetColor(bgColor)
	icmptype.Move(100, 0)
	details.AddChild(icmptype)

	labelcode := sws.NewLabelWidget(100, 25, "Code:")
	labelcode.SetColor(bgColor)
	labelcode.Move(0, 25)
	details.AddChild(labelcode)

	icmpcode := sws.NewLabelWidget(100, 25, fmt.Sprintf("%d", packet.IcmpHeader[1]))
	icmpcode.SetColor(bgColor)
	icmpcode.Move(100, 25)
	details.AddChild(icmpcode)

	labelpayload := sws.NewLabelWidget(100, 25, "Payload size:")
	labelpayload.SetColor(bgColor)
	labelpayload.Move(0, 50)
	details.AddChild(labelpayload)

	max100 := len(packet.Payload)
	if max100 > 100 {
		max100 = 100
	}
	icmppayload := sws.NewLabelWidget(800, 25, fmt.Sprintf("%d (content='%s')", len(packet.Payload), packet.Payload[:max100]))
	icmppayload.SetColor(bgColor)
	icmppayload.Move(100, 50)
	details.AddChild(icmppayload)

	return details
}

type FirewallEventDetailUdp struct {
	sws.CoreWidget
}

func NewFirewallEventDetailUdp(bgColor uint32, packet *firewall.Packet) *FirewallEventDetailUdp {
	corewidget := sws.NewCoreWidget(900, 75)
	corewidget.SetColor(bgColor)
	details := &FirewallEventDetailUdp{
		CoreWidget: *corewidget,
	}
	labelsrcport := sws.NewLabelWidget(100, 25, "Src port:")
	labelsrcport.SetColor(bgColor)
	labelsrcport.Move(0, 0)
	details.AddChild(labelsrcport)

	srcport := sws.NewLabelWidget(100, 25, fmt.Sprintf("%d", packet.SrcPort))
	srcport.SetColor(bgColor)
	srcport.Move(100, 0)
	details.AddChild(srcport)

	labeldstport := sws.NewLabelWidget(100, 25, "Dst port:")
	labeldstport.SetColor(bgColor)
	labeldstport.Move(0, 25)
	details.AddChild(labeldstport)

	dstport := sws.NewLabelWidget(100, 25, fmt.Sprintf("%d", packet.DstPort))
	dstport.SetColor(bgColor)
	dstport.Move(100, 25)
	details.AddChild(dstport)

	labelpayload := sws.NewLabelWidget(100, 25, "Payload:")
	labelpayload.SetColor(bgColor)
	labelpayload.Move(0, 50)
	details.AddChild(labelpayload)

	payload := sws.NewLabelWidget(800, 25, packet.Payload)
	payload.SetColor(bgColor)
	payload.Move(100, 50)
	details.AddChild(payload)

	return details
}

type FirewallEventDetailTcp struct {
	sws.CoreWidget
}

func NewFirewallEventDetailTcp(bgColor uint32, packet *firewall.Packet) *FirewallEventDetailTcp {
	corewidget := sws.NewCoreWidget(900, 100)
	corewidget.SetColor(bgColor)
	details := &FirewallEventDetailTcp{
		CoreWidget: *corewidget,
	}
	labelsrcport := sws.NewLabelWidget(100, 25, "Src port:")
	labelsrcport.SetColor(bgColor)
	labelsrcport.Move(0, 0)
	details.AddChild(labelsrcport)

	srcport := sws.NewLabelWidget(100, 25, fmt.Sprintf("%d", packet.SrcPort))
	srcport.SetColor(bgColor)
	srcport.Move(100, 0)
	details.AddChild(srcport)

	labeldstport := sws.NewLabelWidget(100, 25, "Dst port:")
	labeldstport.SetColor(bgColor)
	labeldstport.Move(0, 25)
	details.AddChild(labeldstport)

	dstport := sws.NewLabelWidget(100, 25, fmt.Sprintf("%d", packet.DstPort))
	dstport.SetColor(bgColor)
	dstport.Move(100, 25)
	details.AddChild(dstport)

	labelflags := sws.NewLabelWidget(300, 25, "Flags:")
	labelflags.SetColor(bgColor)
	labelflags.Move(0, 50)
	details.AddChild(labelflags)

	flagstxt := fmt.Sprintf("0x%x ", packet.Tcpflags)
	if packet.Tcpflags&0x1 != 0 {
		flagstxt += "[FIN] "
	}
	if packet.Tcpflags&0x2 != 0 {
		flagstxt += "[SYN] "
	}
	if packet.Tcpflags&0x4 != 0 {
		flagstxt += "[RST] "
	}
	if packet.Tcpflags&0x8 != 0 {
		flagstxt += "[PSH] "
	}
	if packet.Tcpflags&0x10 != 0 {
		flagstxt += "[ACK] "
	}
	if packet.Tcpflags&0x20 != 0 {
		flagstxt += "[URG] "
	}
	if packet.Tcpflags&0x40 != 0 {
		flagstxt += "[ECE] "
	}
	if packet.Tcpflags&0x80 != 0 {
		flagstxt += "[CWR] "
	}
	flags := sws.NewLabelWidget(500, 25, flagstxt)
	flags.SetColor(bgColor)
	flags.Move(100, 50)
	details.AddChild(flags)

	labelpayload := sws.NewLabelWidget(100, 25, "Payload:")
	labelpayload.SetColor(bgColor)
	labelpayload.Move(0, 75)
	details.AddChild(labelpayload)

	payload := sws.NewLabelWidget(800, 25, packet.Payload)
	payload.SetColor(bgColor)
	payload.Move(100, 75)
	details.AddChild(payload)

	return details
}

type FirewallStatusWidget struct {
	sws.CoreWidget
	firewall   *firewall.Firewall
	pastEvents *ui.TableWithDetails
	sucessrate *sws.LabelWidget
}

func (self *FirewallStatusWidget) addFirewallEvent(event *firewall.FirewallEvent) {
	bgColor := uint32(0xffdddddd)
	if event.Pass != event.Packet.Harmless {
		bgColor = uint32(0xffffaaaa)
	}

	var details sws.Widget
	labels := make([]string, 0, 0)
	labels = append(labels, fmt.Sprintf("%d-%d-%d", event.Time.Day(), event.Time.Month(), event.Time.Year()))
	switch event.Packet.PacketType {
	case firewall.PACKET_ICMP:
		details = NewFirewallEventDetailIcmp(bgColor, event.Packet)
		labels = append(labels, "ICMP")
	case firewall.PACKET_UDP:
		details = NewFirewallEventDetailUdp(bgColor, event.Packet)
		labels = append(labels, "UDP")
	case firewall.PACKET_TCP:
		details = NewFirewallEventDetailTcp(bgColor, event.Packet)
		labels = append(labels, "TCP")
	}
	labels = append(labels, event.Packet.Ipsrc)
	labels = append(labels, event.Packet.Ipdst)

	row := ui.NewTableWithDetailsRow(bgColor, labels, details)
	self.pastEvents.AddRowTop(row)

	events := self.firewall.GetLastEvents()
	var success, fail int
	for _, e := range events {
		if e.Packet.Harmless == e.Pass {
			success++
		} else {
			fail++
		}
	}
	if success+fail > 0 {
		self.sucessrate.SetText(fmt.Sprintf("%d%%", 100*success/(success+fail)))
	}
}

func (self *FirewallStatusWidget) SetGame(firewall *firewall.Firewall) {
	self.firewall = firewall
	firewall.AddFirewallSubscriber(self)
	self.pastEvents.ClearRows()
	for _, e := range firewall.GetLastEvents() {
		self.addFirewallEvent(e)
	}
}

func (self *FirewallStatusWidget) PacketFiltered(event *firewall.FirewallEvent) {
	self.addFirewallEvent(event)
	self.pastEvents.Truncate(20)
}

func (self *FirewallStatusWidget) Resize(w, h int32) {
	self.CoreWidget.Resize(w, h)
	if w < 40 {
		w = 40
	}
	if h < 120 {
		h = 120
	}
	self.pastEvents.Resize(800, h-90)
}

func NewFirewallStatusWidget(w, h int32) *FirewallStatusWidget {
	corewidget := sws.NewCoreWidget(w, h)
	widget := &FirewallStatusWidget{
		CoreWidget: *corewidget,
		pastEvents: ui.NewTableWithDetails(525, 200),
		sucessrate: sws.NewLabelWidget(100, 25, "100%"),
	}
	sucesslabel := sws.NewLabelWidget(200, 25, "Firewall effectiveness:")
	sucesslabel.Move(10, 30)
	widget.AddChild(sucesslabel)

	widget.sucessrate.Move(210, 30)
	widget.AddChild(widget.sucessrate)

	detailslablel := sws.NewLabelWidget(200, 25, "Last 20 packets filtered:")
	detailslablel.Move(10, 60)
	widget.AddChild(detailslablel)

	widget.pastEvents.Move(10, 90)
	widget.pastEvents.AddHeader("Date", 100, ui.TableWithDetailsRowByYearMonthDay)
	widget.pastEvents.AddHeader("Protocol", 100, ui.TableWithDetailsRowByString)
	widget.pastEvents.AddHeader("Source IP", 200, nil)
	widget.pastEvents.AddHeader("Dest IP", 200, nil)
	widget.AddChild(widget.pastEvents)

	return widget
}
