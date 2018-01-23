package supplier

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	//	"github.com/nzin/dctycoon/accounting"
	"github.com/nzin/dctycoon/timer"
	log "github.com/sirupsen/logrus"
)

const (
	PRODUCT_SERVER     = iota
	PRODUCT_RACK       = iota
	PRODUCT_AC         = iota
	PRODUCT_GENERATOR  = iota
	PRODUCT_DECORATION = iota
)

const (
	POWERLINE_NONE = iota
	POWERLINE_10K  = iota
	POWERLINE_100K = iota
	POWERLINE_1M   = iota
	POWERLINE_10M  = iota
)

func GetKilowattPowerline(powerline int32) float64 {
	switch powerline {
	case POWERLINE_10K:
		return 10000
	case POWERLINE_100K:
		return 100000
	case POWERLINE_1M:
		return 1000000
	case POWERLINE_10M:
		return 10000000
	}
	return 0
}

type CartItem struct {
	Typeitem   int32
	Serverconf *ServerConf // if it is an PRODUCT_SERVER
	Unitprice  float64
	Nb         int32
}

type InventoryPoolSubscriber interface {
	PoolCreate(ServerPool)
	PoolRemove(ServerPool)
}

//
// The lifecycle of an InventoryItem is
// item is created -> ItemInTransit
// item arrived -> ItemInStock
// item is installed/racked -> ItemRemovedFromStock + ItemInstalled
// item is uninstall/unracked -> ItemUninstalled + ItemInStock
//
type InventorySubscriber interface {
	ItemInTransit(*InventoryItem)
	ItemInStock(*InventoryItem)
	ItemRemoveFromStock(*InventoryItem)
	ItemInstalled(*InventoryItem)
	ItemUninstalled(*InventoryItem)
	ItemChangedPool(*InventoryItem)
}

// InventoryPowerChangeSubscriber interface is used
// to know when the comsumption or the number of generators have changed
type InventoryPowerChangeSubscriber interface {
	PowerChange(time time.Time, consumed, generated, delivered, cooler float64)
}

//
// InventoryItem is used:
// - to know what we have (immobilization)
// - where it is situated (rack)
// - which customer it is linked to
//
type InventoryItem struct {
	Id               int32
	Typeitem         int32
	Serverconf       *ServerConf // if it is a PRODUCT_SERVER
	Buydate          time.Time   // for the amortizement
	Deliverydate     time.Time   // to know when to show it
	Xplaced, Yplaced int32       // -1 if not placed (yet)
	Zplaced          int32       //only for racking servers
	Pool             ServerPool

	//allocation
	Coresallocated int32
	Ramallocated   int32 // in Mo
	Diskallocated  int32 // in Mo
}

//
// HasArrived(time) is an helper method to know if a bough item arrived in the datacenter
// i.e. InnventoryItem.Deliverydate <= now
func (self *InventoryItem) HasArrived(t time.Time) bool {
	return self.Deliverydate.Before(t) || self.Deliverydate.Equal(t)
}

//
// IsPlaced() is an helper method to know if a given item is on the map (in a rack or placed)
func (self *InventoryItem) IsPlaced() bool {
	return self.Xplaced != -1
}

func (self *InventoryItem) GetSprite() string {
	switch self.Typeitem {
	case PRODUCT_SERVER:
		return self.Serverconf.ConfType.ServerSprite
	case PRODUCT_RACK:
		return "rack"
	case PRODUCT_AC:
		return "ac"
	case PRODUCT_GENERATOR:
		return "generator"
	}
	return ""
}

func (self *InventoryItem) Save() string {
	log.Debug("InventoryItem::Save()")
	str := "{"
	switch self.Typeitem {
	case PRODUCT_SERVER:
		poolservertype := "none"
		if self.Pool != nil {
			if self.Pool.IsVps() {
				poolservertype = "vps"
			} else {
				poolservertype = "hardware"
			}
		}
		str += fmt.Sprintf(`"Id": %d, "Typeitem": "SERVER", "Buydate": "%d-%d-%d", "Deliverydate": "%d-%d-%d", "Xplaced":%d, "Yplaced":%d, "Zplaced":%d, "Coresallocated": %d, "Ramallocated": %d, "Diskallocated":%d, "NbProcessors":%d, "NbCore":%d, "VtSupport": "%t", "NbDisks":%d, "NbSlotRam":%d, "DiskSize":%d, "RamSize":%d, "ConfType": "%s", "pooltype": "%s"`,
			self.Id,
			self.Buydate.Year(), self.Buydate.Month(), self.Buydate.Day(),
			self.Deliverydate.Year(), self.Deliverydate.Month(), self.Deliverydate.Day(),
			self.Xplaced, self.Yplaced, self.Zplaced,
			self.Coresallocated,
			self.Ramallocated,
			self.Diskallocated,
			self.Serverconf.NbProcessors,
			self.Serverconf.NbCore,
			self.Serverconf.VtSupport,
			self.Serverconf.NbDisks,
			self.Serverconf.NbSlotRam,
			self.Serverconf.DiskSize,
			self.Serverconf.RamSize,
			self.Serverconf.ConfType.ServerName,
			poolservertype,
		)
	case PRODUCT_RACK:
		str += fmt.Sprintf(`"Id": %d, "Typeitem": "RACK", "Buydate": "%d-%d-%d", "Deliverydate": "%d-%d-%d", "Xplaced":%d, "Yplaced":%d`,
			self.Id,
			self.Buydate.Year(), self.Buydate.Month(), self.Buydate.Day(),
			self.Deliverydate.Year(), self.Deliverydate.Month(), self.Deliverydate.Day(),
			self.Xplaced, self.Yplaced,
		)
	case PRODUCT_AC:
		str += fmt.Sprintf(`"Id": %d, "Typeitem": "AC", "Buydate": "%d-%d-%d", "Deliverydate": "%d-%d-%d", "Xplaced":%d, "Yplaced":%d`,
			self.Id,
			self.Buydate.Year(), self.Buydate.Month(), self.Buydate.Day(),
			self.Deliverydate.Year(), self.Deliverydate.Month(), self.Deliverydate.Day(),
			self.Xplaced, self.Yplaced,
		)
	case PRODUCT_GENERATOR:
		str += fmt.Sprintf(`"Id": %d, "Typeitem": "GENERATOR", "Buydate": "%d-%d-%d", "Deliverydate": "%d-%d-%d", "Xplaced":%d, "Yplaced":%d`,
			self.Id,
			self.Buydate.Year(), self.Buydate.Month(), self.Buydate.Day(),
			self.Deliverydate.Year(), self.Deliverydate.Month(), self.Deliverydate.Day(),
			self.Xplaced, self.Yplaced,
		)
	}
	return str + "}"
}

func (self *InventoryItem) ShortDescription(condensed bool) string {
	switch self.Typeitem {
	case PRODUCT_RACK:
		return "rack"
	case PRODUCT_AC:
		return "Air Conditionner"
	case PRODUCT_GENERATOR:
		return "Generator"
	case PRODUCT_SERVER:
		ramText := fmt.Sprintf("%d Mo", self.Serverconf.NbSlotRam*self.Serverconf.RamSize)
		if self.Serverconf.NbSlotRam*self.Serverconf.RamSize >= 2048 {
			ramText = fmt.Sprintf("%d Go", self.Serverconf.NbSlotRam*self.Serverconf.RamSize/1024)
		}
		diskText := fmt.Sprintf("%d Mo", self.Serverconf.NbDisks*self.Serverconf.DiskSize)
		if self.Serverconf.NbDisks*self.Serverconf.DiskSize > 4096 {
			diskText = fmt.Sprintf("%d Go", self.Serverconf.NbDisks*self.Serverconf.DiskSize/1024)
		}
		if self.Serverconf.NbDisks*self.Serverconf.DiskSize > 4*1024*1024 {
			diskText = fmt.Sprintf("%d To", self.Serverconf.NbDisks*self.Serverconf.DiskSize/(1024*1024))
		}

		if condensed == false {
			return fmt.Sprintf("%d cores %s RAM %s disks",
				self.Serverconf.NbProcessors*self.Serverconf.NbCore,
				ramText,
				diskText)
		} else {
			return fmt.Sprintf("%d cores/%s/%s",
				self.Serverconf.NbProcessors*self.Serverconf.NbCore,
				ramText,
				diskText)

		}
	}
	return "undefined"
}

//
// Inventory structure owns the inventory of
// - InventoryItems: servers, AC, rack, generator
// - pools
// - offers
//
type Inventory struct {
	globaltimer              *timer.GameTimer
	increment                int32
	Cart                     []*CartItem
	Items                    map[int32]*InventoryItem
	pools                    []ServerPool
	offers                   []*ServerOffer
	powerline                [3]int32
	currentMaxPower          int32 // currentMaxPower is the current highest current provided by utility power lines
	consumptionHotspot       map[int32]map[int32]float64
	globalConsumption        float64
	globalGeneration         float64
	globalCooler             float64 // in BTU ~ kJ
	inventorysubscribers     []InventorySubscriber
	inventoryPoolSubscribers []InventoryPoolSubscriber
	powerchangeSubscribers   []InventoryPowerChangeSubscriber
	defaultPhysicalPool      ServerPool
	defaultVpsPool           ServerPool
}

// GetGlobalPower list all machines on the map and compute
// - the power machines consumes (positive number)
// - the power generator can sustain (positive number)
// normaly called by Inventory only
func (self *Inventory) ComputeGlobalPower() {
	self.consumptionHotspot = make(map[int32]map[int32]float64)
	self.globalConsumption = 0
	self.globalGeneration = 0
	self.globalCooler = 0
	for _, item := range self.Items {
		if item.IsPlaced() {
			if item.Typeitem == PRODUCT_AC {
				self.globalCooler += 50000
			}
			if item.Typeitem == PRODUCT_GENERATOR {
				self.globalGeneration += 50000
			}
			if item.Typeitem == PRODUCT_SERVER {
				itemconsumption := item.Serverconf.PowerConsumption()
				self.globalConsumption += itemconsumption
				if _, ok := self.consumptionHotspot[item.Yplaced]; ok == false {
					self.consumptionHotspot[item.Yplaced] = make(map[int32]float64)
				}
				self.consumptionHotspot[item.Yplaced][item.Xplaced] += itemconsumption
			}
		}
	}
}

// GetHotspotValue allow to get the heat map for each tile
func (self *Inventory) GetHotspotValue(y, x int32) float64 {
	if _, ok := self.consumptionHotspot[y]; ok == true {
		return self.consumptionHotspot[y][x]
	}
	return 0
}

// GetGlobalPower allow to get the current consumption and generator capacity
func (self *Inventory) GetGlobalPower() (consumption, generation, cooler float64) {
	return self.globalConsumption, self.globalGeneration, self.globalCooler
}

func (self *Inventory) triggerPowerChange() {
	self.ComputeGlobalPower()
	for _, s := range self.powerchangeSubscribers {
		s.PowerChange(self.globaltimer.CurrentTime, self.globalConsumption, self.globalGeneration, GetKilowattPowerline(self.currentMaxPower), self.globalCooler)
	}
}

func (self *Inventory) BuyCart(buydate time.Time) []*InventoryItem {
	log.Debug("Inventory::BuyCart(", buydate, ")")
	items := make([]*InventoryItem, 0, 0)
	for _, item := range self.Cart {
		for i := 0; i < int(item.Nb); i++ {
			inventoryitem := &InventoryItem{
				Id:           self.increment,
				Typeitem:     item.Typeitem,
				Serverconf:   item.Serverconf,
				Buydate:      buydate,
				Deliverydate: buydate.Add(96 * time.Hour),
				Xplaced:      -1,
				Yplaced:      -1,
				Zplaced:      -1,
			}
			self.increment++
			self.Items[inventoryitem.Id] = inventoryitem
			for _, sub := range self.inventorysubscribers {
				sub.ItemInTransit(inventoryitem)
			}
			self.globaltimer.AddEvent(inventoryitem.Deliverydate, func() {
				for _, sub := range self.inventorysubscribers {
					sub.ItemInStock(inventoryitem)
				}
			})
			items = append(items, inventoryitem)
		}
	}
	//self.Cart=make([]*CartItem,0) // done in CarpPageWidget.Reset()
	return items
}

func (self *Inventory) InstallItem(item *InventoryItem, x, y, z int32) bool {
	if item.Xplaced != -1 {
		return false
	}
	if _, ok := self.Items[item.Id]; ok {
		for _, sub := range self.inventorysubscribers {
			sub.ItemRemoveFromStock(item)
		}
		item.Xplaced = x
		item.Yplaced = y
		item.Zplaced = z
		for _, sub := range self.inventorysubscribers {
			sub.ItemInstalled(item)
		}
		self.triggerPowerChange()
		return true
	}
	return false
}

func (self *Inventory) UninstallItem(item *InventoryItem) {
	for _, sub := range self.inventorysubscribers {
		sub.ItemUninstalled(item)
	}
	item.Xplaced = -1
	item.Yplaced = -1
	item.Zplaced = -1
	for _, sub := range self.inventorysubscribers {
		sub.ItemInStock(item)
	}
	self.triggerPowerChange()
}

//
// to discard an item, it must not be placed
//
func (self *Inventory) DiscardItem(item *InventoryItem) bool {
	if item.Xplaced != -1 {
		return false
	}
	// remove from pool first
	self.AssignPool(item, nil)
	// remove from inventory
	if _, ok := self.Items[item.Id]; ok {
		for _, sub := range self.inventorysubscribers {
			sub.ItemRemoveFromStock(item)
		}
		delete(self.Items, item.Id)
		return true
	}
	return false
}

func (self *Inventory) AssignPool(item *InventoryItem, pool ServerPool) {
	if item.Pool != pool {
		if item.Pool != nil {
			item.Pool.removeInventoryItem(item)
		}
		item.Pool = pool
		if pool != nil {
			pool.addInventoryItem(item)
		}
		for _, sub := range self.inventorysubscribers {
			sub.ItemChangedPool(item)
		}
	}
}

func (self *Inventory) LoadItem(product map[string]interface{}) {
	log.Debug("Inventory::LoadItem(", product, ")")
	typeitem := product["Typeitem"].(string)
	buydate := strings.Split(product["Buydate"].(string), "-")
	buydateY, _ := strconv.Atoi(buydate[0])
	buydateM, _ := strconv.Atoi(buydate[1])
	buydateD, _ := strconv.Atoi(buydate[2])
	deliverydate := strings.Split(product["Deliverydate"].(string), "-")
	deliverydateY, _ := strconv.Atoi(deliverydate[0])
	deliverydateM, _ := strconv.Atoi(deliverydate[1])
	deliverydateD, _ := strconv.Atoi(deliverydate[2])
	item := &InventoryItem{
		Id:           int32(product["Id"].(float64)),
		Buydate:      time.Date(buydateY, time.Month(buydateM), buydateD, 0, 0, 0, 0, time.UTC),
		Deliverydate: time.Date(deliverydateY, time.Month(deliverydateM), deliverydateD, 0, 0, 0, 0, time.UTC),
		Xplaced:      int32(product["Xplaced"].(float64)),
		Yplaced:      int32(product["Yplaced"].(float64)),
	}

	switch typeitem {
	case "SERVER":
		item.Typeitem = PRODUCT_SERVER
		item.Zplaced = int32(product["Zplaced"].(float64))
		item.Coresallocated = int32(product["Coresallocated"].(float64))
		item.Ramallocated = int32(product["Ramallocated"].(float64))
		item.Diskallocated = int32(product["Diskallocated"].(float64))
		item.Serverconf = &ServerConf{
			NbProcessors: int32(product["NbProcessors"].(float64)),
			NbCore:       int32(product["NbCore"].(float64)),
			VtSupport:    product["VtSupport"].(string) == "true",
			NbDisks:      int32(product["NbDisks"].(float64)),
			NbSlotRam:    int32(product["NbSlotRam"].(float64)),
			DiskSize:     int32(product["DiskSize"].(float64)),
			RamSize:      int32(product["RamSize"].(float64)),
			ConfType:     GetServerConfTypeByName(product["ConfType"].(string)),
		}
		switch product["pooltype"] {
		case "hardware":
			self.AssignPool(item, self.GetDefaultPhysicalPool())
		case "vps":
			self.AssignPool(item, self.GetDefaultVpsPool())
		}

	case "RACK":
		item.Typeitem = PRODUCT_RACK
	case "AC":
		item.Typeitem = PRODUCT_AC
	case "GENERATOR":
		item.Typeitem = PRODUCT_GENERATOR
	}

	// now we store it
	self.Items[item.Id] = item

	for _, sub := range self.inventorysubscribers {
		sub.ItemInTransit(item)
	}
	self.globaltimer.AddEvent(item.Deliverydate, func() {
		for _, sub := range self.inventorysubscribers {
			sub.ItemInStock(item)
		}
	})
}

func (self *Inventory) LoadOffer(offer map[string]interface{}) {
	log.Debug("Inventory::LoadOffer(", offer, ")")
	vps := offer["vps"].(bool)

	var pool ServerPool
	for _, p := range self.pools {
		if p.IsVps() == vps {
			pool = p
			break
		}
	}

	nbcores := int32(offer["nbcores"].(float64))
	ramsize := int32(offer["ramsize"].(float64))
	disksize := int32(offer["disksize"].(float64))
	price, _ := offer["price"].(float64)

	o := &ServerOffer{
		Active:    offer["active"].(bool),
		Name:      offer["name"].(string),
		Inventory: self,
		Pool:      pool,
		Vps:       vps,
		Nbcores:   nbcores,
		Ramsize:   ramsize,
		Disksize:  disksize,
		Vt:        offer["vt"].(bool),
		Price:     price,
	}
	self.AddOffer(o)
}

// called by Load() for a given item to publish it to subscribers
func (self *Inventory) loadPublishItem(item *InventoryItem) {
	if item.Xplaced != -1 {
		for _, sub := range self.inventorysubscribers {
			sub.ItemInstalled(item)
		}
	} else {
		if self.globaltimer.CurrentTime.Before(item.Deliverydate) {
			for _, sub := range self.inventorysubscribers {
				instocksub := sub
				sub.ItemInTransit(item)
				self.globaltimer.AddEvent(item.Deliverydate, func() {
					instocksub.ItemInStock(item)
				})
			}
		} else {
			for _, sub := range self.inventorysubscribers {
				sub.ItemInStock(item)
			}
		}
	}
}

func (self *Inventory) LoadPowerlines(power map[string]interface{}) {
	self.SetPowerline(0, int32(power["powerline1"].(float64)))
	self.SetPowerline(1, int32(power["powerline2"].(float64)))
	self.SetPowerline(2, int32(power["powerline3"].(float64)))
}

func (self *Inventory) Load(conf map[string]interface{}) {
	log.Debug("Inventory::Load(", conf, ")")
	self.increment = int32(conf["increment"].(float64))
	self.consumptionHotspot = make(map[int32]map[int32]float64)
	self.Items = make(map[int32]*InventoryItem)
	if itemsinterface, ok := conf["items"]; ok {
		items := itemsinterface.([]interface{})
		for _, item := range items {
			self.LoadItem(item.(map[string]interface{}))
		}
	}
	if offersinterface, ok := conf["offers"]; ok {
		offers := offersinterface.([]interface{})
		for _, offer := range offers {
			self.LoadOffer(offer.(map[string]interface{}))
		}
	}

	// placed first RACK, AC, GENERATOR
	for _, item := range self.Items {
		if item.Typeitem == PRODUCT_RACK || item.Typeitem == PRODUCT_AC || item.Typeitem == PRODUCT_GENERATOR {
			self.loadPublishItem(item)
		}
	}
	// placed second SERVERS (especially rack servers!)
	for _, item := range self.Items {
		if item.Typeitem == PRODUCT_SERVER {
			self.loadPublishItem(item)
		}
	}
	if powerinterface, ok := conf["powerlines"]; ok {
		self.LoadPowerlines(powerinterface.(map[string]interface{}))
	}
	// to compute the hotspot
	self.triggerPowerChange()
}

func (self *Inventory) Save() string {
	log.Debug("Inventory::Save()")
	str := "{"
	str += fmt.Sprintf(`"increment":%d,`, self.increment)
	str += `"offers":[`
	firstitem := true
	for _, offer := range self.offers {
		if firstitem == true {
			firstitem = false
		} else {
			str += ",\n"
		}
		str += offer.Save()
	}
	str += "],"
	str += `"items":[`
	firstitem = true
	for _, item := range self.Items {
		if firstitem == true {
			firstitem = false
		} else {
			str += ",\n"
		}
		str += item.Save()
	}
	str += "],"
	str += fmt.Sprintf(`"powerlines": { "powerline1": %d, "powerline2": %d, "powerline3": %d }`, self.powerline[0], self.powerline[1], self.powerline[2])
	str += "}"
	return str
}

func (self *Inventory) AddInventorySubscriber(subscriber InventorySubscriber) {
	self.inventorysubscribers = append(self.inventorysubscribers, subscriber)
}

func (self *Inventory) RemoveInventorySubscriber(subscriber InventorySubscriber) {
	for i, s := range self.inventorysubscribers {
		if s == subscriber {
			self.inventorysubscribers = append(self.inventorysubscribers[:i], self.inventorysubscribers[i+1:]...)
			break
		}
	}
}

func (self *Inventory) AddPool(pool ServerPool) {
	self.pools = append(self.pools, pool)
	for _, s := range self.inventoryPoolSubscribers {
		s.PoolCreate(pool)
	}
}

func (self *Inventory) RemovePool(pool ServerPool) {
	for i, p := range self.pools {
		if p == pool {
			self.pools = append(self.pools[:i], self.pools[i+1:]...)
			break
		}
	}
	for _, s := range self.inventoryPoolSubscribers {
		s.PoolRemove(pool)
	}
}

func (self *Inventory) GetPools() []ServerPool {
	return self.pools
}

func (self *Inventory) AddInventoryPoolSubscriber(subscriber InventoryPoolSubscriber) {
	self.inventoryPoolSubscribers = append(self.inventoryPoolSubscribers, subscriber)
}

func (self *Inventory) AddOffer(offer *ServerOffer) {
	// check if not already present
	for _, o := range self.offers {
		if o == offer {
			return
		}
	}
	self.offers = append(self.offers, offer)
}

func (self *Inventory) RemoveOffer(offer *ServerOffer) {
	for i, o := range self.offers {
		if o == offer {
			self.offers = append(self.offers[:i], self.offers[i+1:]...)
			break
		}
	}
}

func (self *Inventory) UpdateOffer(offer *ServerOffer) {
	// nothing yet
}

func (self *Inventory) GetOffers() []*ServerOffer {
	return self.offers
}

func (self *Inventory) GetDefaultPhysicalPool() ServerPool {
	return self.defaultPhysicalPool
}

func (self *Inventory) GetDefaultVpsPool() ServerPool {
	return self.defaultVpsPool
}

func (self *Inventory) AddPowerStatSubscriber(subscriber InventoryPowerChangeSubscriber) {
	for _, s := range self.powerchangeSubscribers {
		if s == subscriber {
			return
		}
	}
	self.powerchangeSubscribers = append(self.powerchangeSubscribers, subscriber)
}

func (self *Inventory) RemovePowerChangeSubscriber(subscriber InventoryPowerChangeSubscriber) {
	for i, s := range self.powerchangeSubscribers {
		if s == subscriber {
			self.powerchangeSubscribers = append(self.powerchangeSubscribers[:i], self.powerchangeSubscribers[i+1:]...)
			break
		}
	}
}

// ChangePowerline is used to adjust one of the 3 main power line arrival
// power = [POWERLINE_NONE,POWERLINE_10K,POWERLINE_100K,POWERLINE_1M,POWERLINE_10M]
// we call subscribers systematically
func (self *Inventory) SetPowerline(index, power int32) {
	log.Debug("Inventory::SetPowerline(", index, ",", power, ")")
	if index < 0 || index > 2 {
		return
	}
	self.powerline[index] = power
	newmax := int32(POWERLINE_NONE)
	for _, pl := range self.powerline {
		if pl > newmax {
			newmax = pl
		}
	}
	if newmax != self.currentMaxPower {
		self.currentMaxPower = newmax
	}
	self.triggerPowerChange()
}

// PowerlineOutage is called everyday to see if we have an electricity outage
func (self *Inventory) GeneratePowerlineOutage(probability float64) {
	log.Debug("Inventory::GeneratePowerlineOutage(", probability, ")")
	newmax := int32(POWERLINE_NONE)
	for _, pl := range self.powerline {
		if rand.Float64() < probability {
			continue
		}
		if pl > newmax {
			newmax = pl
		}
	}
	if newmax != self.currentMaxPower {
		self.currentMaxPower = newmax
		self.triggerPowerChange()
	}
}

// GetPowerlines is used to collect the current situation
func (self *Inventory) GetPowerlines() [3]int32 {
	return self.powerline
}

func (self *Inventory) GetMonthlyPowerlinesPrice() float64 {
	price := float64(0)
	for _, line := range self.powerline {
		switch line {
		case POWERLINE_10K:
			price += 10
		case POWERLINE_100K:
			price += 80
		case POWERLINE_1M:
			price += 600
		case POWERLINE_10M:
			price += 5000
		}
	}
	return price
}

func NewInventory(globaltimer *timer.GameTimer) *Inventory {
	log.Debug("NewInventory(", globaltimer, ")")
	inventory := &Inventory{
		globaltimer:            globaltimer,
		increment:              0,
		Cart:                   make([]*CartItem, 0),
		Items:                  make(map[int32]*InventoryItem),
		pools:                  make([]ServerPool, 0),
		offers:                 make([]*ServerOffer, 0),
		inventorysubscribers:   make([]InventorySubscriber, 0),
		powerchangeSubscribers: make([]InventoryPowerChangeSubscriber, 0, 0),
		defaultPhysicalPool:    NewHardwareServerPool("default"),
		defaultVpsPool:         NewVpsServerPool("default", 1.2, 1.0),
		powerline:              [3]int32{POWERLINE_10K, POWERLINE_NONE, POWERLINE_NONE},
		currentMaxPower:        POWERLINE_10K,
		consumptionHotspot:     make(map[int32]map[int32]float64),
	}

	inventory.AddPool(inventory.defaultPhysicalPool)
	inventory.AddPool(inventory.defaultVpsPool)

	return inventory
}
