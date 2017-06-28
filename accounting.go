package dctycoon

import (
	"time"
	"fmt"
)

//
// payable or reception of money (depending of the sign of Amount
//
type LedgerMovement struct {
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
	Movement(t time.Time) *LedgerMovement
}

//
// object interest to know the ledger status
// essential the Dock
//
type LedgerSubscriber interface {
	CurrentLedgerBalance(balance float64)
}

//
// Ledger
//
type Ledger struct {
	Movements      []LedgerMovement
	CurrentBalance float64
	subscribers    []LedgerSubscriber
}

var GlobalLedger *Ledger

func (self *Ledger) AddMovement(ev LedgerMovement) {
	self.Movements=append(self.Movements,ev)
	self.CurrentBalance+=ev.Amount
	for _,s := range(self.subscribers) {
		s.CurrentLedgerBalance(self.CurrentBalance)
	}
}

func (self *Ledger) AddSubscriber(sub LedgerSubscriber) {
	self.subscribers=append(self.subscribers,sub)
	sub.CurrentLedgerBalance(self.CurrentBalance)
}

func CreateLedger() *Ledger {
	ledger := &Ledger{
		Movements: make([]LedgerMovement,0,100),
		CurrentBalance: 0,
		subscribers: make([]LedgerSubscriber,0),
	}
	return ledger
}

func (self *Ledger) Load(game map[string]interface{}) {
	self.Movements = make([]LedgerMovement,0,100)
	self.CurrentBalance= 0

	for _,event := range(game["movements"].([]interface{})) {
		e := event.(map[string]interface{})
		var year, month, day int
		fmt.Sscanf(e["date"].(string), "%d-%d-%d", &year, &month, &day)
		le:=LedgerMovement{
			Description: e["description"].(string),
			Amount: e["amount"].(float64),
			Date: time.Date(year, time.Month(month), day, 0, 0, 0, 0,  time.UTC),
		}
		self.CurrentBalance+=le.Amount
		self.Movements=append(self.Movements,le)
	}
	for _,s := range(self.subscribers) {
		s.CurrentLedgerBalance(self.CurrentBalance)
	}
}

func (self *Ledger) Save() string {
	str := "{\n"
	str += `"movements": [`
	for i,ev := range(self.Movements) {
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

func (self *Bank) Movement(t time.Time) *LedgerMovement {
	if self.moneyOwned==0 {
		return nil
	}
	return &LedgerMovement{
		Description: "mortage",
		Amount:      self.location.bankinterestrate*self.moneyOwned,
		Date:        t,
	}
}

