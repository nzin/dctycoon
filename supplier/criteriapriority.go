package supplier

import(
	"sort"
)

// helper structure for PriortyPoint objects
type SortedOffer struct {
	offer *ServerOffer
	value float64
}

type SortedOffers []*SortedOffer

func (s SortedOffers) Len() int { return len(s) }
func (s SortedOffers) Swap(i,j int) { s[i],s[j] = s[j],s[i] }
func (s SortedOffers) Less(i,j int) bool { return s[i].value < s[j].value }

type CriteriaFilterDisk struct {
	Disksize int32
}

func NewFilterDisk(i map[string]interface{}) *CriteriaFilterDisk {
	filter:=&CriteriaFilterDisk {
		Disksize: 0,
	}
	if v,ok := i["mindisk"]; ok {
		filter.Disksize=int32(v.(float64))
	}
	return filter
}

func (self *CriteriaFilterDisk) Filter(offers []*ServerOffer) []*ServerOffer {
	result := make([]*ServerOffer,0)
	for _,o := range(offers) {
		if o.Disksize>=self.Disksize {
			result=append(result,o)
		}
	}
	return result
}

type PriorityPrice struct {
	weight int32
}

func NewPriorityPrice(value int32) *PriorityPrice {
	priority:=&PriorityPrice {
		weight: value,
	}
	return priority
}

func (self *PriorityPrice) Score(offer []*ServerOffer, points *map[*ServerOffer]float64) {
	s := make(SortedOffers,0)
	for _,o := range(offer) {
		s=append(s,&SortedOffer{ offer:o, value:o.Price})
	}
	sort.Sort(s)
	weight:=float64(self.weight)
	for i:=len(offer)-1 ; i>=0 ; i-- {
		(*points)[s[i].offer]=weight
		weight-=float64(self.weight)/float64(len(offer))
	}
	return
}

