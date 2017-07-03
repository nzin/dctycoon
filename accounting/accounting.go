package accounting

import (
	"github.com/nzin/sws"
	//"github.com/veandco/go-sdl2/sdl"
)

type Accounting struct {
	rootwidget *sws.SWS_RootWidget
	mainwidget *sws.SWS_MainWidget
	tabwidget  *sws.SWS_TabWidget
	bankwidget *BankWidget
}

func (self *Accounting) Show() {
	self.rootwidget.AddChild(self.mainwidget)
	self.rootwidget.SetFocus(self.mainwidget)
}

func (self *Accounting) Hide() {
	self.rootwidget.RemoveChild(self.mainwidget)
	children := self.rootwidget.GetChildren()
	if len(children) > 0 {
		self.rootwidget.SetFocus(children[0])
	}
}

func (self *Accounting) SetBankinterestrate(rate float64) {
	self.bankwidget.SetBankinterestrate(rate)
}

func CreateAccounting(root *sws.SWS_RootWidget) *Accounting {
	mainwidget := sws.CreateMainWidget(650, 400, " Bank and Finance ", true, true)
	tabwidget := sws.CreateTabWidget(650, 400)
	
	ui := &Accounting{
		rootwidget: root,
		mainwidget: mainwidget,
		tabwidget:  tabwidget,
	}
	ui.bankwidget=CreateBankWidget()
	bankScroll:=sws.CreateScrollWidget(650,400)
	bankScroll.SetInnerWidget(ui.bankwidget)
	bankScroll.ShowHorizontalScrollbar(false)
	tabwidget.AddTab("Bank",bankScroll)
	
	balanceScroll:=sws.CreateScrollWidget(650,400)
	balanceScroll.SetInnerWidget(CreateBalanceWidget())
	balanceScroll.ShowHorizontalScrollbar(false)
	tabwidget.AddTab("Balance",balanceScroll)
	
	liabilitiesScroll:=sws.CreateScrollWidget(650,400)
	liabilitiesScroll.SetInnerWidget(CreateLiabilitiesWidget())
	liabilitiesScroll.ShowHorizontalScrollbar(false)
	tabwidget.AddTab("Liabilities",liabilitiesScroll)
	
	assetScroll:=sws.CreateScrollWidget(650,400)
	assetScroll.SetInnerWidget(CreateAssetsWidget())
	assetScroll.ShowHorizontalScrollbar(false)
	tabwidget.AddTab("Assets",assetScroll)
	
	mainwidget.SetCloseCallback(func() {
		ui.Hide()
	})
	mainwidget.SetInnerWidget(tabwidget)
	return ui
}
