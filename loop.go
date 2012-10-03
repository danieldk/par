package par

import (
	"errors"
	"runtime"
)

type empty struct{}

type semaphore chan empty

type ParallelForLoop func(begin, end int, f func(int)) error

// Parallel for-loop that divides the work set in N chunks, where N is the
// number of CPUs. The data is linearly divided. The loop will return an
// error iff begin > end.
func ForChunked(begin, end int, f func(int)) error {
	if begin > end {
		return errors.New("Starting index should not be higher than end index.")
	}

	n := end - begin
	cpus := min(runtime.GOMAXPROCS(0), n)
	chunkSize := max(n/cpus, 1)
	sem := make(semaphore, cpus)

	for i, chunks := begin, 1; i < end; i, chunks = i+chunkSize, chunks+1 {
		if chunks == cpus {
			// Last Goroutine takes leftovers as well.
			go chunkedWorker(sem, i, end, f)
			break
		} else {
			go chunkedWorker(sem, i, i+chunkSize, f)
		}
	}

	for i := 0; i < cpus; i++ {
		<-sem
	}

	return nil
}

func chunkedWorker(sem semaphore, begin, end int, f func(int)) {
	for i := begin; i < end; i++ {
		f(i)
	}

	sem <- empty{}
}

// Parallel for-loop that divides the work set in N chunks, where N is the
// number of CPUs. The data is divided by interleaving. Use when the
// computation time will be uneven over regions of indices. The loop will
// return an error iff begin > end.
func ForInterleaved(begin, end int, f func(int)) error {
	if begin > end {
		return errors.New("Starting index should not be higher than end index.")
	}

	cpus := runtime.GOMAXPROCS(0)
	sem := make(semaphore, cpus)

	for i := 0; i < cpus; i++ {
		go interleavedWorker(sem, cpus, begin+i, end, f)
	}

	// Block until workers are done.
	for i := 0; i < cpus; i++ {
		<-sem
	}

	return nil
}

func interleavedWorker(sem semaphore, cpus, begin, end int, f func(int)) {
	for i := begin; i < end; i += cpus {
		f(i)
	}

	sem <- empty{}
}
