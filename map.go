package par

// Map over a slice of float64s with the given parallel for-loop.
func MapFloat64(loop ParallelForLoop, f func(float64) float64, l []float64) []float64 {
	result := make([]float64, len(l))

	loop(0, len(l), func(idx int) {
		result[idx] = f(l[idx])
	})

	return result
}

// Convenience function: use MapFloat64 with a chunking parallel for loop.
func MapFloat64Chunked(f func(float64) float64, l []float64) []float64 {
	return MapFloat64(ForChunked, f, l)
}

// Convenience function: use MapFloat64 with an interleaving parallel for loop.
func MapFloat64Interleaved(f func(float64) float64, l []float64) []float64 {
	return MapFloat64(ForInterleaved, f, l)
}

func max(l, r int) int {
	if l > r {
		return l
	}

	return r
}

func min(l, r int) int {
	if l < r {
		return l
	}

	return r
}
