package dctycoon

import (
	"fmt"

	"github.com/nzin/dctycoon/firewall"

	"github.com/nzin/dctycoon/accounting"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/dctycoon/timer"
	log "github.com/sirupsen/logrus"
)

type Player struct {
	inventory    *supplier.Inventory
	ledger       *accounting.Ledger
	location     *supplier.LocationType
	reputation   *supplier.Reputation
	firewall     *firewall.Firewall
	locationname string
	companyname  string
	maplevel     int32 // from 0 (3x4 map), to 2 (32x32 map)
	offers       []*supplier.ServerOffer
}

//
// GetReputationScore is part of the Actor interface
func (p *Player) GetReputationScore() float64 {
	return p.reputation.GetScore()
}

func (p *Player) GetFirewall() *firewall.Firewall {
	return p.firewall
}

func (p *Player) GetReputation() *supplier.Reputation {
	return p.reputation
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
// GetName is part of the Actor interface
func (p *Player) GetName() string {
	return "you"
}

func (p *Player) GetLocation() *supplier.LocationType {
	return p.location
}

func (p *Player) GetCompanyName() string {
	return p.companyname
}

func (p *Player) GetMapLevel() int32 {
	return p.maplevel
}

func (p *Player) SetMapLevel(maplevel int32) {
	p.maplevel = maplevel
}

//
// NewPlayer create a new player representation
func NewPlayer() *Player {
	log.Debug("NewPlayer()")
	location := supplier.AvailableLocation["siliconvalley"]

	p := &Player{
		inventory:    nil,
		ledger:       nil,
		location:     location,
		locationname: "siliconvalley",
		reputation:   nil,
		companyname:  "noname",
		maplevel:     0,
		firewall:     firewall.NewFirewall(),
		offers:       make([]*supplier.ServerOffer, 0),
	}

	return p
}

func (self *Player) Init(timer *timer.GameTimer, initialcapital float64, locationname, companyname string, maplevel int32) {
	log.Debug("Player::Init(", timer, ",", initialcapital, ",", locationname, ")")
	location := supplier.AvailableLocation["siliconvalley"]

	if l, ok := supplier.AvailableLocation[locationname]; ok {
		location = l
	} else {
		log.Error("NewPlayer(): location " + locationname + " not found")
		locationname = "siliconvalley"
	}

	self.locationname = locationname
	self.inventory = supplier.NewInventory(timer)
	self.ledger = accounting.NewLedger(timer, location.Taxrate, location.Bankinterestrate)
	self.location = location
	self.reputation = supplier.NewReputation()
	self.companyname = companyname
	self.maplevel = maplevel
	self.firewall = firewall.NewFirewall()

	// add some equity
	self.ledger.AddMovement(accounting.LedgerMovement{
		Description: "initial capital",
		Amount:      initialcapital,
		AccountFrom: "4561",
		AccountTo:   "5121",
		Date:        timer.CurrentTime,
	})
}
func (self *Player) loadOffer(offer map[string]interface{}) {
	log.Debug("Player::LoadOffer(", offer, ")")
	vps := offer["vps"].(bool)

	pool := self.inventory.GetDefaultPhysicalPool()
	if vps {
		pool = self.inventory.GetDefaultVpsPool()
	}

	nbcores := int32(offer["nbcores"].(float64))
	ramsize := int32(offer["ramsize"].(float64))
	disksize := int32(offer["disksize"].(float64))
	price, _ := offer["price"].(float64)

	o := &supplier.ServerOffer{
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

func (self *Player) LoadGame(timer *timer.GameTimer, v map[string]interface{}) {
	log.Debug("Player::LoadGame(", timer, ",", v, ")")
	locationname := v["location"].(string)
	location := supplier.AvailableLocation["siliconvalley"]

	if l, ok := supplier.AvailableLocation[locationname]; ok {
		location = l
	} else {
		log.Error("NewPlayer(): location " + locationname + " not found")
		locationname = "siliconvalley"
	}
	self.inventory = supplier.NewInventory(timer)
	self.ledger = accounting.NewLedger(timer, location.Taxrate, location.Bankinterestrate)
	self.location = location
	self.locationname = locationname
	self.reputation = supplier.NewReputation()
	self.companyname = v["companyname"].(string)
	self.maplevel = int32(v["maplevel"].(float64))
	self.firewall = firewall.NewFirewall()

	self.ledger.Load(v["ledger"].(map[string]interface{}), location.Taxrate, location.Bankinterestrate)
	self.inventory.Load(v["inventory"].(map[string]interface{}))
	if offersinterface, ok := v["offers"]; ok {
		offers := offersinterface.([]interface{})
		for _, offer := range offers {
			self.loadOffer(offer.(map[string]interface{}))
		}
	}
	self.reputation.Load(v["reputation"].(map[string]interface{}))
	self.firewall.Load(v["firewall"].(map[string]interface{}))
}

func (self *Player) Save() string {
	save := fmt.Sprintf(`"location": "%s",`, self.locationname) + "\n"
	save += fmt.Sprintf(`"inventory": %s,`, self.inventory.Save()) + "\n"
	save += `"offers":[`
	firstitem := true
	for _, offer := range self.offers {
		if firstitem == true {
			firstitem = false
		} else {
			save += ",\n"
		}
		save += offer.Save()
	}
	save += "],"
	save += fmt.Sprintf(`"companyname": "%s",`, self.companyname) + "\n"
	save += fmt.Sprintf(`"maplevel": %d,`, self.maplevel) + "\n"
	save += fmt.Sprintf(`"reputation": %s,`, self.reputation.Save()) + "\n"
	save += fmt.Sprintf(`"firewall": %s,`, self.firewall.Save()) + "\n"
	save += fmt.Sprintf(`"ledger": %s`, self.ledger.Save()) + "\n"
	return save
}

func (self *Player) AddOffer(offer *supplier.ServerOffer) {
	// check if not already present
	for _, o := range self.offers {
		if o == offer {
			return
		}
	}
	self.offers = append(self.offers, offer)
}

func (self *Player) RemoveOffer(offer *supplier.ServerOffer) {
	for i, o := range self.offers {
		if o == offer {
			self.offers = append(self.offers[:i], self.offers[i+1:]...)
			break
		}
	}
}

func (self *Player) GetOffers() []*supplier.ServerOffer {
	return self.offers
}
