package accounting

import (
	"fmt"
	//"github.com/nzin/sws"
	"github.com/nzin/dctycoon/timer"
	//"github.com/veandco/go-sdl2/sdl"
)

// PASSIF (année N)
//
// capitaux propre (45)
// resultat de l'exercice ??
// dettes (16)
// = passif
type LiabilitiesWidget struct {
	FinanceWidget
}

func (self *LiabilitiesWidget) LedgerChange(ledger *Ledger) {
	self.yearN.SetText(fmt.Sprintf("%d (est.)", timer.GlobalGameTimer.CurrentTime.Year()))
	self.yearN1.SetText(fmt.Sprintf("%d", timer.GlobalGameTimer.CurrentTime.Year()-1))
	yearaccountN := GlobalLedger.GetYearAccount(timer.GlobalGameTimer.CurrentTime.Year())
	yearaccountN1 := GlobalLedger.GetYearAccount(timer.GlobalGameTimer.CurrentTime.Year() - 1)
	// forecast account 44
	_, taxN := computeYearlyTaxes(yearaccountN, ledger.taxrate)

	self.lines["45"].N.SetText(fmt.Sprintf("%.2f $", -yearaccountN["45"]))
	self.lines["45"].N1.SetText(fmt.Sprintf("%.2f $", -yearaccountN1["45"]))

	var profitN, profitN1 float64
	profitN = -yearaccountN["70"]
	profitN1 = -yearaccountN1["70"]
	for _, i := range []string{"60", "61", "62", "63", "64", "65", "66"} {
		profitN += yearaccountN[i]
		profitN1 += yearaccountN1[i]
	}
	profitN -= yearaccountN["28"]
	profitN1 -= yearaccountN1["28"]
	profitN -= yearaccountN["29"]
	profitN1 -= yearaccountN1["29"]
	profitN -= taxN
	profitN1 -= yearaccountN1["44"]

	self.lines["profit"].N.SetText(fmt.Sprintf("%.2f $", profitN))
	self.lines["profit"].N1.SetText(fmt.Sprintf("%.2f $", profitN1))

	self.lines["16"].N.SetText(fmt.Sprintf("%.2f $", yearaccountN["16"]))
	self.lines["16"].N1.SetText(fmt.Sprintf("%.2f $", yearaccountN1["16"]))

	liabN := -yearaccountN["45"] + yearaccountN["16"] + profitN
	liabN1 := -yearaccountN1["45"] + yearaccountN1["16"] + profitN1

	self.lines["liabilities"].N.SetText(fmt.Sprintf("%.2f $", liabN))
	self.lines["liabilities"].N1.SetText(fmt.Sprintf("%.2f $", liabN1))
}

func NewLiabilitiesWidget() *LiabilitiesWidget {
	widget := &LiabilitiesWidget{
		FinanceWidget: *NewFinanceWidget(),
	}
	widget.addCategory("Liabilities")
	widget.addLine("45", "Capital")
	widget.addLine("profit", "Profit/Lost")
	widget.addLine("16", "Debts")
	widget.addSeparator()
	widget.addLine("liabilities", "Total")

	GlobalLedger.AddSubscriber(widget)
	return widget
}
