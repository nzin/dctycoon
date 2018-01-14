package dctycoon

import (
	"fmt"
	"testing"

	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/sws"
	"github.com/stretchr/testify/assert"
)

func TestDemandAttribution(t *testing.T) {
	quit := false
	root := &sws.RootWidget{}

	game := NewGame(&quit, root)
	game.InitGame("siliconvalley", DIFFICULTY_EASY)

	assert.NotEmpty(t, game, "Game created")
	assert.Equal(t, 3, len(game.npactors), "3 NPC created")
	assert.NotEmpty(t, game.player, "1 player created")

	// reduce to one opponent
	opponent := NewNPDatacenter()
	opponent.Init(game.timer, 30000, "siliconvalley", game.trends, "mono_r100_r200.json")
	game.npactors = make([]*NPDatacenter, 1, 1)
	game.npactors[0] = opponent

	// reduce demand template to one template
	demandtemplate := supplier.DemandTemplateAssetLoad("001_basicwebserver.json")
	assert.NotEmpty(t, demandtemplate, "read a basic demand template")

	// generate new demand
	actors := make([]supplier.Actor, 0, 0)
	actors = append(actors, opponent)
	actors = append(actors, game.player)
	demand := demandtemplate.InstanciateDemand()
	serverbundle := demand.FindOffer(actors, game.timer.CurrentTime)
	assert.Empty(t, serverbundle, "demand created but not atributed")

	// add server in inventory
	opponent.NewYearOperations()
	assert.Equal(t, 1, len(opponent.inventory.Items), "new year passed, we bought some servers")
	for i := 0; i < 10; i++ {
		game.timer.TimerClock()
	}
	fmt.Println(opponent.inventory.Items[0])
	fmt.Println(game.timer.CurrentTime)
	fmt.Println(opponent.inventory.GetOffers()[0])
	serverbundle = demand.FindOffer(actors, game.timer.CurrentTime)
	assert.NotEmpty(t, serverbundle, "demand created and attributed")
}