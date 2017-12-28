package supplier

import (
	"testing"
	"time"

	"github.com/nzin/dctycoon/timer"
	"github.com/stretchr/testify/assert"
)

/*
type EventPublisherServiceMock struct{}

func (e *EventPublisherServiceMock) Publish(shortdesc string, longdesc string) {
}*/

func TestInventory(t *testing.T) {
	gt := timer.NewGameTimer()
	/*
		ps := &EventPublisherServiceMock{}
		json := map[string]interface{}{
			"cpupricenoise":  make([]interface{}, 0, 0),
			"diskpricenoise": make([]interface{}, 0, 0),
			"rampricenoise":  make([]interface{}, 0, 0),
		}
		trend := TrendLoad(json, ps, gt)
	*/

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
	for _, i := range inventory.Items {
		assert.Equal(t, false, i.HasArrived(time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)), "we bought 2 items but there are not yet on site")
	}
	for _, i := range inventory.Items {
		assert.Equal(t, true, i.HasArrived(time.Date(1999, 1, 5, 0, 0, 0, 0, time.UTC)), "we bought 2 items and there are on site")
	}

	// we now assigned them to a pool
	var defaultPhysicalPool ServerPool
	for _, p := range inventory.GetPools() {
		if p.GetName() == "default" && p.IsVps() == false {
			defaultPhysicalPool = p
		}
	}
	assert.NotEmpty(t, defaultPhysicalPool, "we found the default physical pool")

	for _, i := range inventory.Items {
		inventory.AssignPool(i, defaultPhysicalPool)
		assert.Equal(t, defaultPhysicalPool, i.Pool, "correct pool assignment")
	}

	// discard items
	for _, i := range inventory.Items {
		inventory.DiscardItem(i)
		assert.Equal(t, nil, i.Pool, "discard correctly the item")
	}
	assert.Equal(t, 0, len(inventory.Items), "we removed the 2 items")
}

// todo: test if assignining to the wrong pool works or not (i.e. regarding VT ...)
