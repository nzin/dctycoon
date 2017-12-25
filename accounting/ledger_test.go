package accounting

import (
	"testing"
	"time"

	"github.com/nzin/dctycoon/timer"
	"github.com/stretchr/testify/assert"
)

func TestLedger(t *testing.T) {

	// timer start in 1990
	timer := timer.NewGameTimer()
	ledger1 := NewLedger(timer, 0.15, 0.03)
	accounts := ledger1.runLedger()

	assert.Equal(t, 0, len(accounts), "empty ledger shouldn't have populated years")

	// add some equity
	ledger1.AddMovement(LedgerMovement{
		Description: "initial opening",
		Amount:      12000,
		AccountFrom: "4561",
		AccountTo:   "5121",
		Date:        time.Date(1991, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	accounts = ledger1.runLedger()
	// timer is in 1990, so we just have one year: 1990, without equity
	assert.Equal(t, 1, len(accounts), "empty ledger shouldn't have populated years")
	assert.Equal(t, 0.0, accounts[1990]["51"], "ledger equity is not yet 12000")
	assert.Equal(t, 0.0, accounts[1990]["45"], "ledger equity is not yet 12000")

	// now we will jump to 1994, we have equity poured in 1991
	timer.CurrentTime = time.Date(1994, 1, 1, 0, 0, 0, 0, time.UTC)
	accounts = ledger1.runLedger()
	assert.Equal(t, 4, len(accounts), "empty ledger shouldn't have populated years")
	assert.Equal(t, 12000.0, accounts[1993]["51"], "ledger equity is not yet 12000")
	assert.Equal(t, -12000.0, accounts[1993]["45"], "ledger equity is not yet 12000")

}
