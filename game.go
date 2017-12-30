package dctycoon

import (
	"math/rand"

	"github.com/nzin/dctycoon/accounting"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/dctycoon/timer"
	log "github.com/sirupsen/logrus"
)

// Actor -> Player or NPDatacenter
type Actor interface {
	GetInventory() *supplier.Inventory
	GetLedger() *accounting.Ledger
	//	Save()
	//	Load(gametimer)
}

//
// the global picture of the game, i.e.
// - actors (player and nonplayer)
// - marketplace items (DemandTemplate, ServerBundle)
// - load/save game functions
type Game struct {
	timer           *timer.GameTimer
	actors          []Actor
	demandtemplates []*supplier.DemandTemplate
	serverbundles   []*supplier.ServerBundle
}

//
// RegisterActor used to register user AND NPDatacenter inventory
func (self *Game) RegisterActor(actor Actor) {
	log.Debug("Game::RegisterActor(", actor, ")")
	// check if not already present
	for _, a := range self.actors {
		if a == actor {
			return
		}
	}
	self.actors = append(self.actors, actor)
}

//
// NewGame create a common place where we generate customer demand
// and assign them to the best provider (aka Datacenter)
func NewGame(timer *timer.GameTimer) *Game {
	log.Debug("NewGame(", timer, ")")

	g := &Game{
		timer:           timer,
		actors:          make([]Actor, 0, 0),
		demandtemplates: make([]*supplier.DemandTemplate, 0, 0),
		serverbundles:   make([]*supplier.ServerBundle, 0, 0),
	}

	timer.AddCron(-1, -1, -1, func() {
		g.GenerateDemandAndFee()
	})

	return g
}

//
// GenerateDemandAndFee generate randomly new demand (and check if a datacenter can handle it)
func (self *Game) GenerateDemandAndFee() {
	log.Debug("Game::GenerateDemandAndFee()")
	// check if a bundle has to be paid (and renewed)
	for _, sb := range self.serverbundles {
		//
		if sb.Date.Day() == self.timer.CurrentTime.Day() {
			// pay monthly fee

			// TBD
		}

		if sb.Date.Month() == self.timer.CurrentTime.Month() &&
			sb.Date.Day() == self.timer.CurrentTime.Day() {
			// if we don't renew, we drop it
			if rand.Float64() >= sb.Renewalrate {
				for _, c := range sb.Contracts {
					c.Item.Pool.Release(c.Item, c.Nbcores, c.Ramsize, c.Disksize)
				}
			}
		}
	}

	inventories := make([]*supplier.Inventory, 0, 0)
	for _, a := range self.actors {
		inventories = append(inventories, a.GetInventory())
	}
	// generate new demand
	for _, d := range self.demandtemplates {
		if d.Beginningdate.After(self.timer.CurrentTime) {
			if rand.Float64() >= (365.0-float64(d.Howoften))/365.0 {
				// here we will generate a new server demand (and see if it is fulfill)
				demand := d.InstanciateDemand()
				serverbundle := demand.FindOffer(inventories, self.timer.CurrentTime)
				if serverbundle != nil {
					self.serverbundles = append(self.serverbundles, serverbundle)
				}
			}
		}
	}
}
