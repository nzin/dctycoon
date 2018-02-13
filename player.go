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
	supplier.ActorAbstract
	reputation  *supplier.Reputation
	firewall    *firewall.Firewall
	companyname string
	maplevel    int32 // from 0 (3x4 map), to 2 (32x32 map)
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

	actorabstract := supplier.NewActorAbstract()
	p := &Player{
		ActorAbstract: *actorabstract,
		reputation:    nil,
		companyname:   "noname",
		maplevel:      0,
		firewall:      firewall.NewFirewall(),
	}

	return p
}

func (self *Player) Init(timer *timer.GameTimer, initialcapital float64, locationname, companyname string, maplevel int32) {
	log.Debug("Player::Init(", timer, ",", initialcapital, ",", locationname, ")")
	self.ActorAbstract.Init(timer, locationname, "you")

	self.reputation = supplier.NewReputation()
	self.companyname = companyname
	self.maplevel = maplevel
	self.firewall = firewall.NewFirewall()

	// add some equity
	self.GetLedger().AddMovement(accounting.LedgerMovement{
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

	self.ActorAbstract.Init(timer, locationname, "you")
	self.reputation = supplier.NewReputation()
	self.companyname = v["companyname"].(string)
	self.maplevel = int32(v["maplevel"].(float64))
	self.firewall = firewall.NewFirewall()

	self.GetLedger().Load(v["ledger"].(map[string]interface{}), self.GetLocation().Taxrate, self.GetLocation().Bankinterestrate)
	self.GetInventory().Load(v["inventory"].(map[string]interface{}))
	if offersinterface, ok := v["offers"]; ok {
		offers := offersinterface.([]interface{})
		for _, offer := range offers {
			self.LoadOffer(offer.(map[string]interface{}))
		}
	}
	self.reputation.Load(v["reputation"].(map[string]interface{}))
	self.firewall.Load(v["firewall"].(map[string]interface{}))
}

func (self *Player) Save() string {
	save := fmt.Sprintf(`"location": "%s",`, self.GetLocationName()) + "\n"
	save += fmt.Sprintf(`"inventory": %s,`, self.GetInventory().Save()) + "\n"
	save += `"offers":[`
	firstitem := true
	for _, offer := range self.GetOffers() {
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
	save += fmt.Sprintf(`"ledger": %s`, self.GetLedger().Save()) + "\n"
	return save
}
