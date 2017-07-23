package supplier

import (
	"fmt"
	"time"
	"strings"
	"strconv"
	//	"github.com/nzin/dctycoon/accounting"
	"github.com/nzin/dctycoon/timer"
)

const (
	PRODUCT_SERVER    = iota
	PRODUCT_RACK      = iota
	PRODUCT_AC        = iota
	PRODUCT_GENERATOR = iota
)

type CartItem struct {
	Typeitem   int32
	Serverconf *ServerConf // if it is an PRODUCT_SERVER
	Unitprice  float64
	Nb         int32
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
}

//
// an InventoryItem is used:
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

	//allocation
	Coresallocated int32
	Ramallocated   int32 // in Mo
	Diskallocated  int32 // in Mo
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
	str := "{"
	switch self.Typeitem {
		case PRODUCT_SERVER:
			str+=fmt.Sprintf(`"Id": %d, "Typeitem": "SERVER", "Buydate": "%d-%d-%d", "Deliverydate": "%d-%d-%d", "Xplaced":%d, "Yplaced":%d, "Zplaced":%d, "Coresallocated": %d, "Ramallocated": %d, "Diskallocated":%d, "NbProcessors":%d, "NbCore":%d, "VtSupport": "%t", "NbDisks":%d, "NbSlotRam":%d, "DiskSize":%d, "RamSize":%d, "ConfType": "%s"`,
				self.Id,
				self.Buydate.Year(),self.Buydate.Month(),self.Buydate.Day(),
				self.Deliverydate.Year(),self.Deliverydate.Month(),self.Deliverydate.Day(),
				self.Xplaced,self.Yplaced,self.Zplaced,
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
			)
 		case PRODUCT_RACK:
			str+=fmt.Sprintf(`"Id": %d, "Typeitem": "RACK", "Buydate": "%d-%d-%d", "Deliverydate": "%d-%d-%d", "Xplaced":%d, "Yplaced":%d`,
				self.Id,
				self.Buydate.Year(),self.Buydate.Month(),self.Buydate.Day(),
				self.Deliverydate.Year(),self.Deliverydate.Month(),self.Deliverydate.Day(),
				self.Xplaced,self.Yplaced,
			)
 		case PRODUCT_AC:
			str+=fmt.Sprintf(`"Id": %d, "Typeitem": "AC", "Buydate": "%d-%d-%d", "Deliverydate": "%d-%d-%d", "Xplaced":%d, "Yplaced":%d`,
				self.Id,
				self.Buydate.Year(),self.Buydate.Month(),self.Buydate.Day(),
				self.Deliverydate.Year(),self.Deliverydate.Month(),self.Deliverydate.Day(),
				self.Xplaced,self.Yplaced,
			)
 		case PRODUCT_GENERATOR:
			str+=fmt.Sprintf(`"Id": %d, "Typeitem": "GENERATOR", "Buydate": "%d-%d-%d", "Deliverydate": "%d-%d-%d", "Xplaced":%d, "Yplaced":%d`,
				self.Id,
				self.Buydate.Year(),self.Buydate.Month(),self.Buydate.Day(),
				self.Deliverydate.Year(),self.Deliverydate.Month(),self.Deliverydate.Day(),
				self.Xplaced,self.Yplaced,
			)
	}
	return str + "}"
}

func (self *InventoryItem) ShortDescription() string {
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

	return fmt.Sprintf("%d cores %s RAM %s disks",
		self.Serverconf.NbProcessors*self.Serverconf.NbCore,
		ramText,
		diskText)
}

type Inventory struct {
	increment   int32
	Cart        []*CartItem
	Items       map[int32]*InventoryItem
	pools       []*ServerPool
	offers      []*ServerOffer
	subscribers []InventorySubscriber
}

var GlobalInventory *Inventory

func (self *Inventory) BuyCart(buydate time.Time) {
	for _, item := range self.Cart {
		for i := 0; i < int(item.Nb); i++ {
			inventory := &InventoryItem{
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
			self.Items[inventory.Id] = inventory
			for _, sub := range self.subscribers {
				instocksub := sub
				sub.ItemInTransit(inventory)
				timer.GlobalGameTimer.AddEvent(inventory.Deliverydate, func() {
					instocksub.ItemInStock(inventory)
				})
			}
		}
	}
	//self.Cart=make([]*CartItem,0) // done in CarpPageWidget.Reset()
}

func (self *Inventory) InstallItem(item *InventoryItem, x, y, z int32) bool {
	if item.Xplaced != -1 {
		return false
	}
	if _, ok := self.Items[item.Id]; ok {
		for _, sub := range self.subscribers {
			sub.ItemRemoveFromStock(item)
		}
		item.Xplaced = x
		item.Yplaced = y
		item.Zplaced = z
		for _, sub := range self.subscribers {
			sub.ItemInstalled(item)
		}
		return true
	}
	return false
}

func (self *Inventory) UninstallItem(item *InventoryItem) {
	for _, sub := range self.subscribers {
		sub.ItemUninstalled(item)
	}
	item.Xplaced = -1
	item.Yplaced = -1
	item.Zplaced = -1
	for _, sub := range self.subscribers {
		sub.ItemInStock(item)
	}
}

//
// to discard an item, it must not be placed
//
func (self *Inventory) DiscardItem(item *InventoryItem) bool {
	if item.Xplaced != -1 {
		return false
	}
	if _, ok := self.Items[item.Id]; ok {
		for _, sub := range self.subscribers {
			sub.ItemRemoveFromStock(item)
		}
		delete(self.Items, item.Id)
		return true
	}
	return false
}

func (self *Inventory) LoadItem(product map[string]interface{}) {
	typeitem := product["Typeitem"].(string)
	buydate := strings.Split(product["Buydate"].(string),"-")
	buydateY, _ := strconv.Atoi(buydate[0])
	buydateM, _ := strconv.Atoi(buydate[1])
	buydateD, _ := strconv.Atoi(buydate[2])
	deliverydate := strings.Split(product["Deliverydate"].(string),"-")
	deliverydateY,_ := strconv.Atoi(deliverydate[0])
	deliverydateM,_ := strconv.Atoi(deliverydate[1])
	deliverydateD,_ := strconv.Atoi(deliverydate[2])
	item:=&InventoryItem{
		Id: int32(product["Id"].(float64)),
		Buydate: time.Date(buydateY,time.Month(buydateM),buydateD, 0, 0, 0, 0, time.UTC),
		Deliverydate: time.Date(deliverydateY,time.Month(deliverydateM),deliverydateD, 0, 0, 0, 0, time.UTC),
		Xplaced: int32(product["Xplaced"].(float64)),
		Yplaced: int32(product["Yplaced"].(float64)),
	}
	
	switch typeitem {
		case "SERVER":
			item.Typeitem = PRODUCT_SERVER
			item.Zplaced = int32(product["Zplaced"].(float64))
			item.Coresallocated = int32(product["Coresallocated"].(float64))
			item.Ramallocated = int32(product["Ramallocated"].(float64))
			item.Diskallocated = int32(product["Diskallocated"].(float64))
			item.Serverconf = &ServerConf {
				NbProcessors: int32(product["NbProcessors"].(float64)),
				NbCore: int32(product["NbCore"].(float64)),
				VtSupport: product["VtSupport"].(string)=="true",
				NbDisks: int32(product["NbDisks"].(float64)),
				NbSlotRam: int32(product["NbSlotRam"].(float64)),
				DiskSize: int32(product["DiskSize"].(float64)),
				RamSize: int32(product["RamSize"].(float64)),
				ConfType: GetServerConfTypeByName(product["ConfType"].(string)),
			}

		case "RACK":
			item.Typeitem = PRODUCT_RACK
		case "AC":
			item.Typeitem = PRODUCT_AC
		case "GENERATOR":
			item.Typeitem = PRODUCT_GENERATOR
	}
	
	// now we store it
	self.Items[item.Id]=item
}

func (self *Inventory) LoadPublishItems() {
	// placed first RACK, AC, GENERATOR
	for _,item := range self.Items {
		if item.Typeitem == PRODUCT_RACK || item.Typeitem == PRODUCT_AC || item.Typeitem == PRODUCT_GENERATOR {
			if (item.Xplaced!=-1) {
				for _, sub := range self.subscribers {
					sub.ItemInstalled(item)
				}
			} else {
				if (timer.GlobalGameTimer.CurrentTime.Before(item.Deliverydate)) {
					for _, sub := range self.subscribers {
						instocksub:=sub
						sub.ItemInTransit(item)
						timer.GlobalGameTimer.AddEvent(item.Deliverydate, func() {
							instocksub.ItemInStock(item)
						})
					}
				} else {
					for _, sub := range self.subscribers {
						sub.ItemInStock(item)
					}
				}
			}
		}
	}
	// placed second SERVERS (especially rack servers!)
	for _,item := range self.Items {
		if item.Typeitem == PRODUCT_SERVER {
			if (item.Xplaced!=-1) {
				for _, sub := range self.subscribers {
					sub.ItemInstalled(item)
				}
			} else {
				if (timer.GlobalGameTimer.CurrentTime.Before(item.Deliverydate)) {
					for _, sub := range self.subscribers {
						instocksub:=sub
						sub.ItemInTransit(item)
						timer.GlobalGameTimer.AddEvent(item.Deliverydate, func() {
							instocksub.ItemInStock(item)
						})
					}
				} else {
					for _, sub := range self.subscribers {
						sub.ItemInStock(item)
					}
				}
			}
		}
	}
}

func (self *Inventory) Load(conf map[string]interface{}) {
	self.increment = int32(conf["increment"].(float64))
	self.Items = make(map[int32]*InventoryItem)
	items := conf["items"].([]interface{})
	for _,item := range(items) {
		self.LoadItem(item.(map[string]interface{}))
	}
	self.LoadPublishItems()
}

func (self *Inventory) Save() string {
	str := "{"
	str+=fmt.Sprintf(`"increment":%d,`,self.increment)
	str += `"items":[`
	firstitem:=true
	for _,item := range self.Items {
		if (firstitem == true) {
			firstitem = false
		} else {
			str += ",\n"
		}
		str += item.Save()
	}
	str += "]}"
	return str
}

func (self *Inventory) AddSubscriber(subscriber InventorySubscriber) {
	self.subscribers = append(self.subscribers, subscriber)
}

func NewInventory() *Inventory {
	inventory := &Inventory{
		increment: 0,
		Cart:      make([]*CartItem, 0),
		Items:     make(map[int32]*InventoryItem),
		//Items: make([]*InventoryItem,0),
		pools:       make([]*ServerPool, 0),
		offers:      make([]*ServerOffer, 0),
		subscribers: make([]InventorySubscriber, 0),
	}
	return inventory
}
