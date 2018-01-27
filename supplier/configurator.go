package supplier

import (
	"math"
	"time"
)

//
// different hardware vendor?
//
// 1,2 processors: 1u,2u.
// 2,4 processors: 2u,4u.
// blade: 2 processors/blade + 1 disk.
//
// each component: price x ~ 1.05
//
// power consumption:
// - processor+fan: 100W .
// - motherboard: 60W.
// - ram: 4 W / slot.
// - disk (spindle): 7W/disk.
//

type ServerConfType struct {
	ServerName     string
	NbProcessors   [2]int32
	NbDisks        [2]int32
	NbSlotRam      [2]int32
	BackplanePrice float64
	ServerSprite   string
	NbU            int32
	scrap          bool
}

var AvailableConfs = []ServerConfType{
	ServerConfType{
		ServerName:     "T1000",
		NbProcessors:   [2]int32{1, 2},
		NbDisks:        [2]int32{1, 2},
		NbSlotRam:      [2]int32{1, 4},
		BackplanePrice: 200,
		ServerSprite:   "tower",
		NbU:            -1,
		scrap:          false,
	},
	ServerConfType{
		ServerName:     "R100",
		NbProcessors:   [2]int32{1, 2},
		NbDisks:        [2]int32{1, 1},
		NbSlotRam:      [2]int32{1, 2},
		BackplanePrice: 1000,
		ServerSprite:   "server.1u",
		NbU:            1,
		scrap:          false,
	},
	ServerConfType{
		ServerName:     "R200",
		NbProcessors:   [2]int32{1, 2},
		NbDisks:        [2]int32{1, 4},
		NbSlotRam:      [2]int32{1, 4},
		BackplanePrice: 2000,
		ServerSprite:   "server.2u",
		NbU:            2,
		scrap:          false,
	},
	ServerConfType{
		ServerName:     "R400",
		NbProcessors:   [2]int32{1, 2},
		NbDisks:        [2]int32{1, 10},
		NbSlotRam:      [2]int32{1, 8},
		BackplanePrice: 3000,
		ServerSprite:   "server.4u",
		NbU:            4,
		scrap:          false,
	},
	ServerConfType{
		ServerName:     "R600",
		NbProcessors:   [2]int32{2, 4},
		NbDisks:        [2]int32{1, 6},
		NbSlotRam:      [2]int32{1, 8},
		BackplanePrice: 3000,
		ServerSprite:   "server.4u",
		NbU:            4,
		scrap:          false,
	},
	ServerConfType{
		ServerName:     "B100",
		NbProcessors:   [2]int32{8, 8},
		NbDisks:        [2]int32{8, 8},
		NbSlotRam:      [2]int32{32, 32},
		BackplanePrice: 6000,
		ServerSprite:   "server.blade.8u",
		NbU:            8,
		scrap:          false,
	},
	ServerConfType{
		ServerName:     "B200",
		NbProcessors:   [2]int32{16, 16},
		NbDisks:        [2]int32{8, 8},
		NbSlotRam:      [2]int32{32, 32},
		BackplanePrice: 8000,
		ServerSprite:   "server.blade.8u",
		NbU:            8,
		scrap:          false,
	},
	ServerConfType{
		ServerName:     "scrap1U",
		NbProcessors:   [2]int32{0, 0},
		NbDisks:        [2]int32{0, 0},
		NbSlotRam:      [2]int32{0, 0},
		BackplanePrice: 0,
		ServerSprite:   "server.1u",
		NbU:            1,
		scrap:          true,
	},
	ServerConfType{
		ServerName:     "scrap2U",
		NbProcessors:   [2]int32{0, 0},
		NbDisks:        [2]int32{0, 0},
		NbSlotRam:      [2]int32{0, 0},
		BackplanePrice: 0,
		ServerSprite:   "server.2u",
		NbU:            2,
		scrap:          true,
	},
	ServerConfType{
		ServerName:     "scrap4U",
		NbProcessors:   [2]int32{0, 0},
		NbDisks:        [2]int32{0, 0},
		NbSlotRam:      [2]int32{0, 0},
		BackplanePrice: 0,
		ServerSprite:   "server.4u",
		NbU:            4,
		scrap:          true,
	},
	ServerConfType{
		ServerName:     "scrap8U",
		NbProcessors:   [2]int32{0, 0},
		NbDisks:        [2]int32{0, 0},
		NbSlotRam:      [2]int32{0, 0},
		BackplanePrice: 0,
		ServerSprite:   "server.blade.8u",
		NbU:            8,
		scrap:          true,
	},
}

func GetServerConfTypeByName(name string) *ServerConfType {
	for _, conftype := range AvailableConfs {
		if conftype.ServerName == name {
			return &conftype
		}
	}
	return nil
}

//
// based on the different type of chassis available
// and the vendor(s) options, the final server conf
// will have these caracteristics
//
type ServerConf struct {
	NbProcessors int32 //chosen
	NbCore       int32 // depend on the current trend
	VtSupport    bool  // depend on the current trend
	NbDisks      int32 // chosen
	NbSlotRam    int32 // chosen
	DiskSize     int32 // 3 options: Trend.Disksize: 1,1/2,1/4
	RamSize      int32 // 4 options: Trend.Ramsize: 1,1/2,1/4,1/8
	ConfType     *ServerConfType
	//	PricePaid    float64
}

// PowerConsumption return the Wh consumption
func (self *ServerConf) PowerConsumption() float64 {
	if self.ConfType.scrap == true {
		return 0
	}
	var consumption float64
	// http://www.vertatique.com/average-power-use-server
	consumption = float64(self.NbProcessors)*50.0 +
		float64(self.NbDisks)*42.0 +
		float64(self.NbSlotRam)*22.0 +
		218
	return consumption
}

func (self *ServerConf) Price(trend *Trend, now time.Time) float64 {
	var price float64
	complexity := float64(self.NbProcessors)/10 + float64(self.NbDisks)/20 + float64(self.NbSlotRam)/40 + 1
	price = self.ConfType.BackplanePrice +
		trend.Cpuprice.CurrentValue(now)*float64(self.NbProcessors)*float64(self.NbCore) +
		trend.Diskprice.CurrentValue(now)*float64(self.NbDisks*self.DiskSize)/1000 +
		trend.Ramprice.CurrentValue(now)*float64(self.NbSlotRam*self.RamSize)/1000
	return math.Floor(price * complexity)
}
