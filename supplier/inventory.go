package supplier

import(
	"fmt"
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
// id used to link back a rack element to an inventory element
//
type InventoryItem struct {
	typeitem     int32
	serverconf   *ServerConf // if it is an PRODUCT_SERVER
	buydate      time.Time // for the amortizement
	deliverydate time.Time // to know when to show it
}

type Inventory struct {
	increment int32
	Cart  []*CartItem
	Items map[int32]*InventoryItem
}

var GlobalInventory *Inventory

func (self *Inventory) Load(conf map[string]interface{}) {
	self.increment=int32(conf["increment"].(float64))
}

func (self *Inventory) Save() string {
	str:=""
	str+=fmt.Sprintf(`"increment":%d`,self.increment)
	str+=`"items":[`

	str+="]"
	return str
}

func CreateInventory() *Inventory {
	inventory:=&Inventory{
		Cart: make([]*CartItem,0),
		Items: make(map[int32]*InventoryItem),
	}
	return inventory
}
