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
	trend := supplier.NewTrend()
	trend.Load(j, ps, gt)

	data, err := global.Asset("assets/npdatacenter/mono_r100_r200.json")
	assert.Empty(t, err, "load mono_r100_r200 profile asset")

	profile := make(map[string]BuyoutProfile)
	err = json.Unmarshal(data, &profile)
	assert.Empty(t, err, "parse mono_r100_r200 profile asset")

	npd := NewNPDatacenter()
	npd.Init(gt, 10000, "siliconvalley", trend, "mono_r100_r200.json", "John Doe", true)
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
	trend := supplier.NewTrend()
	trend.Load(j, ps, gt)

	npd := NewNPDatacenter()
	npd.Init(gt, 20000, "siliconvalley", trend, "mono_r100_r200.json", "John Doe", true)
	assert.NotEmpty(t, npd, "NPDatacenter mono_r100_r200 profile loaded")
	assert.Equal(t, "R100", npd.buyoutprofile["R100physical"].Servertype, "check if the profile is correctly loaded")

	npd.NewYearOperations()
	assert.Equal(t, 1, len(npd.GetInventory().Items), "new year passed, we bought some servers")
	assert.Equal(t, 1, len(npd.GetOffers()), "we have one offer for R100 server")
	assert.Equal(t, float64(126.49499999999999), npd.GetOffers()[0].Price, "offer is priced as 81$")

	save := npd.Save()

	npd = NewNPDatacenter()
	var data interface{}
	err := json.Unmarshal([]byte(save), &data)
	assert.NoError(t, err, "unmarshall correctly npd save")
	npd.LoadGame(gt, trend, data.(map[string]interface{}))
	assert.Equal(t, 1, len(npd.GetInventory().Items), "reload: we still have 1 item in the inventory")
}
