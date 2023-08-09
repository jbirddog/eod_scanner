package main

import (
	"runtime"
)

func batch[T any](elems []T) [][]T {
	batchSize := len(elems) / runtime.NumCPU()
	batches := make([][]T, (len(elems)+batchSize-1)/batchSize)
	prev := 0
	i := 0
	till := len(elems) - batchSize

	for prev < till {
		next := prev + batchSize
		batches[i] = elems[prev:next]
		prev = next
		i++
	}

	batches[i] = elems[prev:]

	return batches
}

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
