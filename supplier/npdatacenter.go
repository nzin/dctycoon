package supplier

import (
	"encoding/json"

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
	inventory      *Inventory
	trend          *Trend
	timer          *timer.GameTimer
	initialcapital float64
	location       *LocationType
	buyoutprofile  map[string]BuyoutProfile
}

//
// NewNPDatacenter() create a new NonPlayerDatacenter that will compete with the player
func NewNPDatacenter(timer *timer.GameTimer, trend *Trend, initialcapital float64, locationid string, profilename string) *NPDatacenter {
	log.Debug("NewNPDatacenter(", timer, ",", trend, ",", initialcapital, ",", locationid, ",", profilename, ")")
	// a default value
	location := AvailableLocation["siliconvalley"]

	if l, ok := AvailableLocation[locationid]; ok {
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
		inventory:      NewInventory(timer),
		trend:          trend,
		timer:          timer,
		initialcapital: initialcapital,
		location:       location,
		buyoutprofile:  profile,
	}

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
	cartitems := []*CartItem{}
	for profilename, profile := range self.buyoutprofile {

		// if the offer is a VPS and we don't have yet VT enable hardware... we skip for next year
		if profile.Vps == true && self.trend.Vt.CurrentValue(self.timer.CurrentTime) == 0 {
			continue
		}

		conftype := GetServerConfTypeByName(profile.Servertype)
		var serverconf ServerConf
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

			nb := int32((self.initialcapital * profile.Buyperyear) / unitprice)
			// if we can afford, then we buy it
			if nb > 0 {
				item := &CartItem{
					Typeitem:   PRODUCT_SERVER,
					Serverconf: &serverconf,
					Unitprice:  unitprice,
					Nb:         nb,
				}
				cartitems = append(cartitems, item)

				// create the offer
				var pool ServerPool
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
				offer := &ServerOffer{
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
	self.inventory.Cart = cartitems
	self.inventory.BuyCart(self.timer.CurrentTime)
}
