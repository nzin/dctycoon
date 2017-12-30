package supplier

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"time"

	"github.com/nzin/dctycoon/timer"
)

//
// We will have different trends to follow:
//
// - number of cores/CPU.
//
// - vt.
//
// - disk size / plateau.
//
// - ram size / slot.
//
type TrendItem struct {
	Pit       time.Time
	Value     int32
	ShortDesc string
	LongDesc  string
}

type TrendList []TrendItem

func (self TrendList) Len() int {
	return len(self)
}

func (self TrendList) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self TrendList) Less(i, j int) bool {
	return self[i].Pit.Before(self[j].Pit)
}

func (self TrendList) Sort() {
	sort.Sort(self)
}

func (self TrendList) CurrentValue(now time.Time) int32 {
	if len(self) == 0 {
		panic("no elements in the array")
	}

	index := 0
	for index < len(self) && self[index].Pit.Before(now) {
		index++
	}
	if index == 0 {
		return self[0].Value
	}
	return self[index-1].Value
}

// Other type of trend: price trends (+noise):
//
// - cpu price / core + noise.
//
// - disk price / Go + noise.
//
// - ram price / Go + noise.
//
type PriceTrendItem struct {
	Pit   time.Time
	Value float64
}

type PriceTrendList []PriceTrendItem

func (self PriceTrendList) Len() int {
	return len(self)
}

func (self PriceTrendList) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self PriceTrendList) Less(i, j int) bool {
	return self[i].Pit.Before(self[j].Pit)
}

func (self PriceTrendList) Sort() {
	sort.Sort(self)
}

type PriceTrend struct {
	Trend PriceTrendList
	Noise PriceTrendList
}

//
// function to compute the price+noise for a given date
//
func (self *PriceTrend) CurrentValue(now time.Time) float64 {
	if (self.Trend == nil) || (len(self.Trend) == 0) {
		temp := make(PriceTrendList, 1)
		temp[0].Pit = now
		temp[0].Value = 0.0
		self.Trend = temp
	}

	index := 0
	for index < len(self.Trend) && self.Trend[index].Pit.Before(now) {
		index++
	}

	var Value float64
	if index == 0 {
		Value = self.Trend[0].Value
	} else if index == len(self.Trend) {
		Value = self.Trend[index-1].Value
	} else {
		interval := (self.Trend[index].Pit.Sub(self.Trend[index-1].Pit)).Hours()
		since := now.Sub(self.Trend[index-1].Pit).Hours()
		Value = self.Trend[index-1].Value*((interval-since)/interval) + self.Trend[index].Value*(since/interval)
	}

	// now compute the noise

	if (self.Noise == nil) || len(self.Noise) == 0 {
		temp := make(PriceTrendList, 1)
		temp[0].Pit = now
		temp[0].Value = 0.0
		self.Noise = temp
	}
	endarray := len(self.Noise) - 1
	for now.Before((self.Noise)[endarray].Pit) == false {
		random := rand.Float64()
		if random < 0.1 {
			random = 0.1
		}
		elt := PriceTrendItem{
			Pit:   (self.Noise)[endarray].Pit.AddDate(0, 0, int(100*random)),
			Value: 1.0 - random,
		}
		self.Noise = append(self.Noise, elt)
		endarray = len(self.Noise) - 1
	}

	index = 0
	for index < len(self.Noise) && (self.Noise)[index].Pit.Before(now) {
		index++
	}
	var noise float64
	if index == 0 {
		noise = (self.Noise)[0].Value
	} else if index == len(self.Noise) {
		return (self.Noise)[index-1].Value
	} else {
		fmt.Println(index)
		interval := ((self.Noise)[index].Pit.Sub((self.Noise)[index-1].Pit)).Hours()
		since := now.Sub((self.Noise)[index-1].Pit).Hours()
		noise = (self.Noise)[index-1].Value*((interval-since)/interval) + (self.Noise)[index].Value*(since/interval)
	}
	return Value * (noise + 1.0)
}

func PriceTrendListLoad(noise []interface{}, trend []PriceTrendItem) PriceTrend {
	pt := PriceTrend{}

	pt.Trend = trend

	nl := make(PriceTrendList, 0)

	for _, n := range noise {
		if data, ok := n.(map[string]interface{}); ok {
			ne := data
			var year, month, day int
			fmt.Sscanf(ne["pit"].(string), "%d-%d-%d", &year, &month, &day)
			nl = append(nl, PriceTrendItem{
				Pit:   time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC),
				Value: ne["value"].(float64),
			})
		}
	}
	nl.Sort()
	pt.Noise = nl

	return pt
}

func PriceTrendListSave(pt PriceTrend) string {
	str := `[`
	for i, ne := range pt.Noise {
		if i > 0 {
			str += ","
		}
		str += fmt.Sprintf(`{"pit":"%d-%d-%d", "value":%g}`, ne.Pit.Year(), ne.Pit.Month(), ne.Pit.Day(), ne.Value)
	}
	str += `]`
	return str
}

//
// global structure to store all the trends
//
type Trend struct {
	Corepercpu TrendList
	Vt         TrendList
	Disksize   TrendList
	Ramsize    TrendList

	Cpuprice  PriceTrend
	Diskprice PriceTrend
	Ramprice  PriceTrend
}

var initVt = []TrendItem{
	TrendItem{Pit: time.Date(1979, time.Month(01), 01, 0, 0, 0, 0, time.UTC), Value: 0, ShortDesc: "", LongDesc: ""},
	TrendItem{Pit: time.Date(2005, time.Month(01), 01, 0, 0, 0, 0, time.UTC), Value: 1, ShortDesc: "UMD introduces VT", LongDesc: "UMD introduce Virtualization capability with VT enabled microprocessors!"},
}

var initCorepercpu = []TrendItem{
	TrendItem{Pit: time.Date(1979, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 1, ShortDesc: "", LongDesc: ""},
	TrendItem{Pit: time.Date(2006, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 2, ShortDesc: "multicore processor arrived", LongDesc: "Untel reveals the first dual core processor architecture for x86"},
	TrendItem{Pit: time.Date(2009, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 4, ShortDesc: "4 core processor on sales", LongDesc: "Untel strikes again, but now with a 4 core-processor architecture, on sales now!"},
	TrendItem{Pit: time.Date(2011, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 6, ShortDesc: "6 core processor on sales", LongDesc: "Not satisfied with 4 core processors? The limits has been pushed to a new 6-core processors, and rumors talked about a 8-core processor soon"},
	TrendItem{Pit: time.Date(2012, time.Month(6), 1, 0, 0, 0, 0, time.UTC), Value: 8, ShortDesc: "8 core processor on sales", LongDesc: "A 8 core processor has been released, where and when Moore law will comes to an end?"},
	TrendItem{Pit: time.Date(2017, time.Month(6), 1, 0, 0, 0, 0, time.UTC), Value: 12, ShortDesc: "12 core processor ready", LongDesc: "A 12 core processor have been reveleaed with a 16 core on the horizon"},
	TrendItem{Pit: time.Date(2018, time.Month(9), 1, 0, 0, 0, 0, time.UTC), Value: 16, ShortDesc: "16 core processor just for you", LongDesc: "After the 12 core, here is the 16 core to cram more CPU into your blade servers"},
}

// size : http://www.pcworld.com/article/127105/article.html
var initDisksize = []TrendItem{ // in MB
	TrendItem{Pit: time.Date(1983, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 10, ShortDesc: "RO352 10MB diskdrive", LongDesc: "Rodime releases the first 3.5-inch hard drive; the RO352 includes two platters and stores 10MB"},
	TrendItem{Pit: time.Date(1992, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 170, ShortDesc: "Western Digital 170MB", LongDesc: "Western Digital 170MB arrives"},
	TrendItem{Pit: time.Date(1993, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 270, ShortDesc: "Seagate 270MB", LongDesc: "Seagate 270MB disk on sales"},
	TrendItem{Pit: time.Date(1994, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 540, ShortDesc: "Western Digital 540MB", LongDesc: "Western Digital 540MB disk drive on sales"},
	TrendItem{Pit: time.Date(1995, time.Month(6), 1, 0, 0, 0, 0, time.UTC), Value: 1000, ShortDesc: "Seagate unveils 1G DD", LongDesc: "The first 1G diskdrive from Seagate for ~ $849"},
	TrendItem{Pit: time.Date(1997, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 3200, ShortDesc: "Quantum 3.2G DD on sale", LongDesc: "Quantum unveils a 3200MB disk"},
	TrendItem{Pit: time.Date(1998, time.Month(7), 1, 0, 0, 0, 0, time.UTC), Value: 6400, ShortDesc: "Quantum 6.4G DD", LongDesc: "Quantum strikes back with a 6.4G disk drive"},
	TrendItem{Pit: time.Date(2000, time.Month(2), 1, 0, 0, 0, 0, time.UTC), Value: 12000, ShortDesc: "Fujitsu 12G disk drive", LongDesc: "Fujitsu (and soon competitors) reveals a 12000MB disk!"},
	TrendItem{Pit: time.Date(2001, time.Month(7), 1, 0, 0, 0, 0, time.UTC), Value: 60000, ShortDesc: "Maxtor 60G disk drive", LongDesc: "Quantum leap for Maxtor with this 60G disk drive"},
	TrendItem{Pit: time.Date(2003, time.Month(4), 1, 0, 0, 0, 0, time.UTC), Value: 120000, ShortDesc: "Maxtor 120G DD unveils", LongDesc: "Maxtor enhances his offer with a new 120G disk drive"},
	TrendItem{Pit: time.Date(2005, time.Month(12), 1, 0, 0, 0, 0, time.UTC), Value: 500000, ShortDesc: "Seagate 500G disk", LongDesc: "Seagate Barracuda 7200.9 500GB is on sale, now!"},
	TrendItem{Pit: time.Date(2006, time.Month(9), 1, 0, 0, 0, 0, time.UTC), Value: 750000, ShortDesc: "Seagate 750G disk", LongDesc: "Seagate Barracuda 7200.10 750GB"},
	TrendItem{Pit: time.Date(2007, time.Month(3), 1, 0, 0, 0, 0, time.UTC), Value: 1000000, ShortDesc: "First Terabyte DD", LongDesc: "Beyond Micro Monster Mobile 1TB"},
	TrendItem{Pit: time.Date(2008, time.Month(4), 1, 0, 0, 0, 0, time.UTC), Value: 1500000, ShortDesc: "1500GB disk drive", LongDesc: "New leap forward with new 1.5TB disk on sale"},
	TrendItem{Pit: time.Date(2009, time.Month(10), 1, 0, 0, 0, 0, time.UTC), Value: 2000000, ShortDesc: "2To disk drive", LongDesc: "New density barrier reached with new 2To models"},
	TrendItem{Pit: time.Date(2010, time.Month(7), 1, 0, 0, 0, 0, time.UTC), Value: 3000000, ShortDesc: "3To disk drive", LongDesc: "A new record: 3To in one disk"},
	TrendItem{Pit: time.Date(2011, time.Month(2), 1, 0, 0, 0, 0, time.UTC), Value: 4000000, ShortDesc: "4To disk drive", LongDesc: "How many movies do you want to stores? :-)"},
	TrendItem{Pit: time.Date(2013, time.Month(2), 1, 0, 0, 0, 0, time.UTC), Value: 6000000, ShortDesc: "6To disk drive", LongDesc: "You can rush again on your local vendor to acquire new 6To disks"},
	TrendItem{Pit: time.Date(2014, time.Month(11), 1, 0, 0, 0, 0, time.UTC), Value: 8000000, ShortDesc: "8To disk drive", LongDesc: "The future has no limit with this new 8To disk drive"},
	TrendItem{Pit: time.Date(2017, time.Month(5), 1, 0, 0, 0, 0, time.UTC), Value: 12000000, ShortDesc: "12To disk drive", LongDesc: "The bigger the better with this new 12To disk drive"},
}
var initRamsize = []TrendItem{ // in MB
	TrendItem{Pit: time.Date(1994, time.Month(04), 01, 0, 0, 0, 0, time.UTC), Value: 4, ShortDesc: "4MB SDRAM on sale", LongDesc: "4MB SDRAM memory on sale"},
	TrendItem{Pit: time.Date(1996, time.Month(04), 01, 0, 0, 0, 0, time.UTC), Value: 8, ShortDesc: "8MB SDRAM on sale", LongDesc: "8MB SDRAM memory on sale"},
	TrendItem{Pit: time.Date(1997, time.Month(11), 01, 0, 0, 0, 0, time.UTC), Value: 16, ShortDesc: "16MB SDRAM on sale", LongDesc: "16MB SDRAM memory on sale"},
	TrendItem{Pit: time.Date(1999, time.Month(04), 01, 0, 0, 0, 0, time.UTC), Value: 32, ShortDesc: "32MB SDRAM on sale", LongDesc: "32MB SDRAM memory on sale"},
	TrendItem{Pit: time.Date(2002, time.Month(11), 01, 0, 0, 0, 0, time.UTC), Value: 64, ShortDesc: "64MB SDRAM on sale", LongDesc: "64MB SDRAM memory on sale"},
	TrendItem{Pit: time.Date(2004, time.Month(04), 01, 0, 0, 0, 0, time.UTC), Value: 128, ShortDesc: "128MB SDRAM on sale", LongDesc: "128MB SDRAM memory on sale"},
	TrendItem{Pit: time.Date(2006, time.Month(11), 01, 0, 0, 0, 0, time.UTC), Value: 256, ShortDesc: "256MB SDRAM on sale", LongDesc: "256MB SDRAM memory on sale"},
	TrendItem{Pit: time.Date(2008, time.Month(04), 01, 0, 0, 0, 0, time.UTC), Value: 512, ShortDesc: "512MB SDRAM on sale", LongDesc: "512MB SDRAM memory on sale"},
	TrendItem{Pit: time.Date(2009, time.Month(11), 01, 0, 0, 0, 0, time.UTC), Value: 1024, ShortDesc: "1GB SDRAM on sale", LongDesc: "1GB SDRAM memory on sale"},
	TrendItem{Pit: time.Date(2011, time.Month(04), 01, 0, 0, 0, 0, time.UTC), Value: 2048, ShortDesc: "2GB SDRAM on sale", LongDesc: "2GB SDRAM memory     on sale"},
	TrendItem{Pit: time.Date(2012, time.Month(11), 01, 0, 0, 0, 0, time.UTC), Value: 4096, ShortDesc: "4GB SDRAM on sale", LongDesc: "4GB SDRAM memory     on sale"},
	TrendItem{Pit: time.Date(2014, time.Month(04), 01, 0, 0, 0, 0, time.UTC), Value: 8192, ShortDesc: "8GB SDRAM on sale", LongDesc: "8GB SDRAM memory     on sale"},
	TrendItem{Pit: time.Date(2015, time.Month(11), 01, 0, 0, 0, 0, time.UTC), Value: 16384, ShortDesc: "16GB SDRAM on sale", LongDesc: "16GB SDRAM memory on sale"},
	TrendItem{Pit: time.Date(2017, time.Month(04), 01, 0, 0, 0, 0, time.UTC), Value: 32768, ShortDesc: "32GB SDRAM on sale", LongDesc: "32GB SDRAM memory on sale"},
	TrendItem{Pit: time.Date(2018, time.Month(11), 01, 0, 0, 0, 0, time.UTC), Value: 65536, ShortDesc: "64GB SDRAM on sale", LongDesc: "64GB SDRAM memory on sale"},
}

// source: http://www.mkomo.com/cost-per-gigabyte
var diskPriceTrend = []PriceTrendItem{ // $/Go
	PriceTrendItem{Pit: time.Date(1981, time.Month(11), 1, 0, 0, 0, 0, time.UTC), Value: 340000}, // Seagate 5M
	PriceTrendItem{Pit: time.Date(1983, time.Month(12), 1, 0, 0, 0, 0, time.UTC), Value: 190000}, // Xcomp 10M
	PriceTrendItem{Pit: time.Date(1984, time.Month(3), 1, 0, 0, 0, 0, time.UTC), Value: 170000},  // Tandon 10M
	PriceTrendItem{Pit: time.Date(1984, time.Month(5), 1, 0, 0, 0, 0, time.UTC), Value: 80000},   // Pegaus 23M
	PriceTrendItem{Pit: time.Date(1985, time.Month(7), 1, 0, 0, 0, 0, time.UTC), Value: 71000},   // First class peripherals 10M
	PriceTrendItem{Pit: time.Date(1987, time.Month(10), 1, 0, 0, 0, 0, time.UTC), Value: 45000},  // Iomega 40 M
	PriceTrendItem{Pit: time.Date(1988, time.Month(5), 1, 0, 0, 0, 0, time.UTC), Value: 30000},   // ?? 60M
	PriceTrendItem{Pit: time.Date(1989, time.Month(9), 1, 0, 0, 0, 0, time.UTC), Value: 12000},   // ??
	PriceTrendItem{Pit: time.Date(1990, time.Month(9), 1, 0, 0, 0, 0, time.UTC), Value: 9000},    // ??
	PriceTrendItem{Pit: time.Date(1991, time.Month(9), 1, 0, 0, 0, 0, time.UTC), Value: 7000},    // ??
	PriceTrendItem{Pit: time.Date(1992, time.Month(9), 1, 0, 0, 0, 0, time.UTC), Value: 4000},    // ??
	PriceTrendItem{Pit: time.Date(1993, time.Month(9), 1, 0, 0, 0, 0, time.UTC), Value: 2000},    // ??
	PriceTrendItem{Pit: time.Date(1994, time.Month(9), 1, 0, 0, 0, 0, time.UTC), Value: 950},     // ??
	PriceTrendItem{Pit: time.Date(1995, time.Month(4), 1, 0, 0, 0, 0, time.UTC), Value: 756},     // ?? 1000M
	PriceTrendItem{Pit: time.Date(1996, time.Month(6), 1, 0, 0, 0, 0, time.UTC), Value: 295},     // Western Digital 1600 M
	PriceTrendItem{Pit: time.Date(1997, time.Month(8), 13, 0, 0, 0, 0, time.UTC), Value: 141},    // Western Digital 4000 M
	PriceTrendItem{Pit: time.Date(1998, time.Month(1), 16, 0, 0, 0, 0, time.UTC), Value: 95.20},  // Maxtor 6400 M
	PriceTrendItem{Pit: time.Date(1998, time.Month(5), 11, 0, 0, 0, 0, time.UTC), Value: 58.90},  // Fujitsu 6400 M
	PriceTrendItem{Pit: time.Date(1999, time.Month(2), 26, 0, 0, 0, 0, time.UTC), Value: 37.70},  // Maxtor 8400 M
	PriceTrendItem{Pit: time.Date(1999, time.Month(2), 26, 0, 0, 0, 0, time.UTC), Value: 37.70},  // Maxtor 8400 M
	PriceTrendItem{Pit: time.Date(1999, time.Month(5), 27, 0, 0, 0, 0, time.UTC), Value: 24.50},  // Fujitsu UDMA 17.3 G
	PriceTrendItem{Pit: time.Date(1999, time.Month(10), 1, 0, 0, 0, 0, time.UTC), Value: 20.60},  // Western Digital 27.3G
	PriceTrendItem{Pit: time.Date(1999, time.Month(10), 1, 0, 0, 0, 0, time.UTC), Value: 20.60},  // Western Digital 27.3G
	PriceTrendItem{Pit: time.Date(1999, time.Month(12), 1, 0, 0, 0, 0, time.UTC), Value: 16.30},  // Fujitsu IDE 27.3G
	PriceTrendItem{Pit: time.Date(2000, time.Month(4), 1, 0, 0, 0, 0, time.UTC), Value: 13.00},   // Maxtor UDMA 36.5G
	PriceTrendItem{Pit: time.Date(2000, time.Month(8), 1, 0, 0, 0, 0, time.UTC), Value: 10.90},   // Maxtor 40.9G
	PriceTrendItem{Pit: time.Date(2000, time.Month(10), 27, 0, 0, 0, 0, time.UTC), Value: 7.30},  // Maxtor 81.9G
	PriceTrendItem{Pit: time.Date(2001, time.Month(11), 30, 0, 0, 0, 0, time.UTC), Value: 2.99},  // Western Digital 100G
	PriceTrendItem{Pit: time.Date(2002, time.Month(9), 6, 0, 0, 0, 0, time.UTC), Value: 2.59},    // Western Digital 120G
	PriceTrendItem{Pit: time.Date(2003, time.Month(11), 29, 0, 0, 0, 0, time.UTC), Value: 1.61},  // Maxtor Seria ATA 120G
	PriceTrendItem{Pit: time.Date(2004, time.Month(3), 27, 0, 0, 0, 0, time.UTC), Value: 1.70},   // Western Digital Caviar 250G
	PriceTrendItem{Pit: time.Date(2004, time.Month(12), 4, 0, 0, 0, 0, time.UTC), Value: 0.70},   // Barracuda 400G
	PriceTrendItem{Pit: time.Date(2005, time.Month(8), 29, 0, 0, 0, 0, time.UTC), Value: 0.75},   // ?? 400G
	PriceTrendItem{Pit: time.Date(2006, time.Month(7), 5, 0, 0, 0, 0, time.UTC), Value: 0.60},    // Seagate Barracuda 500G
	PriceTrendItem{Pit: time.Date(2008, time.Month(1), 13, 0, 0, 0, 0, time.UTC), Value: 0.27},   // Seagate Barracuda 750G
	PriceTrendItem{Pit: time.Date(2009, time.Month(7), 24, 0, 0, 0, 0, time.UTC), Value: 0.14},   // HITACHI 1000G
	PriceTrendItem{Pit: time.Date(2010, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 0.07},    // ?? 2000G
	PriceTrendItem{Pit: time.Date(2010, time.Month(7), 1, 0, 0, 0, 0, time.UTC), Value: 0.05},    // ?? 3000G
	PriceTrendItem{Pit: time.Date(2011, time.Month(2), 1, 0, 0, 0, 0, time.UTC), Value: 0.03},    // ?? 4000G
	PriceTrendItem{Pit: time.Date(2013, time.Month(2), 1, 0, 0, 0, 0, time.UTC), Value: 0.02},    // ?? 6000G
	PriceTrendItem{Pit: time.Date(2014, time.Month(11), 1, 0, 0, 0, 0, time.UTC), Value: 0.015},  // ?? 8000G
	PriceTrendItem{Pit: time.Date(2017, time.Month(5), 1, 0, 0, 0, 0, time.UTC), Value: 0.01},    // ?? 12000G
}

// based on https://arstechnica.com/gadgets/2016/11/how-cheap-ram-changes-computing/
var ramPriceTrend = []PriceTrendItem{ // $/Go
	PriceTrendItem{Pit: time.Date(1980, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 10000000},
	PriceTrendItem{Pit: time.Date(1985, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 1000000},
	PriceTrendItem{Pit: time.Date(1990, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 100000},
	PriceTrendItem{Pit: time.Date(1995, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 50000},
	PriceTrendItem{Pit: time.Date(2000, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 1000},
	PriceTrendItem{Pit: time.Date(2005, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 500},
	PriceTrendItem{Pit: time.Date(2010, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 50},
	PriceTrendItem{Pit: time.Date(2015, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 8},
	PriceTrendItem{Pit: time.Date(2020, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 1},
}

// imaginaries values (I didn't find good data on internet)
var cpucorePriceTrend = []PriceTrendItem{ // $/core
	PriceTrendItem{Pit: time.Date(1979, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 4000},
	PriceTrendItem{Pit: time.Date(1990, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 1000},
	PriceTrendItem{Pit: time.Date(2000, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 500},
	PriceTrendItem{Pit: time.Date(2006, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 400},
	PriceTrendItem{Pit: time.Date(2012, time.Month(1), 1, 0, 0, 0, 0, time.UTC), Value: 100},
}

func TrendLoad(json map[string]interface{}, publishservice timer.EventPublisherService, timer *timer.GameTimer) *Trend {
	cpupricenoise := make([]interface{}, 0, 0)
	diskpricenoise := make([]interface{}, 0, 0)
	rampricenoise := make([]interface{}, 0, 0)

	if array, ok := json["cpupricenoise"]; ok == true {
		if reflect.TypeOf(array).Kind() == reflect.Slice {
			s := reflect.ValueOf(array)
			for i := 0; i < s.Len(); i++ {
				cpupricenoise = append(cpupricenoise, s.Index(i).Interface())
			}
		}
	}

	if array, ok := json["diskpricenoise"]; ok == true {
		if reflect.TypeOf(array).Kind() == reflect.Slice {
			s := reflect.ValueOf(array)
			for i := 0; i < s.Len(); i++ {
				diskpricenoise = append(diskpricenoise, s.Index(i).Interface())
			}
		}
	}

	if array, ok := json["rampricenoise"]; ok == true {
		if reflect.TypeOf(array).Kind() == reflect.Slice {
			s := reflect.ValueOf(array)
			for i := 0; i < s.Len(); i++ {
				rampricenoise = append(rampricenoise, s.Index(i).Interface())
			}
		}
	}

	t := &Trend{
		Corepercpu: initCorepercpu,
		Vt:         initVt,
		Disksize:   initDisksize,
		Ramsize:    initRamsize,

		Cpuprice:  PriceTrendListLoad(cpupricenoise, cpucorePriceTrend),
		Diskprice: PriceTrendListLoad(diskpricenoise, diskPriceTrend),
		Ramprice:  PriceTrendListLoad(rampricenoise, ramPriceTrend),
	}

	for _, core := range initCorepercpu {
		c := core
		timer.AddEvent(c.Pit, func() {
			publishservice.Publish(c.ShortDesc, c.LongDesc)
		})
	}

	for _, vt := range initVt {
		v := vt
		timer.AddEvent(v.Pit, func() {
			publishservice.Publish(v.ShortDesc, v.LongDesc)
		})
	}

	for _, dd := range initDisksize {
		d := dd
		timer.AddEvent(d.Pit, func() {
			publishservice.Publish(d.ShortDesc, d.LongDesc)
		})
	}

	for _, ram := range initRamsize {
		r := ram
		timer.AddEvent(r.Pit, func() {
			publishservice.Publish(r.ShortDesc, r.LongDesc)
		})
	}

	return t
}

func TrendSave(t *Trend) string {
	str := "{\n"
	str += fmt.Sprintf(`"cpupricenoise": %s,`, PriceTrendListSave(t.Cpuprice)) + "\n"
	str += fmt.Sprintf(`"diskpricenoise": %s,`, PriceTrendListSave(t.Diskprice)) + "\n"
	str += fmt.Sprintf(`"rampricenoise": %s`, PriceTrendListSave(t.Ramprice)) + "\n"
	str += "}"
	return str
}
