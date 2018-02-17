package accounting

import (
	"fmt"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/timer"
	"github.com/nzin/sws"
	//"github.com/veandco/go-sdl2/sdl"
	"strconv"
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
	askInput        *sws.InputWidget
	paybackInput    *sws.InputWidget
	timer           *timer.GameTimer
	ledger          *Ledger
}

func (self *BankWidget) LedgerChange() {
	yearaccount := self.ledger.GetYearAccount(self.timer.CurrentTime.Year())
	self.interestRate = self.ledger.loanrate
	self.interestRateL.SetText(fmt.Sprintf("%.2f %%/y", self.ledger.loanrate*100))
	self.accountPosition.SetText(fmt.Sprintf("%.2f $", yearaccount["51"]))
	self.accountDebt.SetText(fmt.Sprintf("%.2f $", -yearaccount["16"]))
	self.currentInterest.SetText(fmt.Sprintf("%.2f $/y", -yearaccount["16"]*self.interestRate))
	self.PostUpdate()
}

func NewBankWidget(root *sws.RootWidget) *BankWidget {
	bankwidget := &BankWidget{
		CoreWidget: *sws.NewCoreWidget(420, 280),
		timer:      nil,
		ledger:     nil,
	}

	title := sws.NewLabelWidget(190, 25, "Your Bank account")
	title.Move(80, 10)
	bankwidget.AddChild(title)

	accountPositionTitle := sws.NewLabelWidget(190, 25, "Your current position")
	accountPositionTitle.Move(80, 40)
	bankwidget.AddChild(accountPositionTitle)

	accountPosition := sws.NewLabelWidget(100, 25, "0 $")
	accountPosition.Move(270, 40)
	bankwidget.AddChild(accountPosition)
	bankwidget.accountPosition = accountPosition

	accountDebtTitle := sws.NewLabelWidget(190, 25, "Your current debt")
	accountDebtTitle.Move(80, 65)
	bankwidget.AddChild(accountDebtTitle)

	accountDebt := sws.NewLabelWidget(100, 25, "0 $")
	accountDebt.Move(270, 65)
	bankwidget.AddChild(accountDebt)
	bankwidget.accountDebt = accountDebt

	interestRate := sws.NewLabelWidget(100, 25, "0 %/y")
	interestRate.Move(270, 90)
	bankwidget.AddChild(interestRate)
	bankwidget.interestRateL = interestRate

	currentInterestL := sws.NewLabelWidget(190, 25, "Current interest/y")
	currentInterestL.Move(80, 115)
	bankwidget.AddChild(currentInterestL)

	currentInterest := sws.NewLabelWidget(100, 25, "0 $")
	currentInterest.Move(270, 115)
	bankwidget.AddChild(currentInterest)
	bankwidget.currentInterest = currentInterest

	hr1 := sws.NewHr(460)
	hr1.Move(10, 145)
	bankwidget.AddChild(hr1)

	askLabel := sws.NewLabelWidget(100, 25, "Ask for a loan")
	askLabel.Move(80, 150)
	bankwidget.AddChild(askLabel)

	askAmountLabel := sws.NewLabelWidget(100, 25, "Amount")
	askAmountLabel.Move(80, 180)
	bankwidget.AddChild(askAmountLabel)

	askInput := sws.NewInputWidget(100, 25, "")
	askInput.Move(180, 180)
	bankwidget.askInput = askInput
	bankwidget.AddChild(askInput)

	bankwidget.askloanbutton = sws.NewButtonWidget(100, 25, "Ask")
	bankwidget.askloanbutton.Move(290, 180)
	bankwidget.AddChild(bankwidget.askloanbutton)
	bankwidget.askloanbutton.SetClicked(func() {
		value := bankwidget.askInput.GetText()
		if asked, err := strconv.ParseFloat(value, 64); err != nil {
			iconsurface, _ := global.LoadImageAsset("assets/ui/paper-bill.png")
			sws.ShowModalErrorSurfaceicon(root, "Amount error", iconsurface, "The amount doesn't seems to be a number", nil)
		} else {
			maxDebtPossibleDebt := bankwidget.ledger.GetMaximumPossibleDdebt(bankwidget.timer.CurrentTime.Year())
			currentDebt := bankwidget.ledger.GetCurrentDebt(bankwidget.timer.CurrentTime.Year())
			if asked+currentDebt > maxDebtPossibleDebt {
				iconsurface, _ := global.LoadImageAsset("assets/ui/paper-bill.png")
				sws.ShowModalErrorSurfaceicon(root, "Amount inquiry error", iconsurface, "Seriously? You want to borrow that amount?\n Kid, prove you can run a big business and we will reconsider your demand (check your EBITDA).", nil)
			} else {
				bankwidget.ledger.AskLoan("bank loan", bankwidget.timer.CurrentTime, asked)
			}
			bankwidget.askInput.SetText("")
		}
	})

	hr2 := sws.NewHr(390)
	hr2.Move(80, 215)
	bankwidget.AddChild(hr2)

	paybackLabel := sws.NewLabelWidget(100, 25, "Refund loan")
	paybackLabel.Move(80, 220)
	bankwidget.AddChild(paybackLabel)

	paybackAmountLabel := sws.NewLabelWidget(100, 25, "Amount")
	paybackAmountLabel.Move(80, 250)
	bankwidget.AddChild(paybackAmountLabel)

	paybackInput := sws.NewInputWidget(100, 25, "")
	paybackInput.Move(180, 250)
	bankwidget.paybackInput = paybackInput
	bankwidget.AddChild(paybackInput)

	bankwidget.paybackbutton = sws.NewButtonWidget(100, 25, "Payback")
	bankwidget.paybackbutton.Move(290, 250)
	bankwidget.AddChild(bankwidget.paybackbutton)
	bankwidget.paybackbutton.SetClicked(func() {
		value := bankwidget.paybackInput.GetText()
		if refund, err := strconv.ParseFloat(value, 64); err != nil {
			iconsurface, _ := global.LoadImageAsset("assets/ui/paper-bill.png")
			sws.ShowModalErrorSurfaceicon(root, "Amount error", iconsurface, "The amount doesn't seems to be a number", nil)
		} else {
			yearaccountN := bankwidget.ledger.GetYearAccount(bankwidget.timer.CurrentTime.Year())
			currentDebt := -yearaccountN["16"]
			currentMoney := yearaccountN["51"]
			if refund > currentDebt {
				refund = currentDebt
			}
			if refund > currentMoney {
				iconsurface, _ := global.LoadImageAsset("assets/ui/paper-bill.png")
				sws.ShowModalErrorSurfaceicon(root, "Cashflow problem", iconsurface, "I don't think you can afford to refund so much money, keep working on your business...", nil)
			} else {
				bankwidget.ledger.RefundLoan("payback bank debt", bankwidget.timer.CurrentTime, refund)
			}
			bankwidget.paybackInput.SetText("")
		}
	})

	bankicon := sws.NewLabelWidget(64, 64, "")
	if icon, err := global.LoadImageAsset("assets/ui/icon-bank.big.png"); err == nil {
		bankicon.SetImageSurface(icon)
	}
	bankicon.Move(4, 40)
	bankwidget.AddChild(bankicon)

	loanicon := sws.NewLabelWidget(64, 64, "")
	if icon, err := global.LoadImageAsset("assets/ui/icon-loan.big.png"); err == nil {
		loanicon.SetImageSurface(icon)
	}
	loanicon.Move(4, 183)
	bankwidget.AddChild(loanicon)

	return bankwidget
}

func (self *BankWidget) SetGame(timer *timer.GameTimer, ledger *Ledger) {
	self.timer = timer
	if self.ledger != nil {
		self.ledger.RemoveSubscriber(self)
	}
	self.ledger = ledger
	ledger.AddSubscriber(self)
}
