// Package par provides parallel for-loops
//
// Consider the following for loop:
//
//    for i := 0; i < n; i++ {
//    	doSomething(i)
//    }
//
// If this loop is embarassingly parallel, in other words if there
// are no data dependencies for different i's, we would often like
// to parallelize such loops to exploit the presence of multiple
// CPU cores. This package provides such parallel loops. E.g. the
// loop above could be replaced by:
//
//     par.ForInterleaved(0, n, 1, doSomething)
//
// The value of GOMAXPROCS is used to determine the number of Goroutines
// being used.
//
// The package also provides some wrappers to map a function f over
// a []float64. For instance:
//
//     result := MapFloat64Interleaved(f, data)
package par
