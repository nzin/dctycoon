package supplier

import(
//	"fmt"
	"time"
//	"github.com/nzin/dctycoon/accounting"
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
	xplaced,yplaced int32 // -1 if not placed (yet)
	zplaced      int32 //only for racking servers
	
	//allocation
	Coresallocated int32
	Ramallocated   int32 // in Mo
	Diskallocated  int32 // in Mo
}

type Inventory struct {
	increment int32
	Cart      []*CartItem
	Items     map[int32]*InventoryItem
	pools     []*ServerPool
	offers    []*ServerOffer
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
				xplaced: -1,
				yplaced: -1,
				zplaced: -1,
			}
			self.increment++
			self.Items[inventory.Id]=inventory
		}
	}
	//self.Cart=make([]*CartItem,0) // done in CarpPageWidget.Reset()
}

//
// to discard an item, it must not be placed
//
func (self *Inventory) DiscardItem(item *InventoryItem) bool {
	if (item.xplaced!=-1) { return false }
	if _, ok := self.Items[item.Id]; ok { 
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

func CreateInventory() *Inventory {
	inventory:=&Inventory{
		increment: 0,
		Cart: make([]*CartItem,0),
		Items: make(map[int32]*InventoryItem),
		//Items: make([]*InventoryItem,0),
		pools: make([]*ServerPool,0),
		offers: make([]*ServerOffer,0),
	}
	return inventory
}

