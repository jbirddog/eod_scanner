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

func setRSIs(i *Indicators, a, b, c, d, e float64) {
	r := i.RSI.ring

	r.Value = a
	r = r.Next()

	r.Value = b
	r = r.Next()

	r.Value = c
	r = r.Next()

	r.Value = d
	r = r.Next()

	r.Value = e
	r = r.Next()

	i.RSI.Value = e
}
