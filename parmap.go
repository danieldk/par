package parmap

import (
	"runtime"
)

type empty struct{}

type semaphore chan empty

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

// Parallel map that divides the work set in N chunks, where N is the
// number of CPUs. The data is linearly divided.
func ParMapChunked(f func(float64) float64, l []float64) []float64 {
	n := len(l)
	result := make([]float64, n)
	cpus := min(runtime.NumCPU(), n)
	sem := make(semaphore, cpus)
	chunkSize := max(n/cpus, 1)

	for i := 0; i < n; i += chunkSize {
		go chunkedWorker(sem, l, result, i, chunkSize, f)
	}

	for i := 0; i < cpus; i++ {
		<-sem
	}

	return result
}

func chunkedWorker(sem semaphore, in, out []float64, idx, chunkSize int, f func(float64) float64) {
	n := len(in)

	for i := idx; i < idx+chunkSize; i++ {
		if i >= n {
			break
		}
		out[i] = f(in[i])
	}

	sem <- empty{}
}

// Parallel map that divides the work set in N chunks, where N is the
// number of CPUs. The data is divided by interleaving. Use when the
// computation time will be uneven for regions of the vector.
func ParMapInterleaved(f func(float64) float64, l []float64) []float64 {
	result := make([]float64, len(l))
	cpus := runtime.NumCPU()
	sem := make(semaphore, cpus)

	for i := 0; i < cpus; i++ {
		go interleavedWorker(sem, l, result, i, cpus, f)
	}

	// Block until workers are done.
	for i := 0; i < cpus; i++ {
		<-sem
	}

	return result
}

func interleavedWorker(sem semaphore, in, out []float64, startIdx, cpus int, f func(float64) float64) {
	n := len(in)

	for i := startIdx; i < n; i += cpus {
		out[i] = f(in[i])
	}

	sem <- empty{}
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
