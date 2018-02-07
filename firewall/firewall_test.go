package firewall

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var icmpDummyCode = `
function filterIcmp( ipsrc, ipdst, header, payload)
	return true;
end
`
var icmpDummyCode2 = `
function filterIcmp( ipsrc, ipdst, header, payload)
	if (ipsrc == ipdst) then
		return false;
	end
	if (string.len(payload)>40) then
		return false;
	end
	return true;
end
`

func TestIcmp(t *testing.T) {
	f := NewFirewall()
	ret := f.SubmitIcmp("192,168.1.1", "192.168.2.1", [8]byte{8, 0, 0, 0, 0, 1, 0, 38}, "payload")
	assert.Equal(t, false, ret, "no filterIcmp function")

	err := f.LoadRules(icmpDummyCode)
	assert.Empty(t, err, "icmpDummyCode loaded correctly")
	ret = f.SubmitIcmp("192,168.1.1", "192.168.2.1", [8]byte{8, 0, 0, 0, 0, 1, 0, 38}, "payload")
	assert.Equal(t, true, ret, "dummy filter icmp")

	err = f.LoadRules(icmpDummyCode2)
	assert.Empty(t, err, "icmpDummyCode2 loaded correctly")
	ret = f.SubmitIcmp("192.168.1.1", "192.168.1.1", [8]byte{8, 0, 0, 0, 0, 1, 0, 38}, "payload")
	assert.Equal(t, false, ret, "ipsrc==ipdst")
	ret = f.SubmitIcmp("192.168.1.1", "192.168.1.2", [8]byte{8, 0, 0, 0, 0, 1, 0, 38}, strings.Repeat("A", 100))
	assert.Equal(t, false, ret, "payload > 40")
	ret = f.SubmitIcmp("192.168.1.1", "192.168.1.2", [8]byte{8, 0, 0, 0, 0, 1, 0, 38}, "payload")
	assert.Equal(t, true, ret, "correct packet")
}
