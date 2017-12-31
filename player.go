package dctycoon

import (
	"fmt"

	"github.com/nzin/dctycoon/accounting"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/dctycoon/timer"
	log "github.com/sirupsen/logrus"
)

type Player struct {
	inventory    *supplier.Inventory
	ledger       *accounting.Ledger
	location     *supplier.LocationType
	locationname string
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
func NewPlayer() *Player {
	log.Debug("NewPlayer()")
	location := supplier.AvailableLocation["siliconvalley"]

	p := &Player{
		inventory:    nil,
		ledger:       nil,
		location:     location,
		locationname: "siliconvalley",
	}

	return p
}

func (self *Player) Init(timer *timer.GameTimer, initialcapital float64, locationname string) {
	log.Debug("Player::Init(", timer, ",", initialcapital, ",", locationname, ")")
	location := supplier.AvailableLocation["siliconvalley"]

	if l, ok := supplier.AvailableLocation[locationname]; ok {
		location = l
		self.locationname = locationname
	} else {
		log.Error("NewPlayer(): location " + locationname + " not found")
		self.locationname = "siliconvalley"
	}

	self.inventory = supplier.NewInventory(timer)
	self.ledger = accounting.NewLedger(timer, location.Taxrate, location.Bankinterestrate)
	self.location = location

	// add some equity
	self.ledger.AddMovement(accounting.LedgerMovement{
		Description: "initial capital",
		Amount:      initialcapital,
		AccountFrom: "4561",
		AccountTo:   "5121",
		Date:        timer.CurrentTime,
	})
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

	self.ledger.Load(v["ledger"].(map[string]interface{}), location.Taxrate, location.Bankinterestrate)
	self.inventory.Load(v["inventory"].(map[string]interface{}))
}

func (self *Player) Save() string {
	save := fmt.Sprintf(`"location": "%s",`, self.locationname) + "\n"
	save += fmt.Sprintf(`"inventory": %s,`, self.inventory.Save()) + "\n"
	save += fmt.Sprintf(`"ledger": %s`, self.ledger.Save()) + "\n"
	return save
}
