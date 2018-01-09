package accounting

import (
	"fmt"
	"time"

	"github.com/google/btree"
	"github.com/nzin/dctycoon/timer"
	log "github.com/sirupsen/logrus"
)

//
// payable or reception of money (depending of the sign of Amount
//
type LedgerMovement struct {
	Id          int32 // filled by Ledger.AddMovement
	Description string
	Amount      float64
	AccountFrom string // from credit
	AccountTo   string // to debit
	Date        time.Time
}

func (self *LedgerMovement) Less(b btree.Item) bool {
	bmvt := b.(*LedgerMovement)
	if self.Date.Equal(bmvt.Date) {
		return self.Id < bmvt.Id
	} else {
		return self.Date.Before(bmvt.Date)
	}
}

//
// this interface is used by any stakeholder
// to invoice (supplier) or pay (customer) the DC
//
// notes for myself:
// https://fr.wikipedia.org/wiki/Plan_comptable_g%C3%A9n%C3%A9ral_(France)
// https://www.compta-facile.com/comptabilisation-des-apports-en-capital/
//
// 1012: Capital souscrit, appelé, non versé
// 4561: Associés – Comptes d’apport en société
//
// 1. Capitaux
// 161: dettes
//
// 2. Immo
// 2312: Immo, Terrains
// 2315: Immo, Installations techniques, matériel et outillage industriels
//
// 2815: Amortissements, Installations techniques, matériel et outillage industriels
// (2931: Dépreciations, Immobilisations corporelles en cours)
//
// 3. Comptes de stocks et en-cours
// 371. Marchandises
//
// 4. Compte tiers
// 4011. Fournisseurs - Achats de biens ou de prestations de services
//
// 444. État - Impôts sur les bénéfices
// 44571. TVA collectée
//
// 5. Comptes financiers
// 5121. Comptes en monnaie nationale (aka compte bancaire)
//
// 6. Charges
// 60612. Électricité
// 607. achat de marchandise
// 6132. Location immobiliere
// 6231. Publicité, Annonces et insertions
// 6233. Publicité, Foires et expositions
// (626.)-> 65 Frais postaux et de télécommunications (aka internet bill)
//
// 6311. Taxe sur les salaires
//
// 6411. Salaires, appointements
//
// 6611. Intérêts des emprunts et dettes
//
// 7.Produits:
// 7083. Locations diverses
//
// questions:
// - achat d'actions (de compagnies externes): 501
// - dettes: d'ou vient l'argent
// - vente de ses parts sous forme d'action?
// actif (immo+stock+) cf https://www.compta-facile.com/l-actif-du-bilan-comptable-en-detail/
// passif (capitaux+dettes) cf https://www.compta-facile.com/le-passif-du-bilan-comptable/
//
// https://www.compta-facile.com/comment-lire-comprendre-interpreter-bilan-comptable/
// EBITDA: https://www.compta-facile.com/ebitda-definition-calcul-utilite/
//
// BILAN
//
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
//
//
//
// PASSIF (année N)
//
// capitaux propre (45)
// resultat de l'exercice ??
// dettes (16)
// = passif
//
//
// ACTIF (année N)
//
// Terrain (2312)
// installation techniques (2315)
// autre immo (37)
// - amortissements (28)
// = actif
//
//
// - bank loan interest -> actif (compte) + passif (due) + un compte de depense pour les interests
// - buy -> actif (immo) + -argent (actif) + (immo -> un compte de depense (pour amortissement sur 4 ans))
// - DC location rent (compte de depense -> fournisseur)
// - internet fiber bill
// - electricity bill
// - salaries
// - maintenances?
//
// return nil if nothing to declare
//

type AccountYearly map[string]float64

//
// object interest to know the ledger status
// in particular the bank account
// essential the Dock
//
type LedgerSubscriber interface {
	LedgerChange()
}

//
// Ledger
//
type Ledger struct {
	autoinc      int32
	Movements    *btree.BTree
	subscribers  []LedgerSubscriber
	accounts     map[int]AccountYearly
	taxrate      float64
	loanrate     float64
	computeMonth func()
	timer        *timer.GameTimer
}

func (self *Ledger) GetYearAccount(year int) AccountYearly {
	if account, ok := self.accounts[year]; ok {
		return account
	}
	self.accounts = self.runLedger()

	return self.accounts[year]
}

func (self *Ledger) AddMovement(ev LedgerMovement) {
	log.Debug("Ledger::AddMovement(", ev, ")")
	ev.Id = self.autoinc
	self.Movements.ReplaceOrInsert(&ev)
	self.accounts = self.runLedger()
	for _, s := range self.subscribers {
		s.LedgerChange()
	}
	self.autoinc++
}

func (self *Ledger) AddSubscriber(sub LedgerSubscriber) {
	self.subscribers = append(self.subscribers, sub)
	sub.LedgerChange()
}

func (self *Ledger) RemoveSubscriber(subscriber LedgerSubscriber) {
	for i, s := range self.subscribers {
		if s == subscriber {
			self.subscribers = append(self.subscribers[:i], self.subscribers[i+1:]...)
			break
		}
	}
}

//
// 607 -> 2315 : product
// 404 -> 2315 : product
// et 5121 -> 4011 : money
// 2315 -> 2815 : amortization
func (self *Ledger) BuyProduct(desc string, t time.Time, price float64) {
	log.Debug("Ledger::BuyProduct(", desc, ",", t, ",", price, ")")
	product := &LedgerMovement{
		Id:          self.autoinc,
		Description: desc,
		Amount:      price,
		AccountFrom: "404",
		AccountTo:   "2315",
		Date:        t,
	}
	self.autoinc++
	self.Movements.ReplaceOrInsert(product)

	money := &LedgerMovement{
		Id:          self.autoinc,
		Description: desc,
		Amount:      price,
		AccountFrom: "5121",
		AccountTo:   "4011",
		Date:        t,
	}
	self.autoinc++
	self.Movements.ReplaceOrInsert(money)

	YN1 := time.Date(t.Year()+1, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	YN2 := time.Date(t.Year()+2, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	YN3 := time.Date(t.Year()+3, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	YN4 := time.Date(t.Year()+3, t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	tillY1 := YN1.Sub(t).Hours() / 24
	Y1toY2 := YN2.Sub(YN1).Hours() / 24
	Y2toY3 := YN3.Sub(YN2).Hours() / 24
	Y3toY4 := YN4.Sub(YN3).Hours() / 24

	self.Movements.ReplaceOrInsert(&LedgerMovement{
		Id:          self.autoinc,
		Description: "ammo " + desc,
		Amount:      price * tillY1 / (tillY1 + Y1toY2 + Y2toY3 + Y3toY4),
		AccountFrom: "2315",
		AccountTo:   "2815",
		Date:        t,
	})
	self.autoinc++
	self.Movements.ReplaceOrInsert(&LedgerMovement{
		Id:          self.autoinc,
		Description: "ammo " + desc,
		Amount:      price * Y1toY2 / (tillY1 + Y1toY2 + Y2toY3 + Y3toY4),
		AccountFrom: "2315",
		AccountTo:   "2815",
		Date:        YN1,
	})
	self.autoinc++
	self.Movements.ReplaceOrInsert(&LedgerMovement{
		Id:          self.autoinc,
		Description: "ammo " + desc,
		Amount:      price * Y2toY3 / (tillY1 + Y1toY2 + Y2toY3 + Y3toY4),
		AccountFrom: "2315",
		AccountTo:   "2815",
		Date:        YN2,
	})
	self.autoinc++
	self.Movements.ReplaceOrInsert(&LedgerMovement{
		Id:          self.autoinc,
		Description: "ammo " + desc,
		Amount:      price * Y3toY4 / (tillY1 + Y1toY2 + Y2toY3 + Y3toY4),
		AccountFrom: "2315",
		AccountTo:   "2815",
		Date:        YN3,
	})
	self.autoinc++

	// compute the ledger

	self.accounts = self.runLedger()
	for _, s := range self.subscribers {
		s.LedgerChange()
	}
}

func (self *Ledger) PayServerRenting(price float64, t time.Time, desc string) {
	self.AddMovement(LedgerMovement{
		Description: desc,
		Amount:      price,
		Date:        t,
		AccountFrom: "7083",
		AccountTo:   "5121",
	})
}

//
// 161 (capital/debt) -> 5121 (current bank account)
// every year (every month?) interest rate: 5121 -> 46 (fictious bank account for interest)
//
func (self *Ledger) AskLoan(desc string, t time.Time, amount float64) {
	log.Debug("Ledger::AskLoan(", desc, ",", t, ",", amount, ")")
	loan := &LedgerMovement{
		Id:          self.autoinc,
		Description: desc,
		Amount:      amount,
		AccountFrom: "161",
		AccountTo:   "5121",
		Date:        t,
	}
	self.autoinc++
	self.Movements.ReplaceOrInsert(loan)

	// compute the ledger

	self.accounts = self.runLedger()
	for _, s := range self.subscribers {
		s.LedgerChange()
	}
}

//
// 5121 (current bank account) -> 161 (capital/debt)
//
func (self *Ledger) RefundLoan(desc string, t time.Time, amount float64) {
	log.Debug("Ledger::RefundLoan(", desc, ",", t, ",", amount, ")")
	loan := &LedgerMovement{
		Id:          self.autoinc,
		Description: desc,
		Amount:      amount,
		AccountFrom: "5121",
		AccountTo:   "161",
		Date:        t,
	}
	self.autoinc++
	self.Movements.ReplaceOrInsert(loan)

	// compute the ledger

	self.accounts = self.runLedger()
	for _, s := range self.subscribers {
		s.LedgerChange()
	}
}

// NewLedger create a global accounting ledger:
// - taxrate are in percent per year (for example 0.15)
// - loanrate are in percent per year (for example 0.03)
func NewLedger(timer *timer.GameTimer, taxrate, loanrate float64) *Ledger {
	log.Debug("NewLedger(", timer, ",", taxrate, ",", loanrate, ")")
	ledger := &Ledger{
		Movements:   btree.New(10),
		accounts:    make(map[int]AccountYearly),
		subscribers: make([]LedgerSubscriber, 0),
		taxrate:     taxrate,
		loanrate:    loanrate,
		timer:       timer,
	}
	// compute fiscal year
	ledger.computeMonth = func() {
		l := ledger
		l.accounts = l.runLedger()
		for _, s := range l.subscribers {
			s.LedgerChange()
		}
	}
	timer.AddCron(1, -1, -1, ledger.computeMonth)

	return ledger
}

func (self *Ledger) Load(game map[string]interface{}, taxrate, loanrate float64) {
	log.Debug("Ledger::Load(", game, ",", taxrate, ",", loanrate, ")")
	self.Movements = btree.New(10)
	self.taxrate = taxrate
	self.loanrate = loanrate
	self.autoinc = 0

	for _, event := range game["movements"].([]interface{}) {
		e := event.(map[string]interface{})
		var year, month, day int
		fmt.Sscanf(e["date"].(string), "%d-%d-%d", &year, &month, &day)
		le := &LedgerMovement{
			Id:          self.autoinc,
			Description: e["description"].(string),
			Amount:      e["amount"].(float64),
			AccountFrom: e["from"].(string),
			AccountTo:   e["to"].(string),
			Date:        time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC),
		}
		self.autoinc++
		self.Movements.ReplaceOrInsert(le)
	}
	self.accounts = self.runLedger()
	for _, s := range self.subscribers {
		s.LedgerChange()
	}

}

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
func computeYearlyTaxes(accounts AccountYearly, taxrate float64) (profitlost, taxes float64) {
	log.Debug("computeYearlyTaxes(", accounts, ",", taxrate, ")")
	var ebt float64
	ebt -= accounts["70"]
	for k, v := range accounts {
		if k[0] == '6' {
			ebt += v
		}
	}
	ebt -= accounts["28"] // amortissements
	ebt -= accounts["29"] // depreciations
	ebt -= accounts["66"] // interets sur la dette
	if ebt > 0 {
		taxes = ebt * taxrate
	}
	return ebt - taxes, taxes
}

func (self *Ledger) runLedgerMonth(accounts map[int]AccountYearly, currentYear *int, currentMonth *int) {
	// compute loan for the month
	loanAmount := -accounts[*currentYear]["16"] * self.loanrate / 12
	accounts[*currentYear]["51"] -= loanAmount
	accounts[*currentYear]["66"] += loanAmount

	*currentMonth++
	if *currentMonth > 12 {
		*currentMonth = 1

		previousYear := *currentYear
		*currentYear++

		//fmt.Println("compute taxes for ", currentyear)
		profitlost, taxes := computeYearlyTaxes(accounts[previousYear], self.taxrate)
		accounts[previousYear]["44"] = taxes

		// 51: current balance, 44: taxes
		accounts[previousYear]["51"] -= accounts[previousYear]["44"]
		accounts[*currentYear] = make(AccountYearly)
		// copy from previous year, accounts 1 to 5 (except 44 => 0)
		for k, v := range accounts[previousYear] {
			if k[0] == '1' || k[0] == '2' || k[0] == '3' || k[0] == '4' || k[0] == '5' {
				accounts[*currentYear][k] = v
			}
		}
		accounts[*currentYear]["44"] = 0
		accounts[*currentYear]["45"] -= profitlost
		accounts[*currentYear]["46"] += accounts[previousYear]["66"]
		//accounts[currentyear]["23"] -= accounts[previousYear]["28"]
		accounts[*currentYear]["28"] = 0
	}
}

//
// runLedger() computes all movement and returns the result.
// it does not store the result in the Ledger object
// Usually you don't call this function but AddMovement()
func (self *Ledger) runLedger() (accounts map[int]AccountYearly) {
	log.Debug("runLedger()")
	accounts = make(map[int]AccountYearly)
	currentyear := -1
	currentMonth := -1
	self.Movements.Ascend(func(i btree.Item) bool {
		ev := i.(*LedgerMovement)
		from := ev.AccountFrom[:2]
		to := ev.AccountTo[:2]

		if currentyear == -1 {
			currentyear = ev.Date.Year()
			currentMonth = int(ev.Date.Month())
			accounts[currentyear] = make(AccountYearly)
		}
		for !(currentyear == ev.Date.Year() && currentMonth == int(ev.Date.Month())) { // we must close the month (and eventually the year)
			self.runLedgerMonth(accounts, &currentyear, &currentMonth)
		}

		accounts[currentyear][from] = accounts[currentyear][from] - ev.Amount
		accounts[currentyear][to] = accounts[currentyear][to] + ev.Amount
		log.Debug("from: ", from, " to:", to, "currentyear: ", currentyear, currentMonth, ",desc: ", ev.Description, accounts[currentyear][from], accounts[currentyear][to])

		return true
	})

	// if we don't have any movements...
	if currentyear == -1 || currentyear > self.timer.CurrentTime.Year() || (currentyear == self.timer.CurrentTime.Year() && currentMonth > int(self.timer.CurrentTime.Month())) {
		return accounts
	}
	//
	// in case of we have no movement from "current year" until today
	//
	for !(currentyear == self.timer.CurrentTime.Year() && currentMonth == int(self.timer.CurrentTime.Month())) { // we must close the month until self.timer.CurrentTime
		self.runLedgerMonth(accounts, &currentyear, &currentMonth)
	}

	return accounts
}

func (self *Ledger) Save() string {
	log.Debug("Ledger::Save()")
	str := "{\n"
	str += `"movements": [`
	self.Movements.Ascend(func(i btree.Item) bool {
		ev := i.(*LedgerMovement)
		if i == self.Movements.Min() {
			str += fmt.Sprintf(`  {"description":"%s", "amount": %v , "date":"%d-%d-%d", "from": "%s", "to":"%s" }`, ev.Description, ev.Amount, ev.Date.Year(), ev.Date.Month(), ev.Date.Day(), ev.AccountFrom, ev.AccountTo)
		} else {
			str += fmt.Sprintf(`,{"description":"%s", "amount": %v , "date":"%d-%d-%d", "from":"%s", "to":"%s" }`, ev.Description, ev.Amount, ev.Date.Year(), ev.Date.Month(), ev.Date.Day(), ev.AccountFrom, ev.AccountTo)
		}
		return true
	})
	str += "\n]}"
	return str
}

//
// Bank stakeholder
//
/*
type Bank struct {
	location   LocationType
	moneyOwned float64
}

func (self *Bank) Movement(t time.Time) *LedgerMovement {
	if self.moneyOwned==0 {
		return nil
	}
	return &LedgerMovement{
		Description: "mortage",
		Amount:      self.location.bankinterestrate*self.moneyOwned,
		Date:        t,
	}
}
*/
