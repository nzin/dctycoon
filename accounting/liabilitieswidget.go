package accounting

import (
	"fmt"
	//"github.com/nzin/sws"
	"github.com/nzin/dctycoon/timer"
	//"github.com/veandco/go-sdl2/sdl"
)

// PASSIF (année N)
//
// capitaux propre (45) + debt interest payed (46)
// resultat de l'exercice ??
// dettes (16)
// = passif
type LiabilitiesWidget struct {
	FinanceWidget
	timer  *timer.GameTimer
	ledger *Ledger
}

func (self *LiabilitiesWidget) LedgerChange() {
	self.yearN.SetText(fmt.Sprintf("%d (est.)", self.timer.CurrentTime.Year()))
	self.yearN1.SetText(fmt.Sprintf("%d", self.timer.CurrentTime.Year()-1))
	yearaccountN := self.ledger.GetYearAccount(self.timer.CurrentTime.Year())
	yearaccountN1 := self.ledger.GetYearAccount(self.timer.CurrentTime.Year() - 1)
	// forecast account 44
	_, taxN := computeYearlyTaxes(yearaccountN, self.ledger.taxrate)

	self.lines["45"].N.SetText(fmt.Sprintf("%.2f $", -yearaccountN["45"]-yearaccountN["46"]))
	self.lines["45"].N1.SetText(fmt.Sprintf("%.2f $", -yearaccountN1["45"]-yearaccountN1["46"]))

	var profitN, profitN1 float64
	profitN = -yearaccountN["70"]
	profitN1 = -yearaccountN1["70"]
	for _, i := range []string{"60", "61", "62", "63", "64", "65"} {
		profitN -= yearaccountN[i]
		profitN1 -= yearaccountN1[i]
	}
	profitN -= yearaccountN["28"]
	profitN1 -= yearaccountN1["28"]
	profitN -= yearaccountN["29"]
	profitN1 -= yearaccountN1["29"]
	// debt interest
	profitN -= yearaccountN["66"]
	profitN1 -= yearaccountN1["66"]
	profitN -= taxN
	profitN1 -= yearaccountN1["44"]

	self.lines["profit"].N.SetText(fmt.Sprintf("%.2f $", profitN))
	self.lines["profit"].N1.SetText(fmt.Sprintf("%.2f $", profitN1))

	self.lines["16"].N.SetText(fmt.Sprintf("%.2f $", -yearaccountN["16"]))
	self.lines["16"].N1.SetText(fmt.Sprintf("%.2f $", -yearaccountN1["16"]))

	liabN := -yearaccountN["45"] - yearaccountN["46"] - yearaccountN["16"] + profitN
	liabN1 := -yearaccountN1["45"] - yearaccountN1["46"] - yearaccountN1["16"] + profitN1

	self.lines["liabilities"].N.SetText(fmt.Sprintf("%.2f $", liabN))
	self.lines["liabilities"].N1.SetText(fmt.Sprintf("%.2f $", liabN1))
}

func NewLiabilitiesWidget() *LiabilitiesWidget {
	widget := &LiabilitiesWidget{
		FinanceWidget: *NewFinanceWidget(),
		timer:         nil,
		ledger:        nil,
	}
	widget.addCategory("Liabilities")
	widget.addLine("45", "Capital")
	widget.addLine("profit", "Profit/Lost")
	widget.addLine("16", "Debts")
	widget.addSeparator()
	widget.addLine("liabilities", "Total")

	return widget
}

func (self *LiabilitiesWidget) SetGame(timer *timer.GameTimer, ledger *Ledger) {
	self.timer = timer
	if self.ledger != nil {
		self.ledger.RemoveSubscriber(self)
	}
	self.ledger = ledger
	ledger.AddSubscriber(self)
}
