package dctycoon

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/dctycoon/timer"
	"github.com/stretchr/testify/assert"
)

type EventPublisherServiceMock struct{}

func (e *EventPublisherServiceMock) Publish(shortdesc string, longdesc string) {
}

func TestNPDatacenter(t *testing.T) {
	gt := timer.NewGameTimer()
	ps := &EventPublisherServiceMock{}
	j := map[string]interface{}{
		"cpupricenoise":  make([]interface{}, 0, 0),
		"diskpricenoise": make([]interface{}, 0, 0),
		"rampricenoise":  make([]interface{}, 0, 0),
	}
	trend := supplier.TrendLoad(j, ps, gt)

	data, err := global.Asset("assets/npdatacenter/mono_r100_r200.json")
	assert.Empty(t, err, "load mono_r100_r200 profile asset")

	profile := make(map[string]BuyoutProfile)
	err = json.Unmarshal(data, &profile)
	assert.Empty(t, err, "parse mono_r100_r200 profile asset")

	npd := NewNPDatacenter(gt, trend, 10000, "siliconvalley", "mono_r100_r200.json")
	assert.NotEmpty(t, npd, "NPDatacenter mono_r100_r200 profile loaded")
	assert.Equal(t, "R100", npd.buyoutprofile["R100physical"].Servertype, "check if the profile is correctly loaded")
}

func TestNPDatacenterBuyout(t *testing.T) {
	gt := timer.NewGameTimer()
	gt.CurrentTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	ps := &EventPublisherServiceMock{}
	j := map[string]interface{}{
		"cpupricenoise":  make([]interface{}, 0, 0),
		"diskpricenoise": make([]interface{}, 0, 0),
		"rampricenoise":  make([]interface{}, 0, 0),
	}
	trend := supplier.TrendLoad(j, ps, gt)

	npd := NewNPDatacenter(gt, trend, 10000, "siliconvalley", "mono_r100_r200.json")
	assert.NotEmpty(t, npd, "NPDatacenter mono_r100_r200 profile loaded")
	assert.Equal(t, "R100", npd.buyoutprofile["R100physical"].Servertype, "check if the profile is correctly loaded")

	npd.NewYearOperations()
	assert.Equal(t, 1, len(npd.inventory.Items), "new year passed, we bought some servers")
	assert.Equal(t, 1, len(npd.inventory.GetOffers()), "we have one offer for R100 server")
	assert.Equal(t, float64(2700.0), npd.inventory.GetOffers()[0].Price, "R100 is priced as 2400$")
}
