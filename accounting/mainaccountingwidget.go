package accounting

import (
	"github.com/nzin/sws"
	//"github.com/veandco/go-sdl2/sdl"
)

type MainAccountingWidget struct {
	rootwidget *sws.RootWidget
	mainwidget *sws.MainWidget
	tabwidget  *sws.TabWidget
	bankwidget *BankWidget
}

func (self *MainAccountingWidget) Show() {
	self.rootwidget.AddChild(self.mainwidget)
	self.rootwidget.SetFocus(self.mainwidget)
}

func (self *MainAccountingWidget) Hide() {
	self.rootwidget.RemoveChild(self.mainwidget)
	children := self.rootwidget.GetChildren()
	if len(children) > 0 {
		self.rootwidget.SetFocus(children[0])
	}
}

//func (self *Accounting) SetBankinterestrate(rate float64) {
//	self.bankwidget.SetBankinterestrate(rate)
//}

func NewMainAccountingWidget(root *sws.RootWidget) *MainAccountingWidget {
	mainwidget := sws.NewMainWidget(650, 400, " Bank and Finance ", true, true)
	tabwidget := sws.NewTabWidget(650, 400)

	ui := &MainAccountingWidget{
		rootwidget: root,
		mainwidget: mainwidget,
		tabwidget:  tabwidget,
	}
	ui.bankwidget = NewBankWidget(root)
	bankScroll := sws.NewScrollWidget(650, 400)
	bankScroll.SetInnerWidget(ui.bankwidget)
	bankScroll.ShowHorizontalScrollbar(false)
	tabwidget.AddTab("Bank", bankScroll)

	balanceScroll := sws.NewScrollWidget(650, 400)
	balanceScroll.SetInnerWidget(NewBalanceWidget())
	balanceScroll.ShowHorizontalScrollbar(false)
	tabwidget.AddTab("Balance", balanceScroll)

	liabilitiesScroll := sws.NewScrollWidget(650, 400)
	liabilitiesScroll.SetInnerWidget(NewLiabilitiesWidget())
	liabilitiesScroll.ShowHorizontalScrollbar(false)
	tabwidget.AddTab("Liabilities", liabilitiesScroll)

	assetScroll := sws.NewScrollWidget(650, 400)
	assetScroll.SetInnerWidget(NewAssetsWidget())
	assetScroll.ShowHorizontalScrollbar(false)
	tabwidget.AddTab("Assets", assetScroll)

	mainwidget.SetCloseCallback(func() {
		ui.Hide()
	})
	mainwidget.SetInnerWidget(tabwidget)
	return ui
}
