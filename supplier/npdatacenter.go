package supplier

import (
	"encoding/json"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/timer"
)

const (
	CONFIGURATION_LOW    = 1
	CONFIGURATION_MEDIUM = 2
	CONFIGURATION_HIGH   = 3
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
	margin        float64
	configuration int
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
	// a default value
	location := AvailableLocation["siliconvalley"]

	if l, ok := AvailableLocation[locationid]; ok {
		location = l
	}

	// load buyout profile
	data, err := global.Asset("assets/npdatacenter/" + profilename)
	if err != nil {
		return nil
	}
	profile := make(map[string]BuyoutProfile)
	err = json.Unmarshal(data, &profile)
	if err != nil {
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
// - create/refresh services
func (self *NPDatacenter) NewYearOperations() {

}
