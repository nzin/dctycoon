package dctycoon

import (
	"time"
)

//
// payable or reception of money (depending of the sign of Amount
//
type LedgerEvent struct {
	Description string
	Amount      float64
	Date        time.Time
}

//
// this interface is used by any stakeholder
// to invoice (supplier) or pay (customer) the DC
//
type StakeHolderMonthly interface {
	Movement(t time.Time) LedgerEvent
}

//
// Ledger
//
type Ledger struct {
	Events         []LedgerEvent
	CurrentBalance float64
}

func (self *Ledger) AddEvent(ev LedgerEvent) {
	self.Events=append(self.Events,ev)
	self.CurrentBalance+=ev.Amount
}

//
// Bank stakeholder
//
type Bank struct {
	location   LocationType
	moneyOwned float64
}

func (self *Bank) Movement(t time.Time) LedgerEvent {
	return LedgerEvent{
		Description: "mortage",
		Amount:      self.location.bankinterestrate*self.moneyOwned,
		Date:        t,
	}
}

