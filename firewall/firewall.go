package firewall

import (
	"fmt"

	"github.com/BixData/gluabit32"
	lua "github.com/yuin/gopher-lua"
)

//Firewall:
//- several "level" (from 1 a 15?)
//- one emiter / one firewall / one collector
//
//Basic
//- Ddos icmp: smurf attack https://blog.cloudflare.com/deep-inside-a-dns-amplification-ddos-attack/) -> solution filter source and destination IP (to avoid internal forged IP)
//- Ddos udp (DNS amplification?) -> to an IP. It comes from some IP (open resolver) to a specifc IP in a time frame  with volume -> block trigger by volume? or by IP (or both) ...
//
//- ssh attack: login/password, wp-admin attack: "admin"/password
//
//
//- cut packet to skip filters?  :-) (ok a bit too far)

type Firewall struct {
	vm *lua.LState
}

func (firewall *Firewall) LoadRules(rules string) error {
	if firewall.vm != nil {
		firewall.vm.Close()
	}
	firewall.vm = lua.NewState()
	gluabit32.Preload(firewall.vm)
	return firewall.vm.DoString(rules)
}

// SubmitIcmp submit an ICMP packet and returns true if it passes through the firewall
func (firewall *Firewall) SubmitIcmp(ipsrc, ipdst string, header [8]byte, payload string) bool {
	if err := firewall.vm.CallByParam(lua.P{
		Fn:      firewall.vm.GetGlobal("filterIcmp"), // name of Lua function
		NRet:    1,                                   // number of returned values
		Protect: true,                                // return err or panic
	}, lua.LString(ipsrc), lua.LString(ipdst), lua.LString(string(header[:])), lua.LString(payload)); err != nil {
		fmt.Println(err)
		return false
	}
	if ret, ok := firewall.vm.Get(-1).(lua.LBool); ok {
		return bool(ret)
	}
	return false
}

// SubmitUdp submit an UDP packet and returns true if it passes through the firewall
func (firewall *Firewall) SubmitUdp(ipsrc, ipdst string, srcPort, dstPort uint16, payload string) bool {
	if err := firewall.vm.CallByParam(lua.P{
		Fn:      firewall.vm.GetGlobal("filterUdp"), // name of Lua function
		NRet:    1,                                  // number of returned values
		Protect: true,                               // return err or panic
	}, lua.LString(ipsrc), lua.LString(ipdst), lua.LNumber(srcPort), lua.LNumber(dstPort), lua.LString(payload)); err != nil {
		fmt.Println(err)
		return false
	}
	if ret, ok := firewall.vm.Get(-1).(lua.LBool); ok {
		return bool(ret)
	}
	return false
}

// SubmitTcp submit a TCP packet and returns true if it passes through the firewall
func (firewall *Firewall) SubmitTcp(ipsrc, ipdst string, srcPort, dstPort uint16, flags uint8, payload string) bool {
	if err := firewall.vm.CallByParam(lua.P{
		Fn:      firewall.vm.GetGlobal("filterTcp"), // name of Lua function
		NRet:    1,                                  // number of returned values
		Protect: true,                               // return err or panic
	}, lua.LString(ipsrc), lua.LString(ipdst), lua.LNumber(srcPort), lua.LNumber(dstPort), lua.LNumber(flags), lua.LString(payload)); err != nil {
		fmt.Println(err)
		return false
	}
	if ret, ok := firewall.vm.Get(-1).(lua.LBool); ok {
		return bool(ret)
	}
	return false
}

func NewFirewall() *Firewall {
	firewall := &Firewall{
		vm: lua.NewState(),
	}
	return firewall
}
