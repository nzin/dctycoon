package dctycoon

import (
	"time"
	"fmt"
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
// - bank loan interest
// - DC location rent
// - internet fiber bill
// - electricity bill
// - salaries
// - maintenances?
//
// return nil if nothing to declare
//
type StakeHolderMonthly interface {
	Movement(t time.Time) *LedgerEvent
}

//
// Ledger
//
type Ledger struct {
	Events         []LedgerEvent
	CurrentBalance float64
}

var globalLedger *Ledger

func (self *Ledger) AddEvent(ev LedgerEvent) {
	self.Events=append(self.Events,ev)
	self.CurrentBalance+=ev.Amount
}

func LedgerLoad(game map[string]interface{}) *Ledger {
	ledger := &Ledger{
		Events: make([]LedgerEvent,0,100),
		CurrentBalance: 0,
	}
	for _,event := range(game["events"].([]interface{})) {
		e := event.(map[string]interface{})
		var year, month, day int
		fmt.Sscanf(e["date"].(string), "%d-%d-%d", &year, &month, &day)
		le:=LedgerEvent{
			Description: e["description"].(string),
			Amount: e["amount"].(float64),
			Date: time.Date(year, time.Month(month), day, 0, 0, 0, 0,  time.UTC),
		}
		ledger.CurrentBalance+=le.Amount
		ledger.Events=append(ledger.Events,le)
	}
	return ledger
}

func (self *Ledger) Save() string {
	str := "{\n"
	str += `"events": [`
	for i,ev := range(self.Events) {
		if i==0 {
			str+=fmt.Sprintf(`\n  {"description":"%s", "amount": %v , "date":"%d-%d-%d"}`, ev.Description, ev.Amount, ev.Date.Year(), ev.Date.Month(), ev.Date.Day())
		} else {
			str+=fmt.Sprintf(`\n,{"description":"%s", "amount": %v , "date":"%d-%d-%d"}`, ev.Description, ev.Amount, ev.Date.Year(), ev.Date.Month(), ev.Date.Day())
		}
	}
	str += "\n]}"
	return str
}

//
// Bank stakeholder
//
type Bank struct {
	location   LocationType
	moneyOwned float64
}

func (self *Bank) Movement(t time.Time) *LedgerEvent {
	if self.moneyOwned==0 {
		return nil
	}
	return &LedgerEvent{
		Description: "mortage",
		Amount:      self.location.bankinterestrate*self.moneyOwned,
		Date:        t,
	}
}

