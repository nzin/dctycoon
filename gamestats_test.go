package dctycoon

import (
	"encoding/json"
	"testing"

	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/dctycoon/timer"

	"github.com/stretchr/testify/assert"
)

func TestGameStats(t *testing.T) {
	gs := NewGameStats()
	gt := timer.NewGameTimer()
	inventory := supplier.NewInventory(gt)
	reputation := supplier.NewReputation()

	sample := `{
		"demandsstats": [
			{"date": "1995-1-16","price": 0.000000,"buyer": "Rob Carlson","servers": [
				{"ramsize": 0, "nbcores":2, "disksize":100, "nb":2},
				{"ramsize": 2048, "nbcores":1, "disksize":0, "nb":1}
				]
			}
		],
		"powerstats": [],
		"reputationstats": []
	}`
	var v map[string]interface{}
	err := json.Unmarshal([]byte(sample), &v)
	assert.Empty(t, err, "correct json payload")
	gs.LoadGame(inventory, reputation, v)

	assert.Equal(t, 1, len(gs.demandsstats), "1 demand stat loaded")
	assert.Equal(t, 2, len(gs.demandsstats[0].serverdemands), "1 demand stat loaded with 2 server confs")

	str := gs.Save()
	assert.NotEmpty(t, str, "save is not empty")
}
