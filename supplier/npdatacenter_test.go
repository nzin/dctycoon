package supplier

import (
	"encoding/json"
	"testing"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/timer"
	"github.com/stretchr/testify/assert"
)

func TestNPDatacenter(t *testing.T) {
	gt := timer.NewGameTimer()
	ps := &EventPublisherServiceMock{}
	j := map[string]interface{}{
		"cpupricenoise":  make([]interface{}, 0, 0),
		"diskpricenoise": make([]interface{}, 0, 0),
		"rampricenoise":  make([]interface{}, 0, 0),
	}
	trend := TrendLoad(j, ps, gt)

	data, err := global.Asset("assets/npdatacenter/mono_r100_r200.json")
	assert.Empty(t, err, "load mono_r100_r200 profile asset")

	profile := make(map[string]BuyoutProfile)
	err = json.Unmarshal(data, &profile)
	assert.Empty(t, err, "parse mono_r100_r200 profile asset")

	npd := NewNPDatacenter(gt, trend, 10000, "siliconvalley", "mono_r100_r200.json")
	assert.NotEmpty(t, npd, "NPDatacenter mono_r100_r200 profile loaded")
	assert.Equal(t, "R100", npd.buyoutprofile["R100physical"].Servertype, "check if the profile is correctly loaded")
}
