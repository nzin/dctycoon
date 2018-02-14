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
		Date:        time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	accounts = ledger1.runLedger()
	// timer is in 1995, so we just have one year: 1990, without equity
	assert.Equal(t, 1, len(accounts), "empty ledger shouldn't have populated years")
	assert.Equal(t, 0.0, accounts[1995]["51"], "ledger equity is not yet 12000")
	assert.Equal(t, 0.0, accounts[1995]["45"], "ledger equity is not yet 12000")

	// now we will jump to 1999, we have equity poured in 1997
	timer.CurrentTime = time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)
	accounts = ledger1.runLedger()
	assert.Equal(t, 2, len(accounts), "empty ledger shouldn't have populated years")
	assert.Equal(t, 12000.0, accounts[1998]["51"], "ledger equity is not yet 12000")
	assert.Equal(t, -12000.0, accounts[1998]["45"], "ledger equity is not yet 12000")

}

func TestLoanLedger(t *testing.T) {

	// timer start in 1990
	timer := timer.NewGameTimer()
	ledger1 := NewLedger(timer, 0.15, 0.03)
	accounts := ledger1.runLedger()

	// ask for a loan
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

func TestBuyLedger(t *testing.T) {

	// timer start in 1990
	timer := timer.NewGameTimer()
	ledger1 := NewLedger(timer, 0.15, 0.03)
	accounts := ledger1.runLedger()

	// add some equity
	ledger1.AddMovement(LedgerMovement{
		Description: "initial opening",
		Amount:      1000,
		AccountFrom: "4561",
		AccountTo:   "5121",
		Date:        time.Date(1991, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	ledger1.BuyProduct("a product", time.Date(1991, 6, 1, 0, 0, 0, 0, time.UTC), 500)

	// now we will jump to 1996, and check the end result of 1991
	timer.CurrentTime = time.Date(1996, 1, 1, 0, 0, 0, 0, time.UTC)
	accounts = ledger1.runLedger()

	assert.Equal(t, 500.0, accounts[1991]["51"], "500$ remaining after the purchase")
	assert.Equal(t, 97.62773722627738, accounts[1991]["28"], "armotization 1991")
	assert.Equal(t, 166.97080291970804, accounts[1992]["28"], "armotization 1992")
	assert.Equal(t, 166.51459854014598, accounts[1993]["28"], "armotization 1993")
	assert.Equal(t, 68.88686131386861, accounts[1994]["28"], "armotization 1994")

	assert.Equal(t, 500.0, accounts[1996]["51"], "500$ remaining after the purchase")
	assert.Equal(t, 0.0, accounts[1996]["28"], "armotization")
}

func TestDifferentPaymentLedger(t *testing.T) {
	// timer start in 1997
	timer := timer.NewGameTimer()
	ledger := NewLedger(timer, 0.15, 0.03)

	ledger.AddMovement(LedgerMovement{
		Description: "initial opening",
		Amount:      10000,
		AccountFrom: "4561",
		AccountTo:   "5121",
		Date:        time.Date(1997, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	assert.Equal(t, float64(10000), ledger.GetYearAccount(1997)["51"], "10000 in cash")

	ledger.PayLandlord(12, 1000, timer.CurrentTime)
	assert.Equal(t, float64(9000), ledger.GetYearAccount(1997)["51"], "9000 in cash")

	ledger.PayMapUpgrade(500, timer.CurrentTime)
	assert.Equal(t, float64(8500), ledger.GetYearAccount(1997)["51"], "8500 in cash")

	ledger.PayMapUpgrade(500, timer.CurrentTime)
	assert.Equal(t, float64(8000), ledger.GetYearAccount(1997)["51"], "8000 in cash")

	assert.Equal(t, float64(2000), ledger.GetYearAccount(1997)["61"], "2000 payed")
}
