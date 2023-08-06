package main

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
