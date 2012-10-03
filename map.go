package par

// Simple parallel map, that starts a Goroutine for every value.
// Do not use for large vectors, the Goroutines will eat all your
// memory!
func ParMap(f func(float64) float64, l []float64) []float64 {
	n := len(l)
	result := make([]float64, n)
	sem := make(semaphore, n)

	for i := 0; i < n; i++ {
		go func(idx int) {
			result[idx] = f(l[idx])
			sem <- empty{}
		}(i)
	}

	for i := 0; i < n; i++ {
		<-sem
	}

	return result
}

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
