package accounting

import (
	"fmt"
	"github.com/google/btree"
	"time"
)

//
// payable or reception of money (depending of the sign of Amount
//
type LedgerMovement struct {
	Description string
	Amount      float64
	AccountFrom string // from credit
	AccountTo   string // to debit
	Date        time.Time
}

func (self *LedgerMovement) Less(b btree.Item) bool {
	return self.Date.Before(b.(*LedgerMovement).Date)
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
// 6132. Location immobiliere
// 6231. Publicité, Annonces et insertions
// 6233. Publicité, Foires et expositions
// 626. Frais postaux et de télécommunications (aka internet bill)
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
	LedgerChange(ledger *Ledger)
}

//
// Ledger
//
type Ledger struct {
	Movements   *btree.BTree
	subscribers []LedgerSubscriber
	accounts    map[int]AccountYearly
	taxrate     float64
}

var GlobalLedger *Ledger

func (self *Ledger) GetYearAccount(year int) AccountYearly {
	return self.accounts[year]
}

func (self *Ledger) AddMovement(ev LedgerMovement) {
	self.Movements.ReplaceOrInsert(&ev)
	self.accounts = self.RunLedger()
	for _, s := range self.subscribers {
		s.LedgerChange(self)
	}
}

func (self *Ledger) AddSubscriber(sub LedgerSubscriber) {
	self.subscribers = append(self.subscribers, sub)
	sub.LedgerChange(self)
}

func CreateLedger(taxrate float64) *Ledger {
	ledger := &Ledger{
		Movements:   btree.New(10),
		accounts:    make(map[int]AccountYearly),
		subscribers: make([]LedgerSubscriber, 0),
		taxrate:     taxrate,
	}
	return ledger
}

func (self *Ledger) Load(game map[string]interface{},taxrate float64) {
	self.Movements = btree.New(10)
	self.taxrate = taxrate

	for _, event := range game["movements"].([]interface{}) {
		e := event.(map[string]interface{})
		var year, month, day int
		fmt.Sscanf(e["date"].(string), "%d-%d-%d", &year, &month, &day)
		le := &LedgerMovement{
			Description: e["description"].(string),
			Amount:      e["amount"].(float64),
			AccountFrom: e["from"].(string),
			AccountTo:   e["to"].(string),
			Date:        time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC),
		}
		self.Movements.ReplaceOrInsert(le)
	}
	self.accounts = self.RunLedger()
	for _, s := range self.subscribers {
		s.LedgerChange(self)
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
func computeYearlyTaxes(accounts AccountYearly) float64 {
	var ebt float64
	ebt += accounts["70"]
	for k, v := range accounts {
		if k[0] == '6' {
			ebt -= v
		}
	}
	ebt -= accounts["28"] // amortissements
	ebt -= accounts["29"] // depreciations
	ebt -= accounts["66"] // interets sur la dette
	if ebt > 0 {
		return ebt * 0.15 // taxes are fixed... for now
	}
	return 0
}

func (self *Ledger) RunLedger() (accounts map[int]AccountYearly) {
	accounts = make(map[int]AccountYearly)
	currentyear := -1
	self.Movements.Ascend(func(i btree.Item) bool {
		ev := i.(*LedgerMovement)
		from := ev.AccountFrom[:2]
		to := ev.AccountTo[:2]

		if currentyear == -1 {
			currentyear = ev.Date.Year()
			accounts[currentyear] = make(AccountYearly)
		}
		if currentyear != ev.Date.Year() { // we must close the year and prepare the new year
			previousYear := currentyear
			currentyear = ev.Date.Year()
			accounts[previousYear]["44"] = computeYearlyTaxes(accounts[previousYear])
			accounts[previousYear]["51"] -= accounts[previousYear]["44"]
			accounts[currentyear] = make(AccountYearly)
			// copy from previous year, accounts 1 to 5 (except 44 => 0)
			for k, v := range accounts[previousYear] {
				if k[0] == '1' || k[0] == '2' || k[0] == '3' || k[0] == '4' || k[0] == '5' {
					accounts[currentyear][k] = v
				}
			}
			accounts[currentyear]["44"] = 0
		}

		accounts[currentyear][from] = accounts[currentyear][from] - ev.Amount
		accounts[currentyear][to] = accounts[currentyear][to] + ev.Amount
		//fmt.Println("from: ",from," to:",to, "currentyear: ",currentyear)

		return true
	})

	return accounts
}

func (self *Ledger) Save() string {
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
