package firewall

import (
	"encoding/base64"
	"fmt"

	"github.com/BixData/gluabit32"
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

type Firewall struct {
	vm    *lua.LState
	rules string
}

var emptyRules = `local bit32 = require 'bit32'

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

func (firewall *Firewall) GetRules() string {
	return firewall.rules
}

func (firewall *Firewall) ResetRules() {
	firewall.SetRulesAndApply(emptyRules)
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
	decoded, err := base64.StdEncoding.DecodeString(data["rules"].(string))
	if err != nil {
		firewall.ResetRules()
	}
	err = firewall.SetRulesAndApply(string(decoded))
	if err != nil {
		firewall.ResetRules()
	}
}

func (firewall *Firewall) Save() string {
	return fmt.Sprintf(`{"rules":"%s"}`, base64.StdEncoding.EncodeToString([]byte(firewall.rules)))
}

func NewFirewall() *Firewall {
	firewall := &Firewall{
		vm: lua.NewState(),
	}
	firewall.SetRulesAndApply(emptyRules)

	return firewall
}
