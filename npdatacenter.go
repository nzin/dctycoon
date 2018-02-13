package dctycoon

import (
	"encoding/json"
	"fmt"
	"math/rand"

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
	supplier.ActorAbstract
	trend           *supplier.Trend
	timer           *timer.GameTimer
	buyoutprofile   map[string]BuyoutProfile
	profilename     string
	cronevent       *timer.GameCronEvent
	picture         string
	reputationscore float64
}

//
// GetReputationScore is part of the Actor interface
func (self *NPDatacenter) GetReputationScore() float64 {
	return self.reputationscore
}

func (self *NPDatacenter) GetPicture() string {
	return self.picture
}

//
// NewNPDatacenter() create a new NonPlayerDatacenter that will compete with the player
func NewNPDatacenter() *NPDatacenter {
	log.Debug("NewNPDatacenter()")
	actorabstract := supplier.NewActorAbstract()

	datacenter := &NPDatacenter{
		ActorAbstract:   *actorabstract,
		trend:           nil,
		timer:           nil,
		buyoutprofile:   nil,
		profilename:     "",
		cronevent:       nil,
		picture:         "",
		reputationscore: 0.0,
	}

	return datacenter
}

func (self *NPDatacenter) Init(timer *timer.GameTimer, initialcapital float64, locationname string, trend *supplier.Trend, profilename string, name string, male bool) {
	log.Debug("NPDatacenter::Init(", timer, ",", initialcapital, ",", locationname, ",", trend, ",", profilename, ",", name, ",", male, ")")

	if self.cronevent != nil {
		self.timer.RemoveCron(self.cronevent)
	}

	self.ActorAbstract.Init(timer, locationname, name)

	// load buyout profile
	data, err := global.Asset("assets/npdatacenter/" + profilename)
	if err != nil {
		log.Error("NewNPDatacenter(): asset assets/npdatacenter/" + profilename + " not found")
		return
	}
	profile := make(map[string]BuyoutProfile)
	err = json.Unmarshal(data, &profile)
	if err != nil {
		log.Error("NewNPDatacenter(): asset assets/npdatacenter/" + profilename + " not json compatible")
		return
	}

	self.timer = timer
	self.trend = trend
	self.profilename = profilename
	self.buyoutprofile = profile
	self.picture = self.findPicture(male)
	self.reputationscore = 0.8 + rand.Float64()*0.2

	self.cronevent = timer.AddCron(1, 1, -1, func() {
		self.NewYearOperations()
	})

	// add some equity
	self.GetLedger().AddMovement(accounting.LedgerMovement{
		Description: "initial capital",
		Amount:      initialcapital,
		AccountFrom: "4561",
		AccountTo:   "5121",
		Date:        timer.CurrentTime,
	})
}

func (self *NPDatacenter) findPicture(male bool) string {
	sex := "male"
	if male == false {
		sex = "female"
	}
	if pictures, err := global.AssetDir("assets/faces/" + sex); err == nil {
		return sex + "/" + pictures[rand.Int()%len(pictures)]
	}
	return ""
}

func (self *NPDatacenter) LoadGame(timer *timer.GameTimer, trend *supplier.Trend, v map[string]interface{}) {
	log.Debug("NPDatacenter::LoadGame(", timer, ",", trend, ",", v, ")")
	if self.cronevent != nil {
		self.timer.RemoveCron(self.cronevent)
	}

	locationname := v["location"].(string)
	name := v["name"].(string)
	self.ActorAbstract.Init(timer, locationname, name)

	profilename := v["profile"].(string)
	// load buyout profile
	data, err := global.Asset("assets/npdatacenter/" + profilename)
	if err != nil {
		log.Error("NewNPDatacenter(): asset assets/npdatacenter/" + profilename + " not found")
		return
	}
	profile := make(map[string]BuyoutProfile)
	err = json.Unmarshal(data, &profile)
	if err != nil {
		log.Error("NewNPDatacenter(): asset assets/npdatacenter/" + profilename + " not json compatible")
		return
	}

	self.trend = trend
	self.profilename = profilename
	self.buyoutprofile = profile
	self.picture = v["picture"].(string)
	self.reputationscore = v["reputation"].(float64)

	self.GetLedger().Load(v["ledger"].(map[string]interface{}), self.GetLocation().Taxrate, self.GetLocation().Bankinterestrate)
	self.GetInventory().Load(v["inventory"].(map[string]interface{}))

	if offersinterface, ok := v["offers"]; ok {
		offers := offersinterface.([]interface{})
		for _, offer := range offers {
			self.LoadOffer(offer.(map[string]interface{}))
		}
	}

	self.cronevent = timer.AddCron(1, 1, -1, func() {
		self.NewYearOperations()
	})
}

func (self *NPDatacenter) Save() string {
	save := fmt.Sprintf(`{"location": "%s",`, self.GetLocationName()) + "\n"
	save += fmt.Sprintf(`"profile": "%s",`, self.profilename) + "\n"
	save += fmt.Sprintf(`"name": "%s",`, self.GetName()) + "\n"
	save += fmt.Sprintf(`"picture": "%s",`, self.picture) + "\n"
	save += fmt.Sprintf(`"reputation": %f,`, self.reputationscore) + "\n"
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
	save += fmt.Sprintf(`"ledger": %s`, self.GetLedger().Save()) + "}\n"
	return save
}

//
// NewYearOperations will trigger every year different actions, in particular
// - buy goods
// - create/refresh offers
func (self *NPDatacenter) NewYearOperations() {
	log.Debug("NPDatacenter::NewYearOperations()")
	// let's begin by removing all offers
	for _, o := range self.GetOffers() {
		self.RemoveOffer(o)
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

			currentCash := self.GetLedger().GetYearAccount(self.timer.CurrentTime.Year())["51"]

			nb := int32((currentCash * profile.Buyperyear) / unitprice)
			// if we can afford, then we buy it
			if nb > 0 {
				item := &supplier.CartItem{
					Typeitem:   supplier.PRODUCT_SERVER,
					Serverconf: &serverconf,
					Unitprice:  unitprice,
					Nb:         nb,
				}
				self.GetInventory().Cart = append(self.GetInventory().Cart, item)
				self.GetLedger().BuyProduct(fmt.Sprintf("buying %s", profilename), self.timer.CurrentTime, item.Unitprice*float64(nb))
				items := self.GetInventory().BuyCart(self.timer.CurrentTime)
				self.GetInventory().Cart = make([]*supplier.CartItem, 0)

				// create the offer
				var pool supplier.ServerPool
				// check the first available pool
				if profile.Vps {
					pool = self.GetInventory().GetDefaultVpsPool()
				} else {
					pool = self.GetInventory().GetDefaultPhysicalPool()
				}

				// it shouldn't happen...
				if pool == nil {
					log.Error("NPDatacenter::NewYearOperations(): We didn't find a correct (default) pool!")
					continue
				}

				// add newly create InventoryItem to the pool
				for _, i := range items {
					self.GetInventory().AssignPool(i, pool)
				}

				offer := &supplier.ServerOffer{
					Active:   true,
					Name:     profilename,
					Actor:    self,
					Pool:     pool,
					Vps:      profile.Vps,
					Nbcores:  serverconf.NbCore * serverconf.NbProcessors,
					Ramsize:  serverconf.NbSlotRam * serverconf.RamSize,
					Disksize: serverconf.NbDisks * serverconf.DiskSize,
					Vt:       serverconf.VtSupport,
					Price:    unitprice * profile.Margin,
				}
				self.AddOffer(offer)
			}
		}
	}

	// hack: we (re)place all servers
	for i, item := range self.GetInventory().Items {
		item.Xplaced = 1
		item.Yplaced = 1
		item.Zplaced = i
	}
}
