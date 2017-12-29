package supplier

import (
	"reflect"
	"sort"
)

// helper structure for PriortyPoint objects
type SortedOffer struct {
	offer *ServerOffer
	value float64
}

type SortedOffers []*SortedOffer

func (s SortedOffers) Len() int           { return len(s) }
func (s SortedOffers) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SortedOffers) Less(i, j int) bool { return s[i].value < s[j].value }

type CriteriaFilterDisk struct {
	Disksize int32
}

func NewFilterDisk(i map[string]interface{}) CriteriaFilter {
	filter := &CriteriaFilterDisk{
		Disksize: 0,
	}
	if v, ok := i["mindisk"]; ok {
		filter.Disksize = int32(v.(float64))
	}
	return filter
}

func (self *CriteriaFilterDisk) Filter(offers []*ServerOffer) []*ServerOffer {
	result := make([]*ServerOffer, 0)
	for _, o := range offers {
		if o.Disksize >= self.Disksize {
			result = append(result, o)
		}
	}
	return result
}

type CriteriaFilterRam struct {
	Ramsize int32
}

func NewFilterRam(i map[string]interface{}) CriteriaFilter {
	filter := &CriteriaFilterRam{
		Ramsize: 0,
	}
	if v, ok := i["minram"]; ok {
		filter.Ramsize = int32(v.(float64))
	}
	return filter
}

func (self *CriteriaFilterRam) Filter(offers []*ServerOffer) []*ServerOffer {
	result := make([]*ServerOffer, 0)
	for _, o := range offers {
		if o.Ramsize >= self.Ramsize {
			result = append(result, o)
		}
	}
	return result
}

type CriteriaFilterNbcores struct {
	Nbcores int32
}

func NewFilterNbcores(i map[string]interface{}) CriteriaFilter {
	filter := &CriteriaFilterNbcores{
		Nbcores: 0,
	}
	if v, ok := i["mincores"]; ok {
		filter.Nbcores = int32(v.(float64))
	}
	return filter
}

func (self *CriteriaFilterNbcores) Filter(offers []*ServerOffer) []*ServerOffer {
	result := make([]*ServerOffer, 0)
	for _, o := range offers {
		if o.Nbcores >= self.Nbcores {
			result = append(result, o)
		}
	}
	return result
}

type CriteriaFilterPrice struct {
	Price float64
}

func NewFilterPrice(i map[string]interface{}) CriteriaFilter {
	filter := &CriteriaFilterPrice{
		Price: 0,
	}
	if v, ok := i["maxprice"]; ok {
		filter.Price = v.(float64)
	}
	return filter
}

func (self *CriteriaFilterPrice) Filter(offers []*ServerOffer) []*ServerOffer {
	result := make([]*ServerOffer, 0)
	for _, o := range offers {
		if o.Price <= self.Price {
			result = append(result, o)
		}
	}
	return result
}

type PriorityPrice struct {
	weight int32
}

func NewPriorityPrice(value int32) PriorityPoint {
	priority := &PriorityPrice{
		weight: value,
	}
	return priority
}

func (self *PriorityPrice) Score(offer []*ServerOffer, points *map[*ServerOffer]float64) {
	s := make(SortedOffers, 0)
	for _, o := range offer {
		s = append(s, &SortedOffer{offer: o, value: o.Price})
	}
	sort.Sort(s)
	weight := float64(self.weight)
	for i := 0; i < len(offer); i++ {
		(*points)[s[i].offer] = weight
		weight -= float64(self.weight) / float64(len(offer))
	}
	return
}

type PriorityDisk struct {
	weight int32
}

func NewPriorityDisk(value int32) PriorityPoint {
	priority := &PriorityDisk{
		weight: value,
	}
	return priority
}

func (self *PriorityDisk) Score(offer []*ServerOffer, points *map[*ServerOffer]float64) {
	s := make(SortedOffers, 0)
	for _, o := range offer {
		s = append(s, &SortedOffer{offer: o, value: float64(o.Disksize)})
	}
	sort.Sort(s)
	weight := float64(self.weight)
	for i := len(offer) - 1; i >= 0; i-- {
		(*points)[s[i].offer] = weight
		weight -= float64(self.weight) / float64(len(offer))
	}
}

type PriorityRam struct {
	weight int32
}

func NewPriorityRam(value int32) PriorityPoint {
	priority := &PriorityRam{
		weight: value,
	}
	return priority
}

func (self *PriorityRam) Score(offer []*ServerOffer, points *map[*ServerOffer]float64) {
	s := make(SortedOffers, 0)
	for _, o := range offer {
		s = append(s, &SortedOffer{offer: o, value: float64(o.Ramsize)})
	}
	sort.Sort(s)
	weight := float64(self.weight)
	for i := len(offer) - 1; i >= 0; i-- {
		(*points)[s[i].offer] = weight
		weight -= float64(self.weight) / float64(len(offer))
	}
}

type PriorityNbcores struct {
	weight int32
}

func NewPriorityNbcores(value int32) PriorityPoint {
	priority := &PriorityNbcores{
		weight: value,
	}
	return priority
}

func (self *PriorityNbcores) Score(offer []*ServerOffer, points *map[*ServerOffer]float64) {
	s := make(SortedOffers, 0)
	for _, o := range offer {
		s = append(s, &SortedOffer{offer: o, value: float64(o.Nbcores)})
	}
	sort.Sort(s)
	weight := float64(self.weight)
	for i := len(offer) - 1; i >= 0; i-- {
		(*points)[s[i].offer] = weight
		weight -= float64(self.weight) / float64(len(offer))
	}
}

func ServerDemandParsingNumbers(m map[string]interface{}) [2]int32 {
	var numbers [2]int32
	if v, ok := m["low"]; ok {
		numbers[0] = int32(v.(float64))
	}
	if v, ok := m["high"]; ok {
		numbers[1] = int32(v.(float64))
	}
	return numbers
}

var serverdemandfilterlist = map[string](func(map[string]interface{}) CriteriaFilter){
	"diskfilter":   NewFilterDisk,
	"ramfilter":    NewFilterRam,
	"nbcorefilter": NewFilterNbcores,
	"pricefilter":  NewFilterPrice,
}

func ServerDemandParsingFilters(m map[string]interface{}) []CriteriaFilter {
	filters := make([]CriteriaFilter, 0, 0)
	for filtername, function := range serverdemandfilterlist {
		if v, ok := m[filtername]; ok {
			if reflect.TypeOf(v).Kind() == reflect.Map {
				filters = append(filters, function(v.(map[string]interface{})))
			}
		}
	}
	return filters
}

var serverdemandprioritylist = map[string](func(int32) PriorityPoint){
	"disk":   NewPriorityDisk,
	"ram":    NewPriorityDisk,
	"nbcore": NewPriorityNbcores,
	"price":  NewPriorityPrice,
}

func ServerDemandParsingPriorities(m map[string]interface{}) []PriorityPoint {
	priorities := make([]PriorityPoint, 0, 0)
	for priorityname, function := range serverdemandprioritylist {
		if v, ok := m[priorityname]; ok {
			if reflect.TypeOf(v).Kind() == reflect.Float64 {
				priorities = append(priorities, function(int32(v.(float64))))
			}
		}
	}
	return priorities
}
