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

func (u *U8LossyLookback) LossyMax() float64 {
	return u.LossyMaxN(8)
}

func (u *U8LossyLookback) LossyMaxN(n int) float64 {
	maxVal := 0.0

	for i := 0; i < n; i++ {
		if val := u.LossyValue(i); val > maxVal {
			maxVal = val
		}
	}

	return maxVal
}

func (u *U8LossyLookback) LossyMin() float64 {
	return u.LossyMinN(8)
}

func (u *U8LossyLookback) LossyMinN(n int) float64 {
	minVal := 255.0

	for i := 0; i < n; i++ {
		if val := u.LossyValue(i); val < minVal {
			minVal = val
		}
	}

	return minVal
}
