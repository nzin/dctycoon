package accounting

import (
	"fmt"

	//"github.com/nzin/sws"
	"github.com/nzin/dctycoon/timer"
	//"github.com/veandco/go-sdl2/sdl"
)

// ACTIF (année N)
//
// Terrain (2312)
// installation techniques (2315)
// autre immo (37)
// - amortissements (28)
// = actif immobilisé
// Disponibilité (51)
// = actif circulant
//
// Total=
type AssetsWidget struct {
	FinanceWidget
	timer  *timer.GameTimer
	ledger *Ledger
}

func (self *AssetsWidget) LedgerChange() {
	self.yearN.SetText(fmt.Sprintf("%d (est.)", self.timer.CurrentTime.Year()))
	self.yearN1.SetText(fmt.Sprintf("%d", self.timer.CurrentTime.Year()-1))
	yearaccountN := self.ledger.GetYearAccount(self.timer.CurrentTime.Year())
	yearaccountN1 := self.ledger.GetYearAccount(self.timer.CurrentTime.Year() - 1)
	// forecast account 44
	_, taxN := computeYearlyTaxes(yearaccountN, self.ledger.taxrate)

	self.lines["23"].N.SetText(fmt.Sprintf("%.2f $", yearaccountN["23"]+yearaccountN["28"]))
	self.lines["23"].N1.SetText(fmt.Sprintf("%.2f $", yearaccountN1["23"]+yearaccountN1["28"]))
	self.lines["37"].N.SetText(fmt.Sprintf("%.2f $", yearaccountN["37"]))
	self.lines["37"].N1.SetText(fmt.Sprintf("%.2f $", yearaccountN1["37"]))
	self.lines["28"].N.SetText(fmt.Sprintf("%.2f $", -yearaccountN["28"]))
	self.lines["28"].N1.SetText(fmt.Sprintf("%.2f $", -yearaccountN1["28"]))

	self.lines["51"].N.SetText(fmt.Sprintf("%.2f $", yearaccountN["51"]-taxN))
	self.lines["51"].N1.SetText(fmt.Sprintf("%.2f $", yearaccountN1["51"]))

	assetsN := yearaccountN["23"] + yearaccountN["37"] + yearaccountN["51"] - taxN
	assetsN1 := yearaccountN1["23"] + yearaccountN1["37"] + yearaccountN1["51"]
	self.lines["assets"].N.SetText(fmt.Sprintf("%.2f $", assetsN))
	self.lines["assets"].N1.SetText(fmt.Sprintf("%.2f $", assetsN1))
}

func NewAssetsWidget() *AssetsWidget {
	widget := &AssetsWidget{
		FinanceWidget: *NewFinanceWidget(),
		timer:         nil,
		ledger:        nil,
	}
	widget.addCategory("Immobilized Assets")
	widget.addLine("23", "Immobilization")
	widget.addLine("37", "Other immo")
	widget.addLine("28", "Amortization")
	widget.addSeparator()
	widget.addCategory("Fluid Assets")
	widget.addLine("51", "Bank account (est.)")
	widget.addSeparator()
	widget.addLine("assets", "Total")

	return widget
}

func (self *AssetsWidget) SetGame(timer *timer.GameTimer, ledger *Ledger) {
	self.timer = timer
	if self.ledger != nil {
		self.ledger.RemoveSubscriber(self)
	}
	self.ledger = ledger
	ledger.AddSubscriber(self)
}
