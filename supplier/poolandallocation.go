package supplier

import (
	"math/rand"
	"sort"
	"time"
)

type ServerPool interface {
	GetName() string
	AddInventoryItem(item *InventoryItem)
	IsInside(item *InventoryItem) bool
	RemoveInventoryItem(item *InventoryItem)
	Allocate(nbcores, ramsize, disksize int32, vt bool) *InventoryItem
	IsAllocated(item *InventoryItem) bool
	Release(item *InventoryItem, nbcores, ramsize, disksize int32)
	IsVps() bool
}

type HardwareServerPool struct {
	Name string
	pool map[int32]*InventoryItem
}

func (self *HardwareServerPool) GetName() string {
	return self.Name
}

func (self *HardwareServerPool) IsVps() bool {
	return false
}

func (self *HardwareServerPool) AddInventoryItem(item *InventoryItem) {
	self.pool[item.Id] = item
}

func (self *HardwareServerPool) IsInside(item *InventoryItem) bool {
	_, ok := self.pool[item.Id]
	return ok
}

//
// we suppose that this server is not allocated
//
func (self *HardwareServerPool) RemoveInventoryItem(item *InventoryItem) {
	delete(self.pool, item.Id)
}

func (self *HardwareServerPool) Allocate(nbcores, ramsize, disksize int32, vt bool) *InventoryItem {
	var selected *InventoryItem
	for _, v := range self.pool {
		if v.Coresallocated == 0 &&
			(v.Serverconf.VtSupport == true || v.Serverconf.VtSupport == vt) &&
			v.Serverconf.NbProcessors*v.Serverconf.NbCore >= nbcores &&
			v.Serverconf.NbSlotRam*v.Serverconf.RamSize >= ramsize &&
			v.Serverconf.NbDisks*v.Serverconf.DiskSize >= disksize {
			if selected == nil {
				selected = v
			} else {
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
	return selected
}

func (self *HardwareServerPool) IsAllocated(item *InventoryItem) bool {
	return item.Coresallocated > 0
}

func (self *HardwareServerPool) Release(item *InventoryItem, nbcores, ramsize, disksize int32) {
	item.Coresallocated = 0
	item.Ramallocated = 0
	item.Diskallocated = 0
}

type VpsServerPool struct {
	Name              string
	pool              map[int32]*InventoryItem
	// by default cpuoverallocation is 1.0 (and can go till 2.0)
	cpuoverallocation float64
	// by default ramoverallocation is 1.0 (and can go till 1.5)
	ramoverallocation float64
}

func (self *VpsServerPool) GetName() string {
	return self.Name
}

func (self *VpsServerPool) IsVps() bool {
	return true
}

func (self *VpsServerPool) AddInventoryItem(item *InventoryItem) {
	self.pool[item.Id] = item
}

func (self *VpsServerPool) IsInside(item *InventoryItem) bool {
	_, ok := self.pool[item.Id]
	return ok
}

//
// we suppose that this server is not allocated
//
func (self *VpsServerPool) RemoveInventoryItem(item *InventoryItem) {
	delete(self.pool, item.Id)
}

func (self *VpsServerPool) Allocate(nbcores, ramsize, disksize int32, vt bool) *InventoryItem {
	var selected *InventoryItem
	for _, v := range self.pool {
		if (v.Serverconf.VtSupport == true) &&
			float64(v.Serverconf.NbProcessors*v.Serverconf.NbCore)*self.cpuoverallocation-float64(v.Coresallocated) >= float64(nbcores) &&
			float64(v.Serverconf.NbSlotRam*v.Serverconf.RamSize)*self.ramoverallocation-float64(v.Ramallocated) >= float64(ramsize) &&
			v.Serverconf.NbDisks*v.Serverconf.DiskSize-v.Diskallocated >= disksize {
			if selected == nil {
				selected = v
			} else {
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
	return selected
}

func (self *VpsServerPool) IsAllocated(item *InventoryItem) bool {
	return item.Coresallocated > 0
}

func (self *VpsServerPool) Release(item *InventoryItem, nbcores, ramsize, disksize int32) {
	item.Coresallocated -= nbcores
	item.Ramallocated -= ramsize
	item.Diskallocated -= disksize
}

type ServerOffer struct {
	Inventory *Inventory
	Pool      ServerPool
	Vps       bool
	Nbcores   int32
	Ramsize   int32
	Disksize  int32
	Vt        bool // only for non vps offer
	Price     float64
	// network float64
}

func (self *ServerOffer) Allocate() *InventoryItem {
	return self.Pool.Allocate(self.Nbcores, self.Ramsize, self.Disksize, self.Vt)
}

func (self *ServerOffer) Release(item *InventoryItem) {
	self.Pool.Release(item, self.Nbcores, self.Ramsize, self.Disksize)
}

//
//json: demand: {
// "spec":{
//  "app": {
//    filters: {
//      diskfilter: { mindisk: 40 // Go
//      }
//    },
//    priority: [ {"price": 2}, {"disk": 1}, {"network":1}, {"image":1},{"captif":2} ]
//    number: { low: 1, high: 4} // tire au hasard?
//    },
//  "db": ...
// },
// "beginningdate": "2000-12-01" // date d'apparition de ce type de demande
// "howoften": 40 // /par an (modulo la courbe de penetration du marché?)
// "renewalfactorperyear": 0.7 // sur 1
// }
//}
//
type ServerDemandTemplate struct {
	filters  []CriteriaFilter
	priority []PriortyPoint
	nb       [2]int32 // low, high
}

type DemandTemplate struct {
	specs         map[string]ServerDemandTemplate
	beginningdate time.Time
	howoften      int32   // per year
	renewalfactor float64 // per year
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
	for appname, app := range self.specs {
		instance.nb[appname] = app.nb[0] + rand.Int31()%(app.nb[1]-app.nb[0])
	}
	return instance
}

//
// we should check across the inventory of different competitors
//  and from these inventory checks across all the offers
//
func (self *DemandInstance) FindOffer(inventories []*Inventory) []*ServerContract {
	selection := make(map[*Inventory]map[string]*ServerOffer)
	for _, inventory := range inventories {
		// for a given inventory we try to create the apps
		selectedoffers := make(map[string]*ServerOffer)
		nooffer := false
		for appname, app := range self.template.specs {
			// we filter
			offers := inventory.offers
			for _, filter := range app.filters {
				offers = filter.Filter(offers)
			}

			// we score it
			points := make(map[*ServerOffer]float64)
			for _, prio := range app.priority {
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
						Pool:     kv.Offer.Pool,
						Item:     inventoryitem,
						Nbcores:  kv.Offer.Nbcores,
						Ramsize:  kv.Offer.Ramsize,
						Disksize: kv.Offer.Disksize,
						Date:     time.Now(),
					})
				}
				for _, contract := range allocated {
					contract.Pool.Release(contract.Item, contract.Nbcores, contract.Ramsize, contract.Disksize)
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
	for appname, app := range self.template.specs {
		points := make(map[*ServerOffer]float64)
		offers := make([]*ServerOffer, 0)
		for _, invSelection := range selection {
			points[invSelection[appname]] = 0.0
			offers = append(offers, invSelection[appname])
		}
		for _, prio := range app.priority {
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
	for appname, _ := range self.template.specs {
		for i := 0; i < int(self.nb[appname]); i++ {
			serveroffer := selection[selectedInventory][appname]
			allocated = append(allocated, &ServerContract{
				Pool:     serveroffer.Pool,
				Item:     serveroffer.Allocate(),
				Nbcores:  serveroffer.Nbcores,
				Ramsize:  serveroffer.Ramsize,
				Disksize: serveroffer.Disksize,
				Date:     time.Now(),
			})
		}
	}
	return allocated
}

type ServerContract struct {
	Pool     ServerPool
	Item     *InventoryItem
	Nbcores  int32
	Ramsize  int32
	Disksize int32
	Date     time.Time
}

type CriteriaFilter interface {
	Filter(offers []*ServerOffer) []*ServerOffer
}

type PriortyPoint interface {
	Score(offer []*ServerOffer, points *map[*ServerOffer]float64)
}
