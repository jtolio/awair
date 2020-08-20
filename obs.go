package awair

import (
	"time"
)

type Reading struct {
	Component string  `json:"comp"`
	Value     float64 `json:"value"`
}

type Readings []Reading

func (r Readings) Lookup(component string) (float64, bool) {
	for _, c := range r {
		if c.Component == component {
			return c.Value, true
		}
	}
	return 0, false
}

type Observation struct {
	Timestamp time.Time `json:"timestamp"`
	Score     float64   `json:"score"`
	Sensors   Readings
	Indices   Readings
}
