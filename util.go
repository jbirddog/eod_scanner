package main

func max(a float64, b float64) float64 {
	if a > b {
		return a
	}

	return b
}

func min(a float64, b float64) float64 {
	if a < b {
		return a
	}

	return b
}

func percentage(a float64, b float64) float64 {
	if b < a {
		return ((a - b) / a) * 100.0
	}

	return -(((b - a) / b) * 100.0)
}

func runningAvg(current float64, n int, new float64) float64 {
	n64 := float64(n)
	return (current*n64 + new) / (n64 + 1.0)
}

type U8LossyLookback struct {
	values uint64
}

func (u *U8LossyLookback) Push(value float64) {
	u.values <<= 8
	u.values |= uint64(value) & 0xFF
}

func (u *U8LossyLookback) LossyValue(n int) float64 {
	return float64((u.values >> (n << 3)) & 0xFF)
}
