package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/danieldk/par"
)

func primeSieve(n uint) []bool {
	sieve := make([]bool, n+1)
	half := n / 2

	for idx := range sieve {
		sieve[idx] = true
	}

	for i := uint(2); i <= half; i++ {
		if sieve[i] {
			for j := i + i; j <= n; j += i {
				sieve[j] = false
			}
		}
	}

	return sieve
}

func primes(n uint) []uint {
	r := make([]uint, 0)

	sieve := primeSieve(n)

	for i := uint(2); i <= n; i++ {
		if sieve[i] {
			r = append(r, i)
		}
	}

	return r
}

func primeFactors(primes []uint, n uint) []uint {
	factors := make([]uint, 0)

	for {
		if n <= 1 {
			break
		}

		d := n / primes[0]
		r := n % primes[0]

		if r == 0 {
			factors = append(factors, primes[0])
			n = d
		} else {
			if len(primes) == 1 {
				panic("not enough primes available to factorize number")
			}
			primes = primes[1:]
		}
	}

	return factors
}

func forLinear(begin, end, step uint, f func(uint)) error {
	for i := begin; i < end; i += step {
		f(i)
	}

	return nil
}

var loop = flag.String("loop", "linear", "Loop type (linear, chunked, or interleaved)")

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Printf("Usage: %s n\n\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	max64, err := strconv.ParseUint(flag.Args()[0], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	max := uint(max64)

	p := primes(uint(max))
	factors := make([][]uint, max)

	f := func(n uint) {
		factors[n] = primeFactors(p, uint(n))
	}

	if *loop == "linear" {
		forLinear(0, max, 2, f)
		forLinear(1, max, 2, f)
	} else if *loop == "chunked" {
		par.ForChunked(0, max, 2, f)
		par.ForChunked(1, max, 2, f)
	} else if *loop == "interleaved" {
		par.ForInterleaved(0, max, 2, f)
		par.ForInterleaved(1, max, 2, f)
	} else {
		log.Fatal("Unknown loop type")
	}

	// Note: comment this loop out if you want to do benchmarks.
	for n, f := range factors {
		fmt.Printf("%d, %v\n", n, f)
	}
}
