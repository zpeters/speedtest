package coords

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	type TestExpect struct {
		PosLat float64
		PosLon float64
		Lat    float64
		Lon    float64
	}

	tests := []TestExpect{
		{63.506144, 9.20091, 1.1083913080456418, 0.16058617367967148},
	}

	for test := range tests {
		res := DegPos(tests[test].PosLat, tests[test].PosLon)
		if (res.φ != tests[test].Lat) || (res.ψ != tests[test].Lon) {
			t.Logf("Got: %#v\n", res)
			t.Errorf("Should be: %#v %#v\n", tests[test].Lat, tests[test].Lon)
		}
	}
}

func TestHsDist(t *testing.T) {
	type TestExpect struct {
		Pos1Lat  float64
		Pos1Lon  float64
		Pos2Lat  float64
		Pos2Lon  float64
		Distance float64
	}

	tests := []TestExpect{
		{0.7102, -1.2923, 0.8527, 0.400, 7174.056241819571},
	}

	for test := range tests {
		pos1 := Pos{tests[test].Pos1Lat, tests[test].Pos1Lon}
		pos2 := Pos{tests[test].Pos2Lat, tests[test].Pos2Lon}
		expect := tests[test].Distance
		res := HsDist(pos1, pos2)
		assert.Equal(t, res, expect)
	}
}
