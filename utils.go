package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"math"
)

func SimpleTestFunction() {
	fmt.Println("I called a test function!")
}

/// Vector math:
func dist(a, b pixel.Vec) float64 {
	return math.Abs(a.X-b.X) + math.Abs(a.Y-b.Y)
}

func div(v pixel.Vec, d float64) pixel.Vec {
	v.X /= d
	v.Y /= d

	return v
}