package par

import (
	"testing"
)

func makeUIntData(n uint) []uint {
	return make([]uint, n)
}

func verifyAllOne(t *testing.T, name string, data []uint) {
	for idx, val := range data {
		if val != 1 {
			t.Errorf("%s[%d] == %d, expected 1", name, idx, val)
		}
	}
}

func genericForTest(t *testing.T, name string, loop ParallelForLoop) {
	for i := uint(1); i < maxsize; i *= 2 {
		data := makeUIntData(i)

		loop(0, uint(len(data)), 1, func(idx uint) {
			data[idx]++
		})

		verifyAllOne(t, name, data)
	}
}

func TestForChunked(t *testing.T) {
	genericForTest(t, "TestForInterleaved", ForChunked)
}

func TestForInterleaved(t *testing.T) {
	genericForTest(t, "TestForInterleaved", ForInterleaved)
}
