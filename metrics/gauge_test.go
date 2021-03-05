package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGaugeAdd(t *testing.T) {
	g := NewGauge()
	g.Add(38)
	g.Add(-8)
	assert.Equal(t, float64(30), g.Value())
}

func TestGaugeSet(t *testing.T) {
	g := NewGauge()
	g.Set(38)
	assert.Equal(t, float64(38), g.Value())
}
