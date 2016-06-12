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

func TestDegPos(t *testing.T) {
	p := DegPos(63.506144, 9.20091)
	lat := 1.1083913080456418
	lon := 0.16058617367967148

	if (p.φ != lat) || (p.ψ != lon) {
		t.Logf("Got: %#v\n", p)
		t.Errorf("Should be: %#v %#v\n", lat, lon)
	}
}
