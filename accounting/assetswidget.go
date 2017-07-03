package accounting

import(
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
}

func (self *AssetsWidget) LedgerChange(ledger *Ledger) {
	self.yearN.SetText(fmt.Sprintf("%d", timer.GlobalGameTimer.CurrentTime.Year()))
	self.yearN1.SetText(fmt.Sprintf("%d", timer.GlobalGameTimer.CurrentTime.Year()-1))
	yearaccountN := GlobalLedger.GetYearAccount(timer.GlobalGameTimer.CurrentTime.Year())
	yearaccountN1 := GlobalLedger.GetYearAccount(timer.GlobalGameTimer.CurrentTime.Year()-1)

	self.lines["23"].N.SetText(fmt.Sprintf("%.2f $",yearaccountN["23"]))
	self.lines["23"].N1.SetText(fmt.Sprintf("%.2f $",yearaccountN1["23"]))
	self.lines["37"].N.SetText(fmt.Sprintf("%.2f $",yearaccountN["37"]))
	self.lines["37"].N1.SetText(fmt.Sprintf("%.2f $",yearaccountN1["37"]))
	self.lines["28"].N.SetText(fmt.Sprintf("%.2f $",yearaccountN["28"]))
	self.lines["28"].N1.SetText(fmt.Sprintf("%.2f $",yearaccountN1["28"]))

	self.lines["51"].N.SetText(fmt.Sprintf("%.2f $",yearaccountN["51"]))
	self.lines["51"].N1.SetText(fmt.Sprintf("%.2f $",yearaccountN1["51"]))
	
	assetsN:=yearaccountN["23"]+yearaccountN["37"]-yearaccountN["28"]+yearaccountN["51"]
	assetsN1:=yearaccountN1["23"]+yearaccountN1["37"]-yearaccountN1["28"]+yearaccountN1["51"]
	self.lines["assets"].N.SetText(fmt.Sprintf("%.2f $",assetsN))
	self.lines["assets"].N1.SetText(fmt.Sprintf("%.2f $",assetsN1))
}

func CreateAssetsWidget() *AssetsWidget {
	widget:=&AssetsWidget{
		FinanceWidget: *CreateFinanceWidget(),
	}
	widget.addCategory("Immobilized Assets")
	widget.addLine("23","Immobilization")
	widget.addLine("37","Other immo")
	widget.addLine("28","Amortization")
	widget.addSeparator()
	widget.addCategory("Fluid Assets")
	widget.addLine("51","Bank account")
	widget.addSeparator()
	widget.addLine("assets","Total")

	
	GlobalLedger.AddSubscriber(widget)
	return widget
}
