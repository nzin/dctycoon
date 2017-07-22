package accounting

import (
	"fmt"
	"github.com/nzin/dctycoon/timer"
	"github.com/nzin/sws"
	//"github.com/veandco/go-sdl2/sdl"
)

// Balance Sheet:
//   Comptes de produits (70)
// - Comptes de charges (6x sauf 66)
// = resultat d'exploitation (EBITDA)
// - amortissement (28)
// - depreciation (29)
// = (EBIT)
// - interet de dettes (66)
// = resultat courant avant impots (EBT)
// - taxes sur les benefices (44 (en fait 444))
// = benefices/resultat net
type BalanceWidget struct {
	FinanceWidget
}

func (self *BalanceWidget) LedgerChange(ledger *Ledger) {
	self.yearN.SetText(fmt.Sprintf("%d (est.)", timer.GlobalGameTimer.CurrentTime.Year()))
	self.yearN1.SetText(fmt.Sprintf("%d", timer.GlobalGameTimer.CurrentTime.Year()-1))
	yearaccountN := GlobalLedger.GetYearAccount(timer.GlobalGameTimer.CurrentTime.Year())
	yearaccountN1 := GlobalLedger.GetYearAccount(timer.GlobalGameTimer.CurrentTime.Year() - 1)
	// forecast account 44
	_, taxN := computeYearlyTaxes(yearaccountN, ledger.taxrate)

	self.lines["70"].N.SetText(fmt.Sprintf("%.2f $", -yearaccountN["70"]))
	self.lines["70"].N1.SetText(fmt.Sprintf("%.2f $", -yearaccountN1["70"]))
	totalSalesN := -yearaccountN["70"]
	totalSalesN1 := -yearaccountN1["70"]

	var total60N float64
	var total60N1 float64
	for _, i := range []string{"60", "61", "62", "63", "64", "65"} {
		self.lines[i].N.SetText(fmt.Sprintf("%.2f $", yearaccountN[i]))
		self.lines[i].N1.SetText(fmt.Sprintf("%.2f $", yearaccountN1[i]))
		total60N += yearaccountN[i]
		total60N1 += yearaccountN1[i]
	}
	self.lines["6"].N.SetText(fmt.Sprintf("%.2f $", total60N))
	self.lines["6"].N1.SetText(fmt.Sprintf("%.2f $", total60N1))

	ebitdaN := totalSalesN + total60N
	ebitdaN1 := totalSalesN1 + total60N1
	self.lines["ebitda"].N.SetText(fmt.Sprintf("%.2f $", ebitdaN))
	self.lines["ebitda"].N1.SetText(fmt.Sprintf("%.2f $", ebitdaN1))

	self.lines["28"].N.SetText(fmt.Sprintf("%.2f $", yearaccountN["28"]))
	self.lines["28"].N1.SetText(fmt.Sprintf("%.2f $", yearaccountN1["28"]))

	self.lines["29"].N.SetText(fmt.Sprintf("%.2f $", yearaccountN["29"]))
	self.lines["29"].N1.SetText(fmt.Sprintf("%.2f $", yearaccountN1["29"]))

	ebitN := ebitdaN - yearaccountN["28"] - yearaccountN["29"]
	ebitN1 := ebitdaN1 - yearaccountN1["28"] - yearaccountN1["29"]
	self.lines["ebit"].N.SetText(fmt.Sprintf("%.2f $", ebitN))
	self.lines["ebit"].N1.SetText(fmt.Sprintf("%.2f $", ebitN1))

	self.lines["66"].N.SetText(fmt.Sprintf("%.2f $", yearaccountN["66"]))
	self.lines["66"].N1.SetText(fmt.Sprintf("%.2f $", yearaccountN1["66"]))

	self.lines["44"].N.SetText(fmt.Sprintf("%.2f $", taxN))
	self.lines["44"].N1.SetText(fmt.Sprintf("%.2f $", yearaccountN1["44"]))

	incomeN := ebitN - yearaccountN["66"] - taxN
	incomeN1 := ebitN1 - yearaccountN1["66"] - yearaccountN1["44"]

	self.lines["income"].N.SetText(fmt.Sprintf("%.2f $", incomeN))
	self.lines["income"].N1.SetText(fmt.Sprintf("%.2f $", incomeN1))
	sws.PostUpdate()
}

func NewBalanceWidget() *BalanceWidget {
	widget := &BalanceWidget{
		FinanceWidget: *NewFinanceWidget(),
	}
	widget.addCategory("Revenue")
	widget.addLine("70", "Sales")
	//widget.addLine("7","Total")
	widget.addSeparator()
	widget.addCategory("Sales product cost")
	widget.addLine("60", "Electricity")
	widget.addLine("61", "Space renting")
	widget.addLine("62", "Advertisement")
	widget.addLine("63", "Salaries taxes")
	widget.addLine("64", "Salaries")
	widget.addLine("65", "Telecom")
	widget.addLine("6", "Total")
	widget.addSeparator()
	widget.addLine("ebitda", "Ebitda")
	widget.addLine("28", "Amortization")
	widget.addLine("29", "Depreciation")
	widget.addSeparator()
	widget.addLine("ebit", "Ebit")
	widget.addLine("66", "Debt interest")
	widget.addLine("44", "Income taxes")
	widget.addSeparator()
	widget.addLine("income", "Net Income")

	GlobalLedger.AddSubscriber(widget)
	return widget
}
