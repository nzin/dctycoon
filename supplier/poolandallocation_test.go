package supplier

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nzin/dctycoon/accounting"
	"github.com/nzin/dctycoon/timer"
	"github.com/stretchr/testify/assert"
)

type ActorMock struct {
	inventory *Inventory
}

func (self *ActorMock) GetInventory() *Inventory {
	return self.inventory
}

func (self *ActorMock) GetLedger() *accounting.Ledger {
	return nil
}

func (self *ActorMock) GetName() string {
	return "mock"
}

func (self *ActorMock) GetLocation() *LocationType {
	return AvailableLocation["siliconvalley"]
}

func TestPool(t *testing.T) {
	gt := timer.NewGameTimer()
	inventory := NewInventory(gt)

	// we create the servers
	r100 := GetServerConfTypeByName("R100")
	assert.NotEmpty(t, r100, "pickup a R100 rack server")

	serverconf := &ServerConf{
		NbProcessors: 1,
		NbCore:       1,
		VtSupport:    false,
		NbDisks:      1,
		NbSlotRam:    1,
		DiskSize:     100,
		RamSize:      2,
		ConfType:     r100,
	}
	inventory.Cart = append(inventory.Cart, &CartItem{
		Typeitem:   PRODUCT_SERVER,
		Serverconf: serverconf,
		Unitprice:  600.0,
		Nb:         2,
	})

	// we buy them and place them into the inventory
	inventory.BuyCart(time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC))
	assert.Equal(t, 2, len(inventory.Items), "we bought 2 items")

	// we now assigned them to a pool
	var defaultPhysicalPool ServerPool
	for _, p := range inventory.GetPools() {
		if p.GetName() == "default" && p.IsVps() == false {
			defaultPhysicalPool = p
		}
	}
	assert.NotEmpty(t, defaultPhysicalPool, "we found the default physical pool")

	// we assign pool and place them
	for i, inventoryitem := range inventory.Items {
		inventory.AssignPool(inventoryitem, defaultPhysicalPool)
		assert.Equal(t, defaultPhysicalPool, inventoryitem.Pool, "correct pool assignment")

		inventoryitem.Xplaced = 1
		inventoryitem.Xplaced = 1
		inventoryitem.Zplaced = i
	}

	// we now have to create an offer
	offer1 := &ServerOffer{
		Active:    true,
		Name:      "my offer",
		Inventory: inventory,
		Pool:      defaultPhysicalPool,
		Vps:       false,
		Nbcores:   1,
		Ramsize:   1,
		Disksize:  50,
		Vt:        false,
		Price:     200.0,
	}
	inventory.AddOffer(offer1)

	// now we have to create an client demand and see if it matchs
	payload := `
		{
			"filters": {
			    "diskfilter": { "mindisk": 40}
			},
			"priorities": {
				"price": 2,
				"disk": 1, 
				"network":1, 
				"image":1, 
				"captive":2
			},
			"numbers": { "low": 1, "high": 4}
		}`
	j := make(map[string]interface{})
	err := json.Unmarshal([]byte(payload), &j)
	assert.Empty(t, err, "correct JSON payload format")
	serverspecs := serverDemandParsing(j)
	assert.Equal(t, 1, len(serverspecs.Filters), "only one filter accepted (disk) ")
	assert.Equal(t, 2, len(serverspecs.priorities), "only 2 priorities accepted (disk,price) ")

	bigpayload := `{
		"specs": {
			"app": {
				"filters": {
					"diskfilter": { "mindisk": 200}
				},
				"priorities": {
					"price": 2,
					"disk": 1, 
					"network":1, 
					"image":1, 
					"captive":2
				},
				"numbers": { "low": 1, "high": 1}
			},
			"db": {
				"filters": {
					"diskfilter": { "mindisk": 40}
				},
				"priorities": {
					"disk": 1, 
					"network":1, 
					"image":1
				},
				"numbers": { "low": 1, "high": 1}
			}
		},
		"beginningdate": "1996-12-01",
		"howoften": 40
	}`
	j2 := make(map[string]interface{})
	err = json.Unmarshal([]byte(bigpayload), &j2)
	assert.Empty(t, err, "correct JSON payload format")
	demandtemplate := DemandParsing(j2)
	assert.Equal(t, 2, len(demandtemplate.Specs), "2 servers asked ")

	demand := demandtemplate.InstanciateDemand()
	actor := &ActorMock{
		inventory: inventory,
	}
	actors := []Actor{actor}
	bundlecontracts := demand.FindOffer(actors, time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC))
	assert.Empty(t, bundlecontracts, "we didn't found servers fitting the demand")

	bigpayload2 := `{
		"specs": {
			"app": {
				"filters": {
					"diskfilter": { "mindisk": 40}
				},
				"priorities": {
					"price": 2,
					"disk": 1, 
					"network":1, 
					"image":1, 
					"captive":2
				},
				"numbers": { "low": 1, "high": 1}
			},
			"db": {
				"filters": {
					"diskfilter": { "mindisk": 40}
				},
				"priorities": {
					"disk": 1, 
					"network":1, 
					"image":1
				},
				"numbers": { "low": 1, "high": 1}
			}
		},
		"beginningdate": "1996-12-01",
		"howoften": 40
	}`
	j3 := make(map[string]interface{})
	err = json.Unmarshal([]byte(bigpayload2), &j3)
	assert.Empty(t, err, "correct JSON payload format")
	demandtemplate2 := DemandParsing(j3)
	assert.Equal(t, 2, len(demandtemplate2.Specs), "2 servers asked ")

	demand2 := demandtemplate2.InstanciateDemand()
	bundlecontracts2 := demand2.FindOffer(actors, time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC))
	assert.NotEmpty(t, bundlecontracts2, "we found servers fitting the demand")
	assert.Equal(t, 2, len(bundlecontracts2.Contracts), "2 servers allocated")

	assert.Equal(t, int32(1), bundlecontracts2.Contracts[0].Item.Coresallocated, "1st server allocated")

	bundlecontracts3 := demand2.FindOffer(actors, time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC))
	assert.Empty(t, bundlecontracts3, "no server left fitting the demand")
}
