package supplier

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReputation(t *testing.T) {
	reputation := NewReputation()

	date := time.Date(2012, 11, 4, 0, 0, 0, 0, time.UTC)

	reputation.RecordPositivePoint(date)
	reputation.RecordNegativePoint(date)
	reputation.RecordNegativePoint(date)
	reputation.RecordNegativePoint(date)

	assert.Equal(t, float64(0.675), reputation.GetScore(), "reputatiom of 0.25 expected")
}
