package main

import (
	"math"
)

func sameFloat(a float64, b float64) bool {
	return math.Abs(a-b) <= 1e-9
}

func sameData(a *EODData, b *EODData) bool {
	return a.Symbol == b.Symbol &&
		a.Date == b.Date &&
		sameFloat(a.Open, b.Open) &&
		sameFloat(a.High, b.High) &&
		sameFloat(a.Low, b.Low) &&
		sameFloat(a.Close, b.Close) &&
		a.Volume == b.Volume
}
