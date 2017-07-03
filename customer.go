package dctycoon

import (
	"time"
)

// a tenant/account will have a specific rack (x,y), rackmount (id)
// a tenant/account can belong to a customer or a virtualisation solution (VPS)
// a customer can have a VM, or a physical machine
// a VM portfolio -> take all physical server, regroup them by max capacity
//   and map it to user demand. The rest is crafted as different offers (to be
//   customized)

type Customer struct {
	id   int32
	name string
	sub  []Subscription
}

type Subscription struct {
	rackx, racky int32
	rackmountid  int32         // can I get the ref directly?
	start        time.Time     // when the subscription began (year/month/day)
	duration     time.Duration // number of days for the "lease"
	renew        int32         // number of time to renew
}
