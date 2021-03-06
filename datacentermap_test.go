package dctycoon

import (
	"encoding/json"
	"testing"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/dctycoon/timer"
	"github.com/stretchr/testify/assert"
)

func TestMigrationMap(t *testing.T) {
	gt := timer.NewGameTimer()
	inventory := supplier.NewInventory(gt)
	datacenter := NewDatacenterMap()
	datacenter.SetGame(inventory, supplier.AvailableLocation["siliconvalley"], gt.CurrentTime)

	datamap1, err := global.Asset("assets/dcmap/3_4_room.json")
	assert.Empty(t, err, "load 3_4_room.json map")

	map1 := make(map[string]interface{})
	err = json.Unmarshal(datamap1, &map1)

	assert.Empty(t, err, "unmarshall 3_4_room.json map")

	datacenter.LoadMap(map1)

	assert.Equal(t, int32(12), datacenter.GetMeterSquares(), "12 m2 room")

	datamap2, err := global.Asset("assets/dcmap/20_20_standard.json")
	assert.Empty(t, err, "load 20_20_standard.json map")

	map2 := make(map[string]interface{})
	err = json.Unmarshal(datamap2, &map2)

	assert.Empty(t, err, "unmarshall 20_20_standard.json map")

	datacenter.MigrateToMap(map2)
}
