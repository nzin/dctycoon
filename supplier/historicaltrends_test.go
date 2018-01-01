package supplier

import (
	"testing"
	"time"

	"github.com/nzin/dctycoon/timer"
	"github.com/stretchr/testify/assert"
)

func TestMega(t *testing.T) {

	corepercpu := TrendList(initCorepercpu)
	vt := TrendList(initVt)

	assert.Equal(t, int32(2), corepercpu.CurrentValue(time.Date(2006, 2, 1, 0, 0, 0, 0, time.UTC)), " cores per CPU in 2006")
	assert.Equal(t, int32(1), vt.CurrentValue(time.Date(2006, 2, 1, 0, 0, 0, 0, time.UTC)), "VT in 2006")

	emptyarray := make([]interface{}, 0, 0)
	cpuprice := PriceTrendListLoad(emptyarray, cpucorePriceTrend)
	diskprice := PriceTrendListLoad(emptyarray, diskPriceTrend)
	ramprice := PriceTrendListLoad(emptyarray, ramPriceTrend)

	assert.Equal(t, float64(750), cpuprice.CurrentValue(time.Date(1995, time.Month(1), 1, 0, 0, 0, 0, time.UTC)), "CPU price in 1980  (withouth noise)")
	assert.Equal(t, float64(71000), diskprice.CurrentValue(time.Date(1985, time.Month(7), 1, 0, 0, 0, 0, time.UTC)), "Diskprice in 1986 (without noise)")
	assert.Equal(t, float64(50000), ramprice.CurrentValue(time.Date(1995, time.Month(1), 1, 0, 0, 0, 0, time.UTC)), "Ram price in 1986 (withouth noise)")
}

type EventPublisherServiceMock struct{}

func (e *EventPublisherServiceMock) Publish(shortdesc string, longdesc string) {
}

func TestTrendLoad(t *testing.T) {
	gt := timer.NewGameTimer()
	ps := &EventPublisherServiceMock{}
	json := map[string]interface{}{
		"cpupricenoise":  make([]interface{}, 0, 0),
		"diskpricenoise": make([]interface{}, 0, 0),
		"rampricenoise":  make([]interface{}, 0, 0),
	}
	trend := NewTrend()
	trend.Load(json, ps, gt)

	assert.Equal(t, int32(2), trend.Corepercpu.CurrentValue(time.Date(2006, 2, 1, 0, 0, 0, 0, time.UTC)), " cores per CPU in 2006")
	assert.Equal(t, int32(1), trend.Vt.CurrentValue(time.Date(2006, 2, 1, 0, 0, 0, 0, time.UTC)), "VT in 2006")

	badjson := map[string]interface{}{}
	trend.Load(badjson, ps, gt)
	assert.Equal(t, 0, len(trend.Cpuprice.Noise), " empty trends noise loaded")

	badjson2 := map[string]interface{}{
		"cpupricenoise":  1,
		"diskpricenoise": 2,
		"rampricenoise":  3,
	}
	trend.Load(badjson2, ps, gt)
	assert.Equal(t, 0, len(trend.Cpuprice.Noise), " empty trends noise loaded")

	badjson3 := map[string]interface{}{
		"cpupricenoise":  []int{1},
		"diskpricenoise": []int{2},
		"rampricenoise":  []int{3},
	}
	trend.Load(badjson3, ps, gt)
	assert.Equal(t, 0, len(trend.Cpuprice.Noise), " empty trends noise loaded")

	goodjson := map[string]interface{}{
		"cpupricenoise":  []map[string]interface{}{{"pit": "1999-9-9", "value": 0.4}},
		"diskpricenoise": []map[string]interface{}{{"pit": "1999-9-9", "value": 0.4}},
		"rampricenoise":  []map[string]interface{}{{"pit": "1999-9-9", "value": 0.4}},
	}
	trend.Load(goodjson, ps, gt)
	assert.Equal(t, 1, len(trend.Cpuprice.Noise), " not empty trends noise loaded")
}
