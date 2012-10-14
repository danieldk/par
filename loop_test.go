package par

import (
	"testing"
)

func makeBoolData(n uint) []bool {
	return make([]bool, n)
}

func verifyAllTrue(t *testing.T, name string, data []bool) {
	for idx, val := range data {
		if !val {
			t.Errorf("%s[%d] == false, expected true", name, idx)
		}
	}
}

func genericForTest(t *testing.T, name string, loop ParallelForLoop) {
	for i := uint(1); i < maxsize; i *= 2 {
		data := makeBoolData(i)

		ForInterleaved(0, uint(len(data)), func(idx uint) {
			data[idx] = true
		})

		verifyAllTrue(t, name, data)
	}
}

func TestForChunked(t *testing.T) {
	genericForTest(t, "TestForInterleaved", ForChunked)
}

func TestForInterleaved(t *testing.T) {
	genericForTest(t, "TestForInterleaved", ForInterleaved)
}
