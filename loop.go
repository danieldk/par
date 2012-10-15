package par

import (
	"errors"
	"runtime"
)

type empty struct{}

type semaphore chan empty

type ParallelForLoop func(begin, end, step uint, f func(uint)) error

// Parallel for-loop that divides the work set in N chunks, where N is the
// number of CPUs. The data is linearly divided. The loop will return an
// error iff begin > end.
func ForChunked(begin, end, step uint, f func(uint)) error {
	if begin > end {
		return errors.New("Starting index should not be higher than end index.")
	}

	n := end - begin
	cpus := min(uint(runtime.GOMAXPROCS(0)), n)
	chunkSize := max(n/cpus, 1)

	sem := make(semaphore, cpus)

	for i, chunks := begin, uint(1); i < end; i, chunks = i+chunkSize, chunks+1 {
		// The worker should start at the first index that is begin + (n * step)
		workerBegin := i
		for ; (workerBegin-begin)%step != 0 && workerBegin < end; workerBegin++ {
		}

		if chunks == cpus {
			// Last Goroutine takes leftovers as well.
			go chunkedWorker(sem, workerBegin, end, step, f)
			break
		} else {
			go chunkedWorker(sem, workerBegin, i+chunkSize, step, f)
		}
	}

	for i := uint(0); i < cpus; i++ {
		<-sem
	}

	return nil
}

func chunkedWorker(sem semaphore, begin, end, step uint, f func(uint)) {
	for i := begin; i < end; i += step {
		f(i)
	}

	sem <- empty{}
}

// Parallel for-loop that divides the work set in N chunks, where N is the
// number of CPUs. The data is divided by interleaving. Use when the
// computation time will be uneven over regions of indices. The loop will
// return an error iff begin > end.
func ForInterleaved(begin, end, step uint, f func(uint)) error {
	if begin > end {
		return errors.New("Starting index should not be higher than end index.")
	}

	cpus := uint(runtime.GOMAXPROCS(0))
	sem := make(semaphore, cpus)

	for i := uint(0); i < cpus; i++ {
		go interleavedWorker(sem, cpus, begin+(i*step), end, step, f)
	}

	// Block until workers are done.
	for i := uint(0); i < cpus; i++ {
		<-sem
	}

	return nil
}

func interleavedWorker(sem semaphore, cpus, begin, end, step uint, f func(uint)) {
	for i := begin; i < end; i += (cpus * step) {
		f(i)
	}

	sem <- empty{}
}
