package par

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

const maxsize = 1025

func expensive(v float64) float64 {
	for i := 0; i < 100; i++ {
		v = math.J0(v)
	}

	return v
}

func cheap(v float64) float64 {
	return math.Sin(v)
}

func makeData(n uint) []float64 {
	data := make([]float64, n)
	for i := uint(0); i < n; i++ {
		data[i] = rand.Float64()
	}
	return data
}

func mapSerial(f func(float64) float64, l []float64) []float64 {
	result := make([]float64, len(l))

	for idx, val := range l {
		result[idx] = f(val)
	}

	return result
}

func TestParMapChunked(t *testing.T) {
	for i := uint(1); i < maxsize; i *= 2 {
		data := makeData(i)

		r := MapFloat64Chunked(cheap, data)
		check := mapSerial(cheap, data)

		compareVectors(t, r, check, fmt.Sprintf("ParMapChunked(%d)", i))
	}
}

func TestParMapInterleaved(t *testing.T) {
	for i := uint(1); i < maxsize; i *= 2 {
		data := makeData(i)

		r := MapFloat64Interleaved(cheap, data)
		check := mapSerial(cheap, data)

		compareVectors(t, r, check, fmt.Sprintf("ParMapInterleaved(%d)", i))
	}
}

func compareVectors(t *testing.T, candidate, check []float64, candidateName string) {
	// Sanity check
	if len(candidate) != len(check) {
		t.Errorf("len(%s) = %d, want %d", candidateName, len(candidate),
			len(check))
		return
	}

	for idx, val := range check {
		if candidate[idx] != val {
			t.Errorf("%s[%d] = %f, want %f", candidateName, idx,
				candidate[idx], val)
		}
	}
}
