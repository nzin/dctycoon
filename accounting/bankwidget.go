package accounting

import (
	"fmt"
	"github.com/nzin/dctycoon/timer"
	"github.com/nzin/sws"
	//"github.com/veandco/go-sdl2/sdl"
)

//
// Page Shop
//
type BankWidget struct {
	sws.CoreWidget
	paybackbutton   *sws.ButtonWidget
	askloanbutton   *sws.ButtonWidget
	accountPosition *sws.LabelWidget
	accountDebt     *sws.LabelWidget
	interestRate    float64
	interestRateL   *sws.LabelWidget
	currentInterest *sws.LabelWidget
}

//func (self *BankWidget) SetBankinterestrate(rate float64) {
//	self.interestRate = rate
//	self.interestRateL.SetText(fmt.Sprintf("%.2f %%/y", rate*100))
//	sws.PostUpdate()
//}

func (self *BankWidget) LedgerChange(ledger *Ledger) {
	yearaccount := GlobalLedger.GetYearAccount(timer.GlobalGameTimer.CurrentTime.Year())
	self.interestRate = ledger.loanrate
	self.interestRateL.SetText(fmt.Sprintf("%.2f %%/y", ledger.loanrate*100))
	self.accountPosition.SetText(fmt.Sprintf("%.2f $", yearaccount["51"]))
	self.accountDebt.SetText(fmt.Sprintf("%.2f $", -yearaccount["16"]))
	self.currentInterest.SetText(fmt.Sprintf("%.2f $/y", -yearaccount["16"]*self.interestRate))
	sws.PostUpdate()
}

func NewBankWidget() *BankWidget {
	bankwidget := &BankWidget{
		CoreWidget: *sws.NewCoreWidget(350, 180),
	}

	title := sws.NewLabelWidget(190, 25, "Your Bank account")
	title.Move(10, 10)
	bankwidget.AddChild(title)

	accountPositionTitle := sws.NewLabelWidget(190, 25, "Your current position")
	accountPositionTitle.Move(10, 40)
	bankwidget.AddChild(accountPositionTitle)

	accountPosition := sws.NewLabelWidget(100, 25, "0 $")
	accountPosition.Move(200, 40)
	bankwidget.AddChild(accountPosition)
	bankwidget.accountPosition = accountPosition

	accountDebtTitle := sws.NewLabelWidget(190, 25, "Your current debt")
	accountDebtTitle.Move(10, 65)
	bankwidget.AddChild(accountDebtTitle)

	accountDebt := sws.NewLabelWidget(100, 25, "0 $")
	accountDebt.Move(200, 65)
	bankwidget.AddChild(accountDebt)
	bankwidget.accountDebt = accountDebt

	interestRate := sws.NewLabelWidget(100, 25, "0 %/y")
	interestRate.Move(200, 90)
	bankwidget.AddChild(interestRate)
	bankwidget.interestRateL = interestRate

	currentInterestL := sws.NewLabelWidget(190, 25, "Current interest/y")
	currentInterestL.Move(10, 115)
	bankwidget.AddChild(currentInterestL)

	currentInterest := sws.NewLabelWidget(100, 25, "0 $")
	currentInterest.Move(200, 115)
	bankwidget.AddChild(currentInterest)
	bankwidget.currentInterest = currentInterest

	bankwidget.paybackbutton = sws.NewButtonWidget(150, 25, "Payback debt")
	bankwidget.paybackbutton.Move(10, 140)
	bankwidget.AddChild(bankwidget.paybackbutton)

	bankwidget.askloanbutton = sws.NewButtonWidget(150, 25, "Ask for loan")
	bankwidget.askloanbutton.Move(170, 140)
	bankwidget.AddChild(bankwidget.askloanbutton)

	GlobalLedger.AddSubscriber(bankwidget)

	return bankwidget
}
