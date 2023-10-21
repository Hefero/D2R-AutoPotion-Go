package utils

import (
	"github.com/Hefero/D2R-AutoPotion-Go/pkg/data"
	"math"
)

func DistanceFromPoint(from data.Position, to data.Position) int {
	first := math.Pow(float64(to.X-from.X), 2)
	second := math.Pow(float64(to.Y-from.Y), 2)

	return int(math.Sqrt(first + second))
}
