package accounting

import(
//	"fmt"
//	"github.com/nzin/sws"
//	"github.com/nzin/dctycoon/timer"
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
//	yearaccountN := GlobalLedger.GetYearAccount(timer.GlobalGameTimer.CurrentTime.Year())
//	yearaccountN1 := GlobalLedger.GetYearAccount(timer.GlobalGameTimer.CurrentTime.Year()-1)
}

func CreateBalanceWidget() *BalanceWidget {
	widget:=&BalanceWidget{
		FinanceWidget: *CreateFinanceWidget(),
	}
	widget.addCategory("Revenue")
	widget.addLine("70","Sales")
	widget.addSeparator()
	widget.addCategory("Sales product cost")
	widget.addLine("60","Sales")
	
	GlobalLedger.AddSubscriber(widget)
	return widget
}
