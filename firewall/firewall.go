package firewall

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"github.com/nzin/dctycoon/supplier"

	"github.com/BixData/gluabit32"
	log "github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
)

//Firewall:
// by default if a filter is broken, everything PASS

// Ideas:
//
//- several "level" (from 1 a 15?)
//- one emiter / one firewall / one collector
//
//Basic
//- Ddos icmp: smurf attack https://blog.cloudflare.com/deep-inside-a-dns-amplification-ddos-attack/) -> solution filter source and destination IP (to avoid internal forged IP)
//- Ddos udp (DNS amplification?) -> to an IP. It comes from some IP (open resolver) to a specifc IP in a time frame  with volume -> block trigger by volume? or by IP (or both) ...
//
//- ssh attack: login/password, wp-admin attack: "admin"/password
//

// FirewallSubscriber is used to be inform when a packet has been filtered by the firewall
// (mainly used by MainFirewallWidget)
type FirewallSubscriber interface {
	// when a rack change from status (RACK_NORMAL_STATE, RACK_OVER_CURRENT, RACK_HEAT_WARNING, RACK_OVER_HEAT, RACK_MELTING)
	PacketFiltered(event *FirewallEvent)
}

type FirewallEvent struct {
	Time   time.Time
	Packet *Packet
	Pass   bool
}

type Firewall struct {
	vm                *lua.LState
	dcClassBNetwork   string
	rules             string
	generator         *PacketGenerator
	lastEvents        []*FirewallEvent
	packetSubscribers []FirewallSubscriber
}

var emptyRules = `local bit32 = require 'bit32'
-- all your datacenter servers have IPs like %s.x.x
datacenterClassB="%s"

-- ipsrc,ipdst are string like '192.168.18.100'
-- header is a [8]bytes string
-- payload is a string
function filterIcmp( ipsrc, ipdst, header, payload)
	return true;
end

-- ipsrc,ipdst are string like '192.168.18.100'
-- srcPort, dstPort are number (from 0 to 65535)
-- payload is a string
function filterUdp( ipsrc, ipdst, srcPort, dstPort, payload)
	return true;
end

-- ipsrc,ipdst are string like '192.168.18.100'
-- srcPort, dstPort are number (from 0 to 65535)
-- flags is a byte (hint: if you want to filter SYN packet do something like 'if (bit32.band(flags,0x02)==2) then ...')
-- payload is a string
function filterTcp( ipsrc, ipdst, srcPort, dstPort, flags, payload)
	return true;
end
`

func (firewall *Firewall) AddFirewallSubscriber(subscriber FirewallSubscriber) {
	for _, s := range firewall.packetSubscribers {
		if s == subscriber {
			return
		}
	}
	firewall.packetSubscribers = append(firewall.packetSubscribers, subscriber)
}

func (firewall *Firewall) RemoveFirewallSubscriber(subscriber FirewallSubscriber) {
	for i, s := range firewall.packetSubscribers {
		if s == subscriber {
			firewall.packetSubscribers = append(firewall.packetSubscribers[:i], firewall.packetSubscribers[i+1:]...)
			break
		}
	}
}

func (firewall *Firewall) GetLastEvents() []*FirewallEvent {
	return firewall.lastEvents
}

func (firewall *Firewall) GetRules() string {
	return firewall.rules
}

// GetDatacenterClassBNetwork returns the first 2 numbers of the A.B.C.D/16 of the Datacenter ipv4 class B network
func (firewall *Firewall) GetDatacenterClassBNetwork() string {
	return firewall.dcClassBNetwork
}

func (firewall *Firewall) ResetRules() {
	firewall.SetRulesAndApply(fmt.Sprintf(emptyRules, firewall.dcClassBNetwork, firewall.dcClassBNetwork))
}

// SetRulesAndApply will try to load the rules into the lua interpreter
// return error if the rules cannot be applied
func (firewall *Firewall) SetRulesAndApply(rules string) error {
	firewall.rules = rules
	if firewall.vm != nil {
		firewall.vm.Close()
	}
	firewall.vm = lua.NewState()
	gluabit32.Preload(firewall.vm)
	return firewall.vm.DoString(firewall.rules)
}

// SubmitIcmp submit an ICMP packet and returns true if it passes through the firewall
func (firewall *Firewall) SubmitIcmp(ipsrc, ipdst string, header [8]byte, payload string) bool {
	if err := firewall.vm.CallByParam(lua.P{
		Fn:      firewall.vm.GetGlobal("filterIcmp"), // name of Lua function
		NRet:    1,                                   // number of returned values
		Protect: true,                                // return err or panic
	}, lua.LString(ipsrc), lua.LString(ipdst), lua.LString(string(header[:])), lua.LString(payload)); err != nil {
		fmt.Println(err)
		return true
	}
	if ret, ok := firewall.vm.Get(-1).(lua.LBool); ok {
		return bool(ret)
	}
	return true
}

// SubmitUdp submit an UDP packet and returns true if it passes through the firewall
func (firewall *Firewall) SubmitUdp(ipsrc, ipdst string, srcPort, dstPort uint16, payload string) bool {
	if err := firewall.vm.CallByParam(lua.P{
		Fn:      firewall.vm.GetGlobal("filterUdp"), // name of Lua function
		NRet:    1,                                  // number of returned values
		Protect: true,                               // return err or panic
	}, lua.LString(ipsrc), lua.LString(ipdst), lua.LNumber(srcPort), lua.LNumber(dstPort), lua.LString(payload)); err != nil {
		fmt.Println(err)
		return true
	}
	if ret, ok := firewall.vm.Get(-1).(lua.LBool); ok {
		return bool(ret)
	}
	return true
}

// SubmitTcp submit a TCP packet and returns true if it passes through the firewall
func (firewall *Firewall) SubmitTcp(ipsrc, ipdst string, srcPort, dstPort uint16, flags uint8, payload string) bool {
	if err := firewall.vm.CallByParam(lua.P{
		Fn:      firewall.vm.GetGlobal("filterTcp"), // name of Lua function
		NRet:    1,                                  // number of returned values
		Protect: true,                               // return err or panic
	}, lua.LString(ipsrc), lua.LString(ipdst), lua.LNumber(srcPort), lua.LNumber(dstPort), lua.LNumber(flags), lua.LString(payload)); err != nil {
		fmt.Println(err)
		return true
	}
	if ret, ok := firewall.vm.Get(-1).(lua.LBool); ok {
		return bool(ret)
	}
	return true
}

func (firewall *Firewall) Load(data map[string]interface{}) {
	firewall.dcClassBNetwork = (data["datacenterClassB"].(string))
	decoded, err := base64.StdEncoding.DecodeString(data["rules"].(string))
	if err != nil {
		firewall.ResetRules()
	}
	err = firewall.SetRulesAndApply(string(decoded))
	if err != nil {
		firewall.ResetRules()
	}
	firewall.generator.SetGame(firewall.dcClassBNetwork)

	// load past events
	firewall.lastEvents = make([]*FirewallEvent, 0, 0)

	array := data["lastevents"].([]interface{})
	for i := 0; i < len(array); i++ {
		event := array[i].(map[string]interface{})
		var year, month, day int
		fmt.Sscanf(event["date"].(string), "%d-%d-%d", &year, &month, &day)

		firewallevent := &FirewallEvent{
			Time:   time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC),
			Packet: NewPacket(event["packet"].(map[string]interface{})),
			Pass:   event["pass"].(bool),
		}
		firewall.lastEvents = append(firewall.lastEvents, firewallevent)
	}
}

func (firewall *Firewall) Save() string {
	str := fmt.Sprintf(`{"datacenterClassB":"%s","rules":"%s", "lastevents":[`, firewall.dcClassBNetwork, base64.StdEncoding.EncodeToString([]byte(firewall.rules)))
	for i := 0; i < len(firewall.lastEvents); i++ {
		event := firewall.lastEvents[i]
		if i > 0 {
			str += ","
		}
		pass := "true"
		if event.Pass == false {
			pass = "false"
		}
		str += fmt.Sprintf(`{"date": "%d-%d-%d","packet":%s, "pass":%s}`, event.Time.Year(), event.Time.Month(), event.Time.Day(), event.Packet.Save(), pass)
	}

	return str + "]}"
}

func (firewall *Firewall) GenerateTraffic(reputation *supplier.Reputation, time time.Time) {
	log.Debug("Firewall::GenerateTraffic(", reputation, ",", time, ")")
	packet := firewall.generator.GenerateRandomPacket()
	if packet == nil {
		return
	}

	var res bool
	switch packet.PacketType {
	case PACKET_ICMP:
		res = firewall.SubmitIcmp(packet.Ipsrc, packet.Ipdst, packet.IcmpHeader, packet.Payload)
	case PACKET_UDP:
		res = firewall.SubmitUdp(packet.Ipsrc, packet.Ipdst, packet.SrcPort, packet.DstPort, packet.Payload)
	case PACKET_TCP:
		res = firewall.SubmitTcp(packet.Ipsrc, packet.Ipdst, packet.SrcPort, packet.DstPort, packet.Tcpflags, packet.Payload)
	}

	// store in the last events
	if len(firewall.lastEvents) == 20 {
		firewall.lastEvents = firewall.lastEvents[1:]
	}
	event := &FirewallEvent{
		Time:   time,
		Packet: packet,
		Pass:   res,
	}
	firewall.lastEvents = append(firewall.lastEvents, event)
	for _, s := range firewall.packetSubscribers {
		s.PacketFiltered(event)
	}

	//	if time.Day()%5 == 1 {
	if res != packet.Harmless {
		// something went wrong
		reputation.RecordNegativePoint(time)
	}
	reputation.RecordPositivePoint(time)
	//	}
}

func NewFirewall() *Firewall {
	dcclassb := fmt.Sprintf("%d.%d", 20+rand.Int()%100, rand.Int()%254)
	firewall := &Firewall{
		vm:              lua.NewState(),
		dcClassBNetwork: dcclassb,
		generator:       NewPacketGenerator(dcclassb),
		lastEvents:      make([]*FirewallEvent, 0, 0),
	}

	firewall.SetRulesAndApply(fmt.Sprintf(emptyRules, firewall.dcClassBNetwork, firewall.dcClassBNetwork))

	return firewall
}
