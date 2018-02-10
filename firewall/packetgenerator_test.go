package firewall

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/nzin/dctycoon/global"
	"github.com/stretchr/testify/assert"
)

func TestPackets(t *testing.T) {
	generator := NewPacketGenerator("18.2")

	ip := generator.generateIP("[IP_IN]")
	assert.Equal(t, true, strings.HasPrefix(ip, "18.2"), "generate datacenter internal IP")

	pingofdeath := generator.generatePayload("[65K]")
	assert.Equal(t, 65536, len(pingofdeath), "Ping of death packet is > 64K")

	asset, err := global.Asset("assets/firewall/pingNormal1.json")
	assert.NoError(t, err, "pingNormal1.json loaded")

	var jsonpacket JsonPacket
	json.Unmarshal(asset, &jsonpacket)
	packet := generator.decodeJsonPacket(jsonpacket)

	assert.NotEmpty(t, packet, "packet generated")
	fmt.Println(packet.Ipdst)
	assert.Equal(t, true, strings.HasPrefix(packet.Ipdst, "18.2"), "to Datacenter IP")
	assert.Equal(t, uint8(8), packet.IcmpHeader[0], "ICMP request")
}
