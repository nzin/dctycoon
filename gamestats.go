package dctycoon

import (
	"fmt"
	"reflect"
	"time"

	"github.com/nzin/dctycoon/supplier"
	log "github.com/sirupsen/logrus"
)

type PowerStat struct {
	consumption float64
	generation  float64
	provided    float64
	date        time.Time
}

func (self *PowerStat) Save() string {
	str := "{"
	str += fmt.Sprintf("\"date\": \"%d-%d-%d\",", self.date.Year(), self.date.Month(), self.date.Day())
	str += fmt.Sprintf("\"consumption\": %f,", self.consumption)
	str += fmt.Sprintf("\"generation\": %f,", self.generation)
	str += fmt.Sprintf("\"provided\": %f", self.provided)
	return str + "}"
}

func NewPowerStat(v map[string]interface{}) *PowerStat {
	var year, month, day int
	fmt.Sscanf(v["date"].(string), "%d-%d-%d", &year, &month, &day)

	ps := &PowerStat{
		consumption: v["consumption"].(float64),
		generation:  v["generation"].(float64),
		provided:    v["provided"].(float64),
		date:        time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC),
	}
	return ps
}

type ServerDemandStat struct {
	ramsize  int32
	nbcores  int32
	disksize int32
	nb       int32
}

type DemandStat struct {
	serverdemands []*ServerDemandStat
	date          time.Time
	price         float64
	buyer         string
}

func (self *DemandStat) Save() string {
	str := "{"
	str += fmt.Sprintf("\"date\": \"%d-%d-%d\",", self.date.Year(), self.date.Month(), self.date.Day())
	str += fmt.Sprintf("\"price\": %f,", self.price)
	str += fmt.Sprintf("\"buyer\": \"%s\",", self.buyer)
	str += "\"servers\": ["
	for i, s := range self.serverdemands {
		if i != 0 {
			str += ","
		}
		str += fmt.Sprintf(`{"ramsize": %d, "nbcores": %d, "disksize": %d, "nb":%d}`, s.ramsize, s.nbcores, s.disksize, s.nb)
	}
	return str + "]}"
}

func NewDemandStat(v map[string]interface{}) *DemandStat {
	var year, month, day int
	fmt.Sscanf(v["date"].(string), "%d-%d-%d", &year, &month, &day)

	ds := &DemandStat{
		price:         v["price"].(float64),
		date:          time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC),
		buyer:         v["buyer"].(string),
		serverdemands: make([]*ServerDemandStat, 0, 0),
	}
	servers := v["servers"].([]interface{})
	for _, s := range servers {
		server := s.(map[string]interface{})
		demand := &ServerDemandStat{
			ramsize:  int32(server["ramsize"].(float64)),
			nbcores:  int32(server["nbcores"].(float64)),
			disksize: int32(server["disksize"].(float64)),
			nb:       int32(server["nb"].(float64)),
		}
		ds.serverdemands = append(ds.serverdemands, demand)
	}
	return ds
}

// DemandStatSubscriber is mainly used by the MainStatsWidget
type DemandStatSubscriber interface {
	NewDemandStat(*DemandStat)
}

// DemandStatSubscriber is mainly used by the MainStatsWidget
type PowerStatSubscriber interface {
	NewPowerStat(*PowerStat)
}

type GameStats struct {
	demandstatsubscribers []DemandStatSubscriber
	demandsstats          []*DemandStat
	powerstatssubscribers []PowerStatSubscriber
	powerstats            []*PowerStat
}

func (self *GameStats) AddDemandStatSubscriber(subscriber DemandStatSubscriber) {
	for _, s := range self.demandstatsubscribers {
		if s == subscriber {
			return
		}
	}
	self.demandstatsubscribers = append(self.demandstatsubscribers, subscriber)
}

func (self *GameStats) RemoveDemandStatSubscriber(subscriber DemandStatSubscriber) {
	for i, s := range self.demandstatsubscribers {
		if s == subscriber {
			self.demandstatsubscribers = append(self.demandstatsubscribers[:i], self.demandstatsubscribers[i+1:]...)
			break
		}
	}
}

func (self *GameStats) TriggerDemandStat(t time.Time, demand *supplier.DemandInstance, serverbundle *supplier.ServerBundle) {
	log.Debug("GameStats::TriggerDemandStat(", t, ",", demand, ",", serverbundle, ")")

	stat := &DemandStat{
		date:          t,
		serverdemands: make([]*ServerDemandStat, 0, 0),
	}
	if serverbundle != nil {
		stat.buyer = serverbundle.Actor.GetName()
		for _, sc := range serverbundle.Contracts {
			stat.price += sc.Price
		}
	}
	for templatename, templatevalue := range demand.Template.Specs {
		sds := &ServerDemandStat{
			nb: demand.Nb[templatename],
		}
		for _, filter := range templatevalue.Filters {
			if reflect.TypeOf(filter) == reflect.TypeOf((**supplier.CriteriaFilterDisk)(nil)).Elem() {
				sds.disksize = (filter.(*supplier.CriteriaFilterDisk)).Disksize
			}
			if reflect.TypeOf(filter) == reflect.TypeOf((**supplier.CriteriaFilterRam)(nil)).Elem() {
				sds.ramsize = (filter.(*supplier.CriteriaFilterRam)).Ramsize
			}
			if reflect.TypeOf(filter) == reflect.TypeOf((**supplier.CriteriaFilterNbcores)(nil)).Elem() {
				sds.disksize = (filter.(*supplier.CriteriaFilterNbcores)).Nbcores
			}
		}
		stat.serverdemands = append(stat.serverdemands, sds)
	}
	self.demandsstats = append(self.demandsstats, stat)

	for _, s := range self.demandstatsubscribers {
		s.NewDemandStat(stat)
	}
}

func (self *GameStats) AddPowerStatSubscriber(subscriber PowerStatSubscriber) {
	for _, s := range self.powerstatssubscribers {
		if s == subscriber {
			return
		}
	}
	self.powerstatssubscribers = append(self.powerstatssubscribers, subscriber)
}

func (self *GameStats) RemovePowerStatSubscriber(subscriber PowerStatSubscriber) {
	for i, s := range self.powerstatssubscribers {
		if s == subscriber {
			self.powerstatssubscribers = append(self.powerstatssubscribers[:i], self.powerstatssubscribers[i+1:]...)
			break
		}
	}
}

func (self *GameStats) PowerChange(t time.Time, consumption, generation, provided float64) {
	stat := &PowerStat{
		date:        t,
		consumption: consumption,
		generation:  generation,
		provided:    provided,
	}
	self.powerstats = append(self.powerstats, stat)

	for _, s := range self.powerstatssubscribers {
		s.NewPowerStat(stat)
	}
}

func NewGameStats() *GameStats {
	gs := &GameStats{
		demandsstats:          make([]*DemandStat, 0, 0),
		demandstatsubscribers: make([]DemandStatSubscriber, 0, 0),
		powerstats:            make([]*PowerStat, 0, 0),
		powerstatssubscribers: make([]PowerStatSubscriber, 0, 0),
	}

	return gs
}

func (self *GameStats) InitGame(inventory *supplier.Inventory) {
	self.demandsstats = make([]*DemandStat, 0, 0)
	self.powerstats = make([]*PowerStat, 0, 0)
	inventory.AddPowerStatSubscriber(self)
}

func (self *GameStats) LoadGame(inventory *supplier.Inventory, stats map[string]interface{}) {
	log.Debug("GameStats::LoadGame(", stats, ")")

	self.demandsstats = make([]*DemandStat, 0, 0)
	inventory.AddPowerStatSubscriber(self)

	demandsstats := stats["demandsstats"].([]interface{})
	for _, d := range demandsstats {
		self.demandsstats = append(self.demandsstats, NewDemandStat(d.(map[string]interface{})))
	}
	self.powerstats = make([]*PowerStat, 0, 0)

	powerstats := stats["powerstats"].([]interface{})
	for _, d := range powerstats {
		self.powerstats = append(self.powerstats, NewPowerStat(d.(map[string]interface{})))
	}
}

func (self *GameStats) Save() string {
	log.Debug("GameStats::Save()")

	str := "{\n"
	str += "\"demandsstats\": ["
	for i, d := range self.demandsstats {
		if i != 0 {
			str += ",\n"
		}
		str += d.Save()
	}
	str += "],\n"
	str += "\"powerstats\": ["
	for i, ps := range self.powerstats {
		if i != 0 {
			str += ",\n"
		}
		str += ps.Save()
	}
	return str + "]}"
}
