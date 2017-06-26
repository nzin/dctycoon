package dctycoon

package (
	"time"
)

//
// payable or reception of money (depending of the sign of Amount
//
type LedgerEvent {
	Description string
	Amount      float64
	Date        time.Time
}

//
// this interface is used by any stakeholder
// to invoice (supplier) or pay (customer) the DC
//
type StakeHolderMonthly interface {
	movement(t time.Time) LedgerEvent
}

//
// Ledger
//
type Ledger struct {
	Events         []LedgerEvent
	CurrentBalance float64
}

func (self *Ledger) AddEvent(ev LedgerEvent) {
	self.events=append(self.events,ev)
	self.currentBalance+=ev.Amount
}

//
// Bank stakeholder
//
type Bank struct {
	location   LocationType
	moneyOwned float64
}

func (self *Bank) movement(t time.Time) LedgerEvent (
	return LedgerEvent{
	        Description: "mortage",
        	Amount:      self.location.bankinterestrate*self.moneyOwned,
        	Date:        t,
	}
}

