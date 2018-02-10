package firewall

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"encoding/json"

	"github.com/nzin/dctycoon/global"
)

const (
	PACKET_ICMP = iota
	PACKET_UDP  = iota
	PACKET_TCP  = iota
)

type Packet struct {
	PacketType int
	Ipsrc      string
	Ipdst      string
	SrcPort    uint16
	DstPort    uint16
	IcmpHeader [8]byte
	Payload    string
	Tcpflags   uint8
	Harmless   bool
}

type JsonPacket struct {
	PacketType string
	Ipsrc      string
	Ipdst      string
	SrcPort    string
	DstPort    string
	IcmpHeader string
	Payload    string
	Tcpflags   string
	Harmless   bool
}

type PacketGenerator struct {
	dcclassb string
}

func (generator *PacketGenerator) generateIP(instruction string) string {
	if instruction == "[IP_IN]" {
		return fmt.Sprintf("%s.%d.%d", generator.dcclassb, rand.Int()%253, rand.Int()%253)
	}
	if instruction == "[IP_OUT]" {
		byte1 := rand.Int() % 253
		byte2 := rand.Int() % 253
		if byte1 == 127 || byte1 == 10 || byte1 == 172 || (fmt.Sprintf("%d.%d", byte1, byte2) == generator.dcclassb) {
			byte1++
		}
		return fmt.Sprintf("%d.%d.%d.%d", byte1, byte2, rand.Int()%253, rand.Int()%253)
	}
	return instruction
}

func (generator *PacketGenerator) generatePort(instruction string) uint16 {
	if instruction == "[RANDOM]" {
		return uint16(rand.Int()%65534 + 1)
	}

	if ret, err := strconv.Atoi(instruction); err != nil {
		return 30000
	} else {
		return uint16(ret)
	}
}

func (generator *PacketGenerator) generatePayload(instruction string) string {
	if instruction == "[65K]" {
		return strings.Repeat("A", 65536)
	}
	return instruction
}

func (generator *PacketGenerator) generateIcmpHeader(instruction string) [8]byte {
	var header [8]byte
	if strings.Contains(instruction, "[ICMP_REQUEST]") {
		header = [8]byte{0x08, 0x00, 0x12, 0x34, 0x00, 0x00, 0x00, 0x00}
	}
	if strings.Contains(instruction, "[ICMP_UNREACHABLE]") {
		header = [8]byte{0x03, 0x01, 0x12, 0x34, 0x00, 0x00, 0x00, 0x01}
	}
	return header
}

func (generator *PacketGenerator) generateTcpflags(instruction string) uint8 {
	var flags uint8
	if strings.Contains(instruction, "[FIN]") {
		flags |= 0x01
	}
	if strings.Contains(instruction, "[SYN]") {
		flags |= 0x02
	}
	if strings.Contains(instruction, "[RST]") {
		flags |= 0x04
	}
	if strings.Contains(instruction, "[PSH]") {
		flags |= 0x08
	}
	if strings.Contains(instruction, "[ACK]") {
		flags |= 0x10
	}
	if strings.Contains(instruction, "[URG]") {
		flags |= 0x20
	}
	if strings.Contains(instruction, "[ECE]") {
		flags |= 0x40
	}
	if strings.Contains(instruction, "[CWR]") {
		flags |= 0x80
	}
	return flags
}

func (generator *PacketGenerator) GenerateRandomPacket() *Packet {

	assets, err := global.AssetDir("assets/firewall")
	if err != nil {
		return nil
	}
	assetname := assets[rand.Int()%len(assets)]
	asset, err := global.Asset("assets/firewall/" + assetname)
	if err != nil {
		return nil
	}

	var jsonpacket JsonPacket
	json.Unmarshal(asset, &jsonpacket)
	return generator.decodeJsonPacket(jsonpacket)
}

func (generator *PacketGenerator) decodeJsonPacket(jsonpacket JsonPacket) *Packet {
	var packettype int
	switch jsonpacket.PacketType {
	case "icmp":
		packettype = PACKET_ICMP
	case "tcp":
		packettype = PACKET_TCP
	case "udp":
		packettype = PACKET_UDP
	default:
		return nil
	}

	packet := &Packet{
		PacketType: packettype,
		Ipsrc:      generator.generateIP(jsonpacket.Ipsrc),
		Ipdst:      generator.generateIP(jsonpacket.Ipdst),
		SrcPort:    generator.generatePort(jsonpacket.SrcPort),
		DstPort:    generator.generatePort(jsonpacket.DstPort),
		IcmpHeader: generator.generateIcmpHeader(jsonpacket.IcmpHeader),
		Payload:    generator.generatePayload(jsonpacket.Payload),
		Tcpflags:   generator.generateTcpflags(jsonpacket.Tcpflags),
		Harmless:   jsonpacket.Harmless,
	}
	return packet
}

func NewPacketGenerator(dcclassb string) *PacketGenerator {
	packetgenerator := &PacketGenerator{
		dcclassb: dcclassb,
	}
	return packetgenerator
}
