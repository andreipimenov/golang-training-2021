package main

import (
	"fmt"
	"math"
	"testing"
)

func TestCalculateSphereVolume(t *testing.T) {
	type Volume = float64
	type Radius = float64
	samples := map[Radius]Volume{
		7:    1436.76,
		3.2:  137.26,
		-3.2: 0.0,
	}

	for radius, volume := range samples {
		testName := fmt.Sprintln("radius", radius)
		t.Run(testName, func(t *testing.T) {
			calculated := calculateSphereVolume(radius)
			if math.Abs(calculated-volume) >= 0.1 {
				t.Errorf("got: %f, want: %f", calculated, volume)
			}
		})
	}
}
