package par

import (
	"fmt"
	"testing"
)

func makeUIntData(n uint) []uint {
	return make([]uint, n)
}

func verify(t *testing.T, name string, begin, step uint, data []uint, check uint) {
	for i := begin; i < uint(len(data)); i += step {
		val := data[i]

		if val != check {
			t.Errorf("%s[%d] == %d, expected %d", name, i, val, check)
		}
	}
}

func genericForTest(t *testing.T, name string, loop ParallelForLoop) {
	for i := uint(1); i < maxsize; i *= 2 {
		data := makeUIntData(i)

		loop(0, uint(len(data)), 1, func(idx uint) {
			data[idx]++
		})

		verify(t, fmt.Sprintf("%s-%d", name, i), 0, 1, data, 1)
	}
}

func genericStep3ForTest(t *testing.T, name string, loop ParallelForLoop) {
	for i := uint(1); i < maxsize; i *= 2 {
		data := makeUIntData(i)

		loop(0, uint(len(data)), 3, func(idx uint) {
			data[idx]++
		})

		verify(t, fmt.Sprintf("%s-%d", name, i), 0, 3, data, 1)
		verify(t, fmt.Sprintf("%s-%d", name, i), 1, 3, data, 0)
		verify(t, fmt.Sprintf("%s-%d", name, i), 2, 3, data, 0)
	}
}

func TestForChunked(t *testing.T) {
	genericForTest(t, "TestForChunked", ForChunked)
	genericStep3ForTest(t, "TestForChunked3", ForChunked)
}

func TestForInterleaved(t *testing.T) {
	genericForTest(t, "TestForInterleaved", ForInterleaved)
	genericStep3ForTest(t, "TestForInterleaved3", ForInterleaved)
}
