package dctycoon

import (
	"encoding/json"
	"fmt"

	"github.com/nzin/dctycoon/accounting"
	"github.com/nzin/dctycoon/supplier"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/timer"
	log "github.com/sirupsen/logrus"
)

//
//  buyoutprofile: {
//	  "R100vps": {
//     "servertype": "R100"
//      "vps": true,
//	    "buyperyear": 0.1,     // 10% of remaining capital
//      "margin": 2.0,         // service price associated
//      "configuration": "low" // [low, medium, high] => low = 1 cpu, 1 ram, 1 disk,... high=> max cpu, max ram, max disk ..., medium: high/2
//    },
//    "R200": {
//      ...
//    }
//  }
//
type BuyoutProfile struct {
	Servertype    string
	Vps           bool
	Buyperyear    float64
	Margin        float64
	Configuration string
}

type NPDatacenter struct {
	inventory     *supplier.Inventory
	ledger        *accounting.Ledger
	trend         *supplier.Trend
	timer         *timer.GameTimer
	location      *supplier.LocationType
	buyoutprofile map[string]BuyoutProfile
}

//
// GetInventory is part of the Actor interface
func (self *NPDatacenter) GetInventory() *supplier.Inventory {
	return self.inventory
}

//
// GetLedger is part of the Actor interface
func (self *NPDatacenter) GetLedger() *accounting.Ledger {
	return self.ledger
}

//
// NewNPDatacenter() create a new NonPlayerDatacenter that will compete with the player
func NewNPDatacenter(timer *timer.GameTimer, trend *supplier.Trend, initialcapital float64, locationid string, profilename string) *NPDatacenter {
	log.Debug("NewNPDatacenter(", timer, ",", trend, ",", initialcapital, ",", locationid, ",", profilename, ")")
	// a default value
	location := supplier.AvailableLocation["siliconvalley"]

	if l, ok := supplier.AvailableLocation[locationid]; ok {
		location = l
	} else {
		log.Error("NewNPDatacenter(): location " + locationid + " not found")
	}

	// load buyout profile
	data, err := global.Asset("assets/npdatacenter/" + profilename)
	if err != nil {
		log.Error("NewNPDatacenter(): asset assets/npdatacenter/" + profilename + " not found")
		return nil
	}
	profile := make(map[string]BuyoutProfile)
	err = json.Unmarshal(data, &profile)
	if err != nil {
		log.Error("NewNPDatacenter(): asset assets/npdatacenter/" + profilename + " not json compatible")
		return nil
	}

	datacenter := &NPDatacenter{
		inventory:     supplier.NewInventory(timer),
		ledger:        accounting.NewLedger(timer, location.Taxrate, location.Bankinterestrate),
		trend:         trend,
		timer:         timer,
		location:      location,
		buyoutprofile: profile,
	}

	// add some equity
	datacenter.ledger.AddMovement(accounting.LedgerMovement{
		Description: "initial capital",
		Amount:      initialcapital,
		AccountFrom: "4561",
		AccountTo:   "5121",
		Date:        timer.CurrentTime,
	})

	timer.AddCron(1, 1, -1, func() {
		datacenter.NewYearOperations()
	})
	return datacenter
}

//
// NewYearOperations will trigger every year different actions, in particular
// - buy goods
// - create/refresh offers
func (self *NPDatacenter) NewYearOperations() {
	log.Debug("NPDatacenter::NewYearOperations()")
	// let's begin by removing all offers
	for _, o := range self.inventory.GetOffers() {
		self.inventory.RemoveOffer(o)
	}

	// loop over the buyoutprofile and buy some hardware + create the corresponding offer
	for profilename, profile := range self.buyoutprofile {

		// if the offer is a VPS and we don't have yet VT enable hardware... we skip for next year
		if profile.Vps == true && self.trend.Vt.CurrentValue(self.timer.CurrentTime) == 0 {
			continue
		}

		conftype := supplier.GetServerConfTypeByName(profile.Servertype)
		var serverconf supplier.ServerConf
		if conftype != nil {
			switch profile.Configuration {
			case "low":
				serverconf.NbProcessors = 1
				serverconf.NbCore = 1
				if self.trend.Vt.CurrentValue(self.timer.CurrentTime) == 1 && profile.Vps == false {
					serverconf.VtSupport = true
				} else {
					serverconf.VtSupport = false
				}
				serverconf.NbDisks = conftype.NbDisks[0]
				serverconf.NbSlotRam = conftype.NbSlotRam[0]
				serverconf.DiskSize = self.trend.Disksize.CurrentValue(self.timer.CurrentTime) / 4 // 3 options: Trend.Disksize: 1,1/2,1/4
				serverconf.RamSize = self.trend.Ramsize.CurrentValue(self.timer.CurrentTime) / 4   // 4 options: Trend.Ramsize: 1,1/2,1/4,1/8
				serverconf.ConfType = conftype
				break
			case "medium":
				serverconf.NbProcessors = (conftype.NbProcessors[0] + conftype.NbProcessors[1]) / 2
				if self.trend.Corepercpu.CurrentValue(self.timer.CurrentTime) > 1 {
					serverconf.NbCore = self.trend.Corepercpu.CurrentValue(self.timer.CurrentTime) / 2
				} else {
					serverconf.NbCore = 1
				}
				if self.trend.Vt.CurrentValue(self.timer.CurrentTime) == 1 && profile.Vps == false {
					serverconf.VtSupport = true
				} else {
					serverconf.VtSupport = false
				}
				serverconf.NbDisks = (conftype.NbDisks[0] + conftype.NbDisks[1]) / 2
				serverconf.NbSlotRam = (conftype.NbSlotRam[0] + conftype.NbSlotRam[1]) / 2
				serverconf.DiskSize = self.trend.Disksize.CurrentValue(self.timer.CurrentTime) / 2 // 3 options: Trend.Disksize: 1,1/2,1/4
				serverconf.RamSize = self.trend.Ramsize.CurrentValue(self.timer.CurrentTime) / 2   // 4 options: Trend.Ramsize: 1,1/2,1/4,1/8
				serverconf.ConfType = conftype
				break
			case "high":
				serverconf.NbProcessors = conftype.NbProcessors[1]
				serverconf.NbCore = self.trend.Corepercpu.CurrentValue(self.timer.CurrentTime)
				if self.trend.Vt.CurrentValue(self.timer.CurrentTime) == 1 && profile.Vps == false {
					serverconf.VtSupport = true
				} else {
					serverconf.VtSupport = false
				}
				serverconf.NbDisks = conftype.NbDisks[1]
				serverconf.NbSlotRam = conftype.NbSlotRam[1]
				serverconf.DiskSize = self.trend.Disksize.CurrentValue(self.timer.CurrentTime) // 3 options: Trend.Disksize: 1,1/2,1/4
				serverconf.RamSize = self.trend.Ramsize.CurrentValue(self.timer.CurrentTime)   // 4 options: Trend.Ramsize: 1,1/2,1/4,1/8
				serverconf.ConfType = conftype
				break
			default: // profile configuration not found
				log.Error("NPDatacenter::NewYearOperations(): profile " + profilename + " has a strange configuration: " + profile.Configuration)
				continue
			}
			unitprice := serverconf.Price(self.trend, self.timer.CurrentTime)
			log.Info("NPDatacenter::NewYearOperations(): profilename", profilename, "unitprice:", unitprice, "profilemargin:", profile.Margin)

			currentCash := self.ledger.GetYearAccount(self.timer.CurrentTime.Year())["51"]

			nb := int32((currentCash * profile.Buyperyear) / unitprice)
			// if we can afford, then we buy it
			if nb > 0 {
				item := &supplier.CartItem{
					Typeitem:   supplier.PRODUCT_SERVER,
					Serverconf: &serverconf,
					Unitprice:  unitprice,
					Nb:         nb,
				}
				self.inventory.Cart = append(self.inventory.Cart, item)
				self.ledger.BuyProduct(fmt.Sprintf("buying %s", profilename), self.timer.CurrentTime, item.Unitprice*float64(nb))
				self.inventory.BuyCart(self.timer.CurrentTime)
				self.inventory.Cart = make([]*supplier.CartItem, 0)

				// create the offer
				var pool supplier.ServerPool
				// check the first available pool
				for _, p := range self.inventory.GetPools() {
					if p.IsVps() == profile.Vps {
						pool = p
						break
					}
				}
				if pool == nil {
					log.Error("NPDatacenter::NewYearOperations(): We didn't find a correct (default) pool!")
					continue
				}
				offer := &supplier.ServerOffer{
					Active:    true,
					Name:      profilename,
					Inventory: self.inventory,
					Pool:      pool,
					Vps:       profile.Vps,
					Nbcores:   serverconf.NbCore * serverconf.NbProcessors,
					Ramsize:   serverconf.NbSlotRam * serverconf.RamSize,
					Disksize:  serverconf.NbDisks * serverconf.DiskSize,
					Vt:        serverconf.VtSupport,
					Price:     unitprice * profile.Margin,
				}
				self.inventory.AddOffer(offer)
			}
		}
	}
}