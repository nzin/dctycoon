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

func TestLoanLedger(t *testing.T) {

	// timer start in 1990
	timer := timer.NewGameTimer()
	ledger1 := NewLedger(timer, 0.15, 0.03)
	accounts := ledger1.runLedger()

	// add some equity
	ledger1.AskLoan("loan", time.Date(1991, 1, 1, 0, 0, 0, 0, time.UTC), 1000.0)

	// now we will jump to 1992, and check the end result of 1991
	timer.CurrentTime = time.Date(1992, 1, 1, 0, 0, 0, 0, time.UTC)
	accounts = ledger1.runLedger()

	// refund equity
	timer.CurrentTime = time.Date(1992, 12, 1, 0, 0, 0, 0, time.UTC)
	ledger1.RefundLoan("loan refunded", time.Date(1992, 7, 1, 0, 0, 0, 0, time.UTC), 1000.0)
	accounts = ledger1.runLedger()

	assert.Equal(t, -1000.0, accounts[1991]["16"], "loan from the bank")
	assert.Equal(t, 970.0, accounts[1991]["51"], "equity - loan interest (3% of 1000$)")
	assert.Equal(t, 30.0, accounts[1991]["66"], "loan interest (3% of 1000$)")

	assert.Equal(t, 0.0, accounts[1992]["16"], "loan from the bank")
	assert.Equal(t, -45.0, accounts[1992]["51"], "equity - loan interest (3% of 1000$)")
	assert.Equal(t, 15.0, accounts[1992]["66"], "loan interest (3% of 1000$)")
}
