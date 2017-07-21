package supplier

import(
	"fmt"
	"time"
//	"github.com/nzin/dctycoon/accounting"
	"github.com/nzin/dctycoon/timer"
)

const(
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
	Id           int32
	Typeitem     int32
	Serverconf   *ServerConf // if it is a PRODUCT_SERVER
	Buydate      time.Time // for the amortizement
	Deliverydate time.Time // to know when to show it
	Xplaced,Yplaced int32 // -1 if not placed (yet)
	Zplaced      int32 //only for racking servers
	
	//allocation
	Coresallocated int32
	Ramallocated   int32 // in Mo
	Diskallocated  int32 // in Mo
}

func (self *InventoryItem) ShortDescription() string {
	ramText:=fmt.Sprintf("%d Mo",self.Serverconf.NbSlotRam*self.Serverconf.RamSize)
	if (self.Serverconf.NbSlotRam*self.Serverconf.RamSize>=2048) {
		ramText=fmt.Sprintf("%d Go",self.Serverconf.NbSlotRam*self.Serverconf.RamSize/1024)
	}
        diskText:=fmt.Sprintf("%d Mo",self.Serverconf.NbDisks*self.Serverconf.DiskSize)
        if self.Serverconf.NbDisks*self.Serverconf.DiskSize>4096 {
                diskText=fmt.Sprintf("%d Go",self.Serverconf.NbDisks*self.Serverconf.DiskSize/1024)
        }
        if self.Serverconf.NbDisks*self.Serverconf.DiskSize>4*1024*1024 {
                diskText=fmt.Sprintf("%d To",self.Serverconf.NbDisks*self.Serverconf.DiskSize/(1024*1024))
        }
	
	return fmt.Sprintf("%d cores %s RAM %s disks",
		self.Serverconf.NbProcessors*self.Serverconf.NbCore,
		ramText,
		diskText)
}

type Inventory struct {
	increment int32
	Cart        []*CartItem
	Items       map[int32]*InventoryItem
	pools       []*ServerPool
	offers      []*ServerOffer
	subscribers []InventorySubscriber
}

var GlobalInventory *Inventory

func (self *Inventory) BuyCart(buydate time.Time) {
	for _,item := range(self.Cart) {
		for i:=0;i<int(item.Nb);i++ {
			inventory:=&InventoryItem{
				Id: self.increment,
				Typeitem: item.Typeitem,
				Serverconf: item.Serverconf,
				Buydate: buydate,
				Deliverydate: buydate.Add(96*time.Hour),
				Xplaced: -1,
				Yplaced: -1,
				Zplaced: -1,
			}
			self.increment++
			self.Items[inventory.Id]=inventory
			for _,sub := range(self.subscribers) {
				instocksub:=sub
				sub.ItemInTransit(inventory)
				timer.GlobalGameTimer.AddEvent(inventory.Deliverydate,func() {
					instocksub.ItemInStock(inventory)
				})
			}
		}
	}
	//self.Cart=make([]*CartItem,0) // done in CarpPageWidget.Reset()
}

func (self *Inventory) InstallItem(item *InventoryItem,x,y,z int32) bool {
	if (item.Xplaced!=-1) { return false }
	if _, ok := self.Items[item.Id]; ok {
		for _,sub := range(self.subscribers) {
			sub.ItemRemoveFromStock(item)
		}
		item.Xplaced=x
		item.Yplaced=y
		item.Zplaced=z
		for _,sub := range(self.subscribers) {
			sub.ItemInstalled(item)
		}
		return true
	}
	return false
}

func (self *Inventory) UninstallItem(item *InventoryItem) {
	for _,sub := range(self.subscribers) {
		sub.ItemUninstalled(item)
	}
	item.Xplaced=-1
	item.Yplaced=-1
	item.Zplaced=-1
	for _,sub := range(self.subscribers) {
		sub.ItemInStock(item)
	}
}

//
// to discard an item, it must not be placed
//
func (self *Inventory) DiscardItem(item *InventoryItem) bool {
	if (item.Xplaced!=-1) { return false }
	if _, ok := self.Items[item.Id]; ok { 
		for _,sub := range(self.subscribers) {
			sub.ItemRemoveFromStock(item)
		}
		delete(self.Items,item.Id)
		return true
	}
	return false
}

func (self *Inventory) Load(conf map[string]interface{}) {
	//self.increment=int32(conf["increment"].(float64))
}

func (self *Inventory) Save() string {
	str:=""
	//str+=fmt.Sprintf(`"increment":%d`,self.increment)
	str+=`"items":[`

	str+="]"
	return str
}

func (self *Inventory) AddSubscriber(subscriber InventorySubscriber) {
	self.subscribers=append(self.subscribers,subscriber)
}

func NewInventory() *Inventory {
	inventory:=&Inventory{
		increment: 0,
		Cart: make([]*CartItem,0),
		Items: make(map[int32]*InventoryItem),
		//Items: make([]*InventoryItem,0),
		pools: make([]*ServerPool,0),
		offers: make([]*ServerOffer,0),
		subscribers: make([]InventorySubscriber,0),
	}
	return inventory
}

