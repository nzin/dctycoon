package supplier

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	REPUTATION_NB_MONTH = 36
)

// ReputationSubscriber is mainly used by the MainStatsWidget
type ReputationSubscriber interface {
	NewReputationScore(t time.Time, score float64)
}

// Reputation stores all positive and negative "review"
// and create a [0-1] score.
// we forget what is older than REPUTATION_NB_MONTH
// It means the score is only based on the last 36 months
type Reputation struct {
	positiveRecords       map[int32]int32
	negativeRecords       map[int32]int32
	lastConsolidatedMonth int32
	reputationsubscribers []ReputationSubscriber
}

func (self *Reputation) AddReputationSubscriber(subscriber ReputationSubscriber) {
	for _, s := range self.reputationsubscribers {
		if s == subscriber {
			return
		}
	}
	self.reputationsubscribers = append(self.reputationsubscribers, subscriber)
}

func (self *Reputation) RemoveReputationSubscriber(subscriber ReputationSubscriber) {
	for i, s := range self.reputationsubscribers {
		if s == subscriber {
			self.reputationsubscribers = append(self.reputationsubscribers[:i], self.reputationsubscribers[i+1:]...)
			break
		}
	}
}

func (self *Reputation) GetScore() float64 {
	var positive int32
	var negative int32
	for i := int32(0); i < REPUTATION_NB_MONTH; i++ {
		positive += self.positiveRecords[self.lastConsolidatedMonth-i]
		negative += self.negativeRecords[self.lastConsolidatedMonth-i]
	}
	if positive+negative == 0 {
		return 1.0
	}
	return float64(positive) / float64(positive+negative)
}

// register a new reputation increment (positive/negative) and inform subscribers
func (self *Reputation) newReputation(time time.Time, positive bool) {
	log.Debug("Reputation::newReputation()")
	yearmonth := int32(time.Year())*12 + int32(time.Month())
	if self.lastConsolidatedMonth < yearmonth {
		self.lastConsolidatedMonth = yearmonth
	}
	if positive {
		self.positiveRecords[yearmonth]++
	} else {
		self.negativeRecords[yearmonth]++
	}

	// now we trigger the result
	score := self.GetScore()
	for _, s := range self.reputationsubscribers {
		s.NewReputationScore(time, score)
	}
}

func (self *Reputation) RecordPositivePoint(time time.Time) {
	self.newReputation(time, true)
}

func (self *Reputation) RecordNegativePoint(time time.Time) {
	self.newReputation(time, false)
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
	str += `], "negative": [`
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
		reputationsubscribers: make([]ReputationSubscriber, 0, 0),
	}
	return reputation
}
