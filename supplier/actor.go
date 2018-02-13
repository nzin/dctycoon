package supplier

import (
	"github.com/nzin/dctycoon/accounting"
	"github.com/nzin/dctycoon/timer"
	log "github.com/sirupsen/logrus"
)

// Actor -> Player or NPDatacenter
type Actor interface {
	GetInventory() *Inventory
	GetLedger() *accounting.Ledger
	GetName() string
	GetLocation() *LocationType
	GetReputationScore() float64
	GetOffers() []*ServerOffer
}

type ActorAbstract struct {
	inventory    *Inventory
	ledger       *accounting.Ledger
	location     *LocationType
	locationname string
	name         string
	offers       []*ServerOffer
}

//
// GetInventory is part of the Actor interface
func (actor *ActorAbstract) GetInventory() *Inventory {
	return actor.inventory
}

//
// GetLedger is part of the Actor interface
func (actor *ActorAbstract) GetLedger() *accounting.Ledger {
	return actor.ledger
}

func (actor *ActorAbstract) GetLocation() *LocationType {
	return actor.location
}

func (actor *ActorAbstract) GetLocationName() string {
	return actor.locationname
}

func (self *ActorAbstract) GetName() string {
	return self.name
}

//
// GetReputationScore is part of the Actor interface
func (self *ActorAbstract) GetReputationScore() float64 {
	return 0
}

func (actor *ActorAbstract) Init(timer *timer.GameTimer, locationname, name string) {
	location := AvailableLocation["siliconvalley"]

	if l, ok := AvailableLocation[locationname]; ok {
		location = l
	} else {
		log.Error("ActorAbstract::Init(): location " + locationname + " not found")
		locationname = "siliconvalley"
	}

	actor.locationname = locationname
	actor.inventory = NewInventory(timer)
	actor.ledger = accounting.NewLedger(timer, location.Taxrate, location.Bankinterestrate)
	actor.location = location
	actor.offers = make([]*ServerOffer, 0)
	actor.name = name
}

func (self *ActorAbstract) LoadOffer(offer map[string]interface{}) {
	log.Debug("ActorAbstract::LoadOffer(", offer, ")")
	vps := offer["vps"].(bool)

	pool := self.inventory.GetDefaultPhysicalPool()
	if vps {
		pool = self.inventory.GetDefaultVpsPool()
	}

	nbcores := int32(offer["nbcores"].(float64))
	ramsize := int32(offer["ramsize"].(float64))
	disksize := int32(offer["disksize"].(float64))
	price, _ := offer["price"].(float64)

	o := &ServerOffer{
		Active:   offer["active"].(bool),
		Name:     offer["name"].(string),
		Actor:    self,
		Pool:     pool,
		Vps:      vps,
		Nbcores:  nbcores,
		Ramsize:  ramsize,
		Disksize: disksize,
		Vt:       offer["vt"].(bool),
		Price:    price,
	}
	self.AddOffer(o)
}

func (self *ActorAbstract) AddOffer(offer *ServerOffer) {
	// check if not already present
	for _, o := range self.offers {
		if o == offer {
			return
		}
	}
	self.offers = append(self.offers, offer)
}

func (self *ActorAbstract) RemoveOffer(offer *ServerOffer) {
	for i, o := range self.offers {
		if o == offer {
			self.offers = append(self.offers[:i], self.offers[i+1:]...)
			break
		}
	}
}

func (self *ActorAbstract) GetOffers() []*ServerOffer {
	return self.offers
}

func NewActorAbstract() *ActorAbstract {
	location := AvailableLocation["siliconvalley"]
	actor := &ActorAbstract{
		inventory:    nil,
		ledger:       nil,
		location:     location,
		locationname: "siliconvalley",
		offers:       make([]*ServerOffer, 0),
	}
	return actor
}
