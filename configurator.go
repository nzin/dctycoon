package dctycoon

import (
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
}

var AvailableConfs = []ServerConfType{
	ServerConfType{
		ServerName:     "small",
		NbProcessors:   [2]int32{1, 2},
		NbDisks:        [2]int32{1, 1},
		NbSlotRam:      [2]int32{1, 2},
		BackplanePrice: 1000,
		ServerSprite:   "server.1u",
		NbU:            1,
	},
	ServerConfType{
		ServerName:     "medium",
		NbProcessors:   [2]int32{1, 2},
		NbDisks:        [2]int32{1, 4},
		NbSlotRam:      [2]int32{1, 4},
		BackplanePrice: 2000,
		ServerSprite:   "server.2u",
		NbU:            2,
	},
	ServerConfType{
		ServerName:     "large.d1",
		NbProcessors:   [2]int32{1, 2},
		NbDisks:        [2]int32{1, 10},
		NbSlotRam:      [2]int32{1, 8},
		BackplanePrice: 3000,
		ServerSprite:   "server.4u",
		NbU:            4,
	},
	ServerConfType{
		ServerName:     "large.c1",
		NbProcessors:   [2]int32{2, 4},
		NbDisks:        [2]int32{1, 6},
		NbSlotRam:      [2]int32{1, 8},
		BackplanePrice: 3000,
		ServerSprite:   "server.4u",
		NbU:            4,
	},
	ServerConfType{
		ServerName:     "blade1",
		NbProcessors:   [2]int32{8, 8},
		NbDisks:        [2]int32{8, 8},
		NbSlotRam:      [2]int32{16, 16},
		BackplanePrice: 6000,
		ServerSprite:   "server.blade",
		NbU:            8,
	},
	ServerConfType{
		ServerName:     "blade2",
		NbProcessors:   [2]int32{16, 16},
		NbDisks:        [2]int32{8, 8},
		NbSlotRam:      [2]int32{16, 16},
		BackplanePrice: 8000,
		ServerSprite:   "server.blade",
		NbU:            8,
	},
}

type ServerConf struct {
	NbProcessors int32
	NbDisks      int32
	NbSlotRam    int32
	DiskSize     int32 // 3 options: Trend.Disksize: 1,1/2,1/4
	RamSize      int32 // 4 options: Trend.Ramsize: 1,1/2,1/4,1/8
	ConfType     *ServerConfType
}

func (self *ServerConf) PowerConsumption() float64 {
	var consumption float64
	consumption = float64(self.NbProcessors) * 100.0 +
		float64(self.NbDisks) * 7.0 +
		float64(self.NbSlotRam) * 4.0 +
		60.0
	return consumption
}

func (self *ServerConf) Price(now time.Time) float64 {
	var price float64
	complexity := float64(self.NbProcessors) / 10 + float64(self.NbDisks) / 20 + float64(self.NbSlotRam) / 40 + 1
	price = self.ConfType.BackplanePrice +
		Trends.Cpuprice.CurrentValue(now) * float64(self.NbProcessors) +
		Trends.Diskprice.CurrentValue(now) * float64(self.NbDisks * self.DiskSize) +
		Trends.Ramprice.CurrentValue(now) * float64(self.NbSlotRam * self.RamSize)
	return price * complexity
}



