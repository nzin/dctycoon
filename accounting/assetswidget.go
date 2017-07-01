package accounting

import(
//	"fmt"
	"github.com/nzin/sws"
//	"github.com/nzin/dctycoon/timer"
	//"github.com/veandco/go-sdl2/sdl"
)

//
// Page Shop
//
type AssetsWidget struct {
	sws.SWS_CoreWidget
}

func (self *AssetsWidget) LedgerChange(ledger *Ledger) {
//	yearaccount := GlobalLedger.GetYearAccount(timer.GlobalGameTimer.CurrentTime.Year())
}

func CreateAssetsWidget() *AssetsWidget {
	widget:=&AssetsWidget{
		SWS_CoreWidget: *sws.CreateCoreWidget(300,100),
	}
	
	GlobalLedger.AddSubscriber(widget)
	return widget
}
