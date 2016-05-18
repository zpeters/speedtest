package coords

import (
	"testing"
)

type HalversineTest struct {
	in  float64
	out float64
}

var halversinetests = []HalversineTest{
	{1.00, 0.22984884706593012},
}

func TestHalversine(t *testing.T) {
	for i, test := range halversinetests {
		output := haversine(test.in)
		if output != test.out {
			t.Errorf("#%d: Input %f; want %f, got %f", i, test.in, test.out, output)
		}
	}
}
