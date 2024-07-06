package main

import (
	"math"
)

func max(x, y, z int) int {
	return int(math.Max(float64(x), math.Max(float64(y), float64(z))))
}
