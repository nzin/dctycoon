package supplier

import (
	"fmt"
	"time"
)

const (
	REPUTATION_NB_MONTH = 36
)

// Reputation stores all positive and negative "review"
// and create a [0-1] score.
// we forget what is older than REPUTATION_NB_MONTH
// It means the score is only based on the last 36 months
type Reputation struct {
	positiveRecords       map[int32]int32
	negativeRecords       map[int32]int32
	lastConsolidatedMonth int32
}

func (self *Reputation) GetScore() float64 {
	var positive int32
	var negative int32
	for i := int32(0); i < REPUTATION_NB_MONTH; i++ {
		positive += self.positiveRecords[self.lastConsolidatedMonth-i]
		negative += self.positiveRecords[self.lastConsolidatedMonth-i]
	}
	if positive+negative == 0 {
		return 1.0
	}
	return float64(positive) / float64(positive+negative)
}

func (self *Reputation) RecordPositivePoint(time time.Time) {
	yearmonth := int32(time.Year())*12 + int32(time.Month())
	if self.lastConsolidatedMonth < yearmonth {
		self.lastConsolidatedMonth = yearmonth
	}
	self.positiveRecords[yearmonth]++
}

func (self *Reputation) RecordNegativePoint(time time.Time) {
	yearmonth := int32(time.Year())*12 + int32(time.Month())
	if self.lastConsolidatedMonth < yearmonth {
		self.lastConsolidatedMonth = yearmonth
	}
	self.positiveRecords[yearmonth]--
}

func (self *Reputation) Load(data map[string]interface{}) {
	self.positiveRecords = make(map[int32]int32)
	self.negativeRecords = make(map[int32]int32)

	self.lastConsolidatedMonth = int32(data["last"].(float64))
	positive := data["positive"].([]interface{})
	for _, p := range positive {
		record := p.(map[string]interface{})
		month := int32(record["month"].(float64))
		value := int32(record["value"].(float64))
		self.positiveRecords[month] = value
	}
	negative := data["negative"].([]interface{})
	for _, n := range negative {
		record := n.(map[string]interface{})
		month := int32(record["month"].(float64))
		value := int32(record["value"].(float64))
		self.negativeRecords[month] = value
	}
}

func (self *Reputation) Save() string {
	str := "{"
	str += fmt.Sprintf(`"last":%d,`, self.lastConsolidatedMonth)
	str += `"positive": [`
	first := true
	for m, v := range self.positiveRecords {
		if first == false {
			str += ","
		}
		first = false
		str += fmt.Sprintf(`{"month":%d, "value":%d}`, m, v)
	}
	str += `}, "negative": [`
	first = true
	for m, v := range self.negativeRecords {
		if first == false {
			str += ","
		}
		first = false
		str += fmt.Sprintf(`{"month":%d, "value":%d}`, m, v)
	}

	return str + "]}"
}

func NewReputation() *Reputation {
	reputation := &Reputation{
		positiveRecords:       make(map[int32]int32),
		negativeRecords:       make(map[int32]int32),
		lastConsolidatedMonth: 0,
	}
	return reputation
}
