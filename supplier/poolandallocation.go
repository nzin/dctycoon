package supplier

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
)

type PoolSubscriber interface {
	InventoryItemAdd(*InventoryItem)
	InventoryItemRemove(*InventoryItem)
	InventoryItemAllocate(*InventoryItem)
	InventoryItemRelease(*InventoryItem)
}

type ServerPool interface {
	GetName() string
	addInventoryItem(item *InventoryItem)
	IsInside(item *InventoryItem) bool
	removeInventoryItem(item *InventoryItem)
	Allocate(nbcores, ramsize, disksize int32, vt bool) *InventoryItem
	IsAllocated(item *InventoryItem) bool
	Release(item *InventoryItem, nbcores, ramsize, disksize int32)
	IsVps() bool
	HowManyFit(nbcores, ramsize, disksize int32, vt bool) int32
	AddPoolSubscriber(subscriber PoolSubscriber)
	RemovePoolSubscriber(subscriber PoolSubscriber)
}

type HardwareServerPool struct {
	Name            string
	pool            map[int32]*InventoryItem
	poolSubscribers []PoolSubscriber
}

func (self *HardwareServerPool) AddPoolSubscriber(subscriber PoolSubscriber) {
	for _, s := range self.poolSubscribers {
		if s == subscriber {
			return
		}
	}
	self.poolSubscribers = append(self.poolSubscribers, subscriber)
}

func (self *HardwareServerPool) RemovePoolSubscriber(subscriber PoolSubscriber) {
	for p, s := range self.poolSubscribers {
		if s == subscriber {
			self.poolSubscribers = append(self.poolSubscribers[:p], self.poolSubscribers[p+1:]...)
		}
	}
}

func (self *HardwareServerPool) GetName() string {
	return self.Name
}

func (self *HardwareServerPool) IsVps() bool {
	return false
}

func (self *HardwareServerPool) addInventoryItem(item *InventoryItem) {
	self.pool[item.Id] = item
	item.Pool = self
	for _, s := range self.poolSubscribers {
		s.InventoryItemAdd(item)
	}
}

func (self *HardwareServerPool) IsInside(item *InventoryItem) bool {
	_, ok := self.pool[item.Id]
	return ok
}

//
// we suppose that this server is not allocated
//
func (self *HardwareServerPool) removeInventoryItem(item *InventoryItem) {
	delete(self.pool, item.Id)
	item.Pool = nil
	for _, s := range self.poolSubscribers {
		s.InventoryItemRemove(item)
	}
}

func (self *HardwareServerPool) Allocate(nbcores, ramsize, disksize int32, vt bool) *InventoryItem {
	log.Debug("HardwareServerPool::Allocate(", nbcores, ",", ramsize, ",", disksize, ",", vt, ")")
	var selected *InventoryItem
	for _, v := range self.pool {
		if v.Coresallocated == 0 &&
			(v.Serverconf.VtSupport == true || v.Serverconf.VtSupport == vt) &&
			v.Serverconf.NbProcessors*v.Serverconf.NbCore >= nbcores &&
			v.Serverconf.NbSlotRam*v.Serverconf.RamSize >= ramsize &&
			v.Serverconf.NbDisks*v.Serverconf.DiskSize >= disksize {
			if selected == nil {
				selected = v
			} else { // try to find the closest
				if selected.Serverconf.NbProcessors*selected.Serverconf.NbCore >
					v.Serverconf.NbProcessors*v.Serverconf.NbCore {
					selected = v
				} else if selected.Serverconf.NbSlotRam*selected.Serverconf.RamSize >
					v.Serverconf.NbSlotRam*v.Serverconf.RamSize {
					selected = v
				} else if selected.Serverconf.NbDisks*selected.Serverconf.DiskSize >
					v.Serverconf.NbDisks*v.Serverconf.DiskSize {
					selected = v
				}
			}
		}
	}
	if selected != nil {
		selected.Coresallocated = selected.Serverconf.NbProcessors * selected.Serverconf.NbCore
		selected.Ramallocated = selected.Serverconf.NbSlotRam * selected.Serverconf.RamSize
		selected.Diskallocated = selected.Serverconf.NbDisks * selected.Serverconf.DiskSize
	}
	for _, s := range self.poolSubscribers {
		s.InventoryItemAllocate(selected)
	}
	return selected
}

func (self *HardwareServerPool) HowManyFit(nbcores, ramsize, disksize int32, vt bool) int32 {
	var howmany int32
	for _, v := range self.pool {
		if v.Coresallocated == 0 &&
			(v.Serverconf.VtSupport == true || v.Serverconf.VtSupport == vt) &&
			v.Serverconf.NbProcessors*v.Serverconf.NbCore >= nbcores &&
			v.Serverconf.NbSlotRam*v.Serverconf.RamSize >= ramsize &&
			v.Serverconf.NbDisks*v.Serverconf.DiskSize >= disksize {
			howmany++
		}
	}
	return howmany
}

func (self *HardwareServerPool) IsAllocated(item *InventoryItem) bool {
	return item.Coresallocated > 0
}

func (self *HardwareServerPool) Release(item *InventoryItem, nbcores, ramsize, disksize int32) {
	log.Debug("HardwareServerPool::Release(", item, ",", nbcores, ",", ramsize, ",", disksize, ")")
	item.Coresallocated = 0
	item.Ramallocated = 0
	item.Diskallocated = 0
	for _, s := range self.poolSubscribers {
		s.InventoryItemRelease(item)
	}
}

func NewHardwareServerPool(name string) *HardwareServerPool {
	return &HardwareServerPool{
		Name:            name,
		pool:            make(map[int32]*InventoryItem),
		poolSubscribers: make([]PoolSubscriber, 0, 0),
	}
}

type VpsServerPool struct {
	Name string
	pool map[int32]*InventoryItem
	// by default cpuoverallocation is 1.0 (and can go till 2.0)
	cpuoverallocation float64
	// by default ramoverallocation is 1.0 (and can go till 1.5)
	ramoverallocation float64
	poolSubscribers   []PoolSubscriber
}

func (self *VpsServerPool) AddPoolSubscriber(subscriber PoolSubscriber) {
	for _, s := range self.poolSubscribers {
		if s == subscriber {
			return
		}
	}
	self.poolSubscribers = append(self.poolSubscribers, subscriber)
}

func (self *VpsServerPool) RemovePoolSubscriber(subscriber PoolSubscriber) {
	for p, s := range self.poolSubscribers {
		if s == subscriber {
			self.poolSubscribers = append(self.poolSubscribers[:p], self.poolSubscribers[p+1:]...)
		}
	}
}

func (self *VpsServerPool) GetName() string {
	return self.Name
}

func (self *VpsServerPool) IsVps() bool {
	return true
}

func (self *VpsServerPool) addInventoryItem(item *InventoryItem) {
	self.pool[item.Id] = item
	item.Pool = self
	for _, s := range self.poolSubscribers {
		s.InventoryItemAdd(item)
	}
}

func (self *VpsServerPool) IsInside(item *InventoryItem) bool {
	_, ok := self.pool[item.Id]
	return ok
}

//
// we suppose that this server is not allocated
//
func (self *VpsServerPool) removeInventoryItem(item *InventoryItem) {
	delete(self.pool, item.Id)
	item.Pool = nil
	for _, s := range self.poolSubscribers {
		s.InventoryItemRemove(item)
	}
}

func (self *VpsServerPool) Allocate(nbcores, ramsize, disksize int32, vt bool) *InventoryItem {
	log.Debug("VpsServerPool::Allocate(", nbcores, ",", ramsize, ",", disksize, ",", vt, ")")
	var selected *InventoryItem
	for _, v := range self.pool {
		if (v.Serverconf.VtSupport == true) &&
			float64(v.Serverconf.NbProcessors*v.Serverconf.NbCore)*self.cpuoverallocation-float64(v.Coresallocated) >= float64(nbcores) &&
			float64(v.Serverconf.NbSlotRam*v.Serverconf.RamSize)*self.ramoverallocation-float64(v.Ramallocated) >= float64(ramsize) &&
			v.Serverconf.NbDisks*v.Serverconf.DiskSize-v.Diskallocated >= disksize {
			if selected == nil {
				selected = v
			} else { // try to find the closest
				if selected.Serverconf.NbProcessors*selected.Serverconf.NbCore >
					v.Serverconf.NbProcessors*v.Serverconf.NbCore {
					selected = v
				} else if selected.Serverconf.NbSlotRam*selected.Serverconf.RamSize >
					v.Serverconf.NbSlotRam*v.Serverconf.RamSize {
					selected = v
				} else if selected.Serverconf.NbDisks*selected.Serverconf.DiskSize >
					v.Serverconf.NbDisks*v.Serverconf.DiskSize {
					selected = v
				}
			}
		}
	}
	if selected != nil {
		selected.Coresallocated += nbcores
		selected.Ramallocated += ramsize
		selected.Diskallocated += disksize
	}
	for _, s := range self.poolSubscribers {
		s.InventoryItemAllocate(selected)
	}
	return selected
}

func (self *VpsServerPool) HowManyFit(nbcores, ramsize, disksize int32, vt bool) int32 {
	var howmany int32
	for _, v := range self.pool {
		if (v.Serverconf.VtSupport == true) &&
			float64(v.Serverconf.NbProcessors*v.Serverconf.NbCore)*self.cpuoverallocation-float64(v.Coresallocated) >= float64(nbcores) &&
			float64(v.Serverconf.NbSlotRam*v.Serverconf.RamSize)*self.ramoverallocation-float64(v.Ramallocated) >= float64(ramsize) &&
			v.Serverconf.NbDisks*v.Serverconf.DiskSize-v.Diskallocated >= disksize {
			cpuX := int32(float64(v.Serverconf.NbProcessors*v.Serverconf.NbCore)*self.cpuoverallocation - float64(v.Coresallocated)/float64(nbcores))
			ramX := int32(float64(v.Serverconf.NbSlotRam*v.Serverconf.RamSize)*self.ramoverallocation - float64(v.Ramallocated)/float64(ramsize))
			diskX := int32(v.Serverconf.NbDisks*v.Serverconf.DiskSize - v.Diskallocated/disksize)

			// how many times we can put nbcores/ramsize/disksize...
			x := cpuX
			if ramX < x {
				x = ramX
			}
			if diskX < x {
				x = diskX
			}

			howmany += x
		}
	}
	return howmany
}

func (self *VpsServerPool) IsAllocated(item *InventoryItem) bool {
	return item.Coresallocated > 0
}

func (self *VpsServerPool) Release(item *InventoryItem, nbcores, ramsize, disksize int32) {
	log.Debug("HardwareServerPool::Release(", item, ",", nbcores, ",", ramsize, ",", disksize, ")")
	item.Coresallocated -= nbcores
	item.Ramallocated -= ramsize
	item.Diskallocated -= disksize
	for _, s := range self.poolSubscribers {
		s.InventoryItemRelease(item)
	}
}

func NewVpsServerPool(name string, cpuoverallocation, ramoverallocation float64) *VpsServerPool {
	return &VpsServerPool{
		Name:              name,
		pool:              make(map[int32]*InventoryItem),
		cpuoverallocation: cpuoverallocation,
		ramoverallocation: ramoverallocation,
		poolSubscribers:   make([]PoolSubscriber, 0, 0),
	}
}

type ServerOffer struct {
	Active    bool
	Name      string
	Inventory *Inventory
	Pool      ServerPool
	Vps       bool
	Nbcores   int32
	Ramsize   int32
	Disksize  int32
	Vt        bool    // only for non vps offer
	Price     float64 // per month
	// network float64
}

func (self *ServerOffer) Allocate() *InventoryItem {
	return self.Pool.Allocate(self.Nbcores, self.Ramsize, self.Disksize, self.Vt)
}

func (self *ServerOffer) Release(item *InventoryItem) {
	self.Pool.Release(item, self.Nbcores, self.Ramsize, self.Disksize)
}

func (self *ServerOffer) Save() string {
	active := "true"
	if self.Active == false {
		active = "false"
	}
	vps := "true"
	if self.Vps == false {
		vps = "false"
	}
	vt := "true"
	if self.Vt == false {
		vt = "false"
	}
	return fmt.Sprintf(`{"active": %s, "name":"%s", "vps":%s, "nbcores": %d, "ramsize":%d, "disksize":%d, "vt":%s, "price":%f }`, active, self.Name, vps, self.Nbcores, self.Ramsize, self.Disksize, vt, self.Price)
}

//
//json: demand: {
// "spec":{
//  "app": {
//    filters: {
//      diskfilter: { mindisk: 40 // Go
//      }
//    },
//    priorities: [ {"price": 2}, {"disk": 1}, {"network":1}, {"image":1},{"captive":2} ]
//    numbers: { low: 1, high: 4} // randomly?
//    },
//  "db": ...
// },
// "beginningdate": "2000-12-01" // when this demand begins to appear
// "howoften": 40 // /par an (modulo la courbe de penetration du marché?)
// "renewalfactorperyear": 0.7 // sur 1
// }
//}
//
type ServerDemandTemplate struct {
	filters    []CriteriaFilter
	priorities []PriorityPoint
	nb         [2]int32 // low, high
}

type DemandTemplate struct {
	Specs         map[string]*ServerDemandTemplate
	Beginningdate time.Time
	Howoften      int32   // per year
	Renewalfactor float64 // per year
}

type DemandInstance struct {
	template *DemandTemplate
	nb       map[string]int32 // number of instance per specs
}

//
// this function is called every "howoften" per year
//
func (self *DemandTemplate) InstanciateDemand() *DemandInstance {
	instance := &DemandInstance{
		template: self,
		nb:       make(map[string]int32),
	}
	for appname, app := range self.Specs {
		if app.nb[1] > app.nb[0] {
			instance.nb[appname] = app.nb[0] + rand.Int31()%(app.nb[1]-app.nb[0])
		} else {
			instance.nb[appname] = app.nb[0]
		}
	}
	return instance
}

//
// we should check across the inventory of different competitors
//  and from these inventory checks across all the offers
//
func (self *DemandInstance) FindOffer(inventories []*Inventory, now time.Time) *ServerBundle {
	log.Debug("DemandInstance::FindOffer(", inventories, ",", now, ")")
	selection := make(map[*Inventory]map[string]*ServerOffer)
	for _, inventory := range inventories {
		// for a given inventory we try to create the apps
		selectedoffers := make(map[string]*ServerOffer)
		nooffer := false
		for appname, app := range self.template.Specs {
			// we filter
			offers := inventory.offers
			for _, filter := range app.filters {
				offers = filter.Filter(offers)
			}

			// we score it
			points := make(map[*ServerOffer]float64)
			for _, prio := range app.priorities {
				prio.Score(offers, &points)
			}

			// we sort
			type kv struct {
				Offer *ServerOffer
				Point float64
			}
			var ss []kv
			for k, v := range points {
				ss = append(ss, kv{k, v})
			}

			sort.Slice(ss, func(i, j int) bool {
				return ss[i].Point > ss[j].Point
			})

			// we try to to see if there is enough capacity for each offer
			var allocated []*ServerContract
			var offer *ServerOffer
			for _, kv := range ss {
				allocated = make([]*ServerContract, 0)
				nb := self.nb[appname]
				filled := true
				offer = kv.Offer

				for i := 0; i < int(nb); i++ {
					inventoryitem := offer.Allocate()
					if inventoryitem == nil {
						filled = false
						break
					}
					allocated = append(allocated, &ServerContract{
						Item:      inventoryitem,
						OfferName: kv.Offer.Name,
						Vps:       kv.Offer.Vps,
						Nbcores:   kv.Offer.Nbcores,
						Ramsize:   kv.Offer.Ramsize,
						Disksize:  kv.Offer.Disksize,
						Vt:        kv.Offer.Vt,
						Price:     kv.Offer.Price,
					})
				}
				for _, contract := range allocated {
					contract.Item.Pool.Release(contract.Item, contract.Nbcores, contract.Ramsize, contract.Disksize)
				}
				if filled == true {
					selectedoffers[appname] = kv.Offer
					break
				} // else we try the next offer
			}

			// if selectedoffers[appname] != nil => we manage to allocate the app section
			if selectedoffers[appname] == nil {
				nooffer = true
				break
			}
		}
		if nooffer == false {
			selection[inventory] = selectedoffers
		}
	}

	// pass 2: we sort serverconf by inventory score

	prioInventory := make(map[*Inventory]float64)
	for appname, app := range self.template.Specs {
		points := make(map[*ServerOffer]float64)
		offers := make([]*ServerOffer, 0)
		for _, invSelection := range selection {
			points[invSelection[appname]] = 0.0
			offers = append(offers, invSelection[appname])
		}
		for _, prio := range app.priorities {
			prio.Score(offers, &points)
		}
		for inventory, invSelection := range selection {
			prioInventory[inventory] += points[invSelection[appname]] * float64(self.nb[appname])
		}
	}
	var selectedInventory *Inventory
	var points float64
	for _, inventory := range inventories {
		if prioInventory[inventory] > points {
			selectedInventory = inventory
			points = prioInventory[inventory]
		}
	}

	// pass 3: we create the contracts

	if selectedInventory == nil {
		return nil
	}

	var allocated []*ServerContract
	for appname, _ := range self.template.Specs {
		for i := 0; i < int(self.nb[appname]); i++ {
			serveroffer := selection[selectedInventory][appname]
			allocated = append(allocated, &ServerContract{
				Item:      serveroffer.Allocate(),
				OfferName: serveroffer.Name,
				Vps:       serveroffer.Vps,
				Nbcores:   serveroffer.Nbcores,
				Ramsize:   serveroffer.Ramsize,
				Disksize:  serveroffer.Disksize,
				Vt:        serveroffer.Vt,
				Price:     serveroffer.Price,
			})
		}
	}
	return &ServerBundle{
		Contracts:   allocated,
		Renewalrate: self.template.Renewalfactor,
		Date:        now,
	}
}

type ServerContract struct {
	Item      *InventoryItem
	OfferName string
	Vps       bool
	Nbcores   int32
	Ramsize   int32
	Disksize  int32
	Vt        bool    // only for non vps offer
	Price     float64 // per month
}

type ServerBundle struct {
	Contracts   []*ServerContract
	Renewalrate float64
	Date        time.Time
}

type CriteriaFilter interface {
	Filter(offers []*ServerOffer) []*ServerOffer
}

type PriorityPoint interface {
	// point = the bigger, the better
	Score(offer []*ServerOffer, points *map[*ServerOffer]float64)
}

//
// serverDemandParsing parse the "specs" sub part of a template demand such as:
//			{
//				"filters": {
//					"diskfilter": { "mindisk": 40}
//				},
//				"priorities": {
//					"price": 2,
//					"disk": 1,
//					"network":1,
//					"image":1,
//					"captive":2
//				},
//				"numbers": { "low": 1, "high": 1}
//			}
func serverDemandParsing(json map[string]interface{}) *ServerDemandTemplate {
	template := &ServerDemandTemplate{
		filters:    make([]CriteriaFilter, 0, 0),
		priorities: make([]PriorityPoint, 0, 0),
	}
	for k, v := range json {
		switch k {
		case "filters":
			if reflect.TypeOf(v).Kind() == reflect.Map {
				if reflect.TypeOf(v).Key().Kind() == reflect.String {
					// see criteriapriority.go
					template.filters = ServerDemandParsingFilters(v.(map[string]interface{}))
				}
			}
			break
		case "priorities":
			if reflect.TypeOf(v).Kind() == reflect.Map {
				if reflect.TypeOf(v).Key().Kind() == reflect.String {
					// see criteriapriority.go
					template.priorities = ServerDemandParsingPriorities(v.(map[string]interface{}))
				}
			}
			break
		case "numbers":
			if reflect.TypeOf(v).Kind() == reflect.Map {
				if reflect.TypeOf(v).Key().Kind() == reflect.String {
					// see criteriapriority.go
					template.nb = ServerDemandParsingNumbers(v.(map[string]interface{}))
				}
			}
			break
		}
	}
	return template
}

//
// DemandParsing parse a complete demand template
// such as
//  {
//		"specs": {
//			"app": {
//				"filters": {
//					"diskfilter": { "mindisk": 40}
//				},
//				"priorities": {
//					"price": 2,
//					"disk": 1,
//					"network":1,
//					"image":1,
//					"captive":2
//				},
//				"numbers": { "low": 1, "high": 1}
//			},
//			"db": {
//				"filters": {
//					"diskfilter": { "mindisk": 40}
//				},
//				"priorities": {
//					"disk": 1,
//					"network":1,
//					"image":1
//				},
//				"numbers": { "low": 1, "high": 1}
//			}
//		},
//		"beginningdate": "1996-12-01",
//		"howoften": 40
//	}
func DemandParsing(j map[string]interface{}) *DemandTemplate {
	template := &DemandTemplate{
		Specs:         make(map[string]*ServerDemandTemplate),
		Beginningdate: time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
		Howoften:      10,
		Renewalfactor: 0.1,
	}

	if data, ok := j["specs"]; ok {
		if reflect.TypeOf(data).Kind() == reflect.Map {
			for k, v := range data.(map[string]interface{}) {
				if reflect.TypeOf(v).Kind() == reflect.Map {
					template.Specs[k] = serverDemandParsing(v.(map[string]interface{}))
				}
			}
		}
	}

	if data, ok := j["beginningdate"]; ok {
		if reflect.TypeOf(data).Kind() == reflect.String {
			var year, month, day int
			fmt.Sscanf(data.(string), "%d-%d-%d", &year, &month, &day)
			template.Beginningdate = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		}
	}

	if data, ok := j["howoften"]; ok {
		if reflect.TypeOf(data).Kind() == reflect.Float64 {
			template.Howoften = int32(data.(float64))
		}
	}

	if data, ok := j["renewalfactorperyear"]; ok {
		if reflect.TypeOf(data).Kind() == reflect.Float64 {
			template.Renewalfactor = data.(float64)
		}
	}

	return template
}
