package dctycoon

import (
	"github.com/nzin/dctycoon/accounting"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/dctycoon/timer"
	log "github.com/sirupsen/logrus"
)

type Player struct {
	inventory *supplier.Inventory
	ledger    *accounting.Ledger
	location  *supplier.LocationType
}

//
// GetInventory is part of the Actor interface
func (p *Player) GetInventory() *supplier.Inventory {
	return p.inventory
}

//
// GetLedger is part of the Actor interface
func (p *Player) GetLedger() *accounting.Ledger {
	return p.ledger
}

//
// NewPlayer create a new player representation
func NewPlayer(timer *timer.GameTimer, initialcapital float64, locationid string) *Player {
	log.Debug("NewPlayer(", timer, ",", initialcapital, ",", locationid, ")")
	location := supplier.AvailableLocation["siliconvalley"]

	if l, ok := supplier.AvailableLocation[locationid]; ok {
		location = l
	} else {
		log.Error("NewPlayer(): location " + locationid + " not found")
	}

	p := &Player{
		inventory: supplier.NewInventory(timer),
		ledger:    accounting.NewLedger(timer, location.Taxrate, location.Bankinterestrate),
		location:  location,
	}

	// add some equity
	p.ledger.AddMovement(accounting.LedgerMovement{
		Description: "initial capital",
		Amount:      initialcapital,
		AccountFrom: "4561",
		AccountTo:   "5121",
		Date:        timer.CurrentTime,
	})

	return p
}
