package coords

import (
	//"math"
	"testing"
)

//const RadiusEarth = 6372.8

//type Coordinate struct {
//	Lat float64
//	Lon float64
//}

type HalversineTest struct {
	in float64
	out float64
}

var halversinetests = []HalversineTest{
	{1.00, 0.22984884706593012},
}

//type pos struct {
//	φ float64 // latitude, radians
//	ψ float64 // longitude, radians
//}



func TestHalversine(t *testing.T) {
	for i, test := range halversinetests {
		output := haversine(test.in)
		if output != test.out {
			t.Errorf("#%d: Input %f; want %f, got %f", i, test.in, test.out, output)
		}
	}
}


//func DegPos(lat, lon float64) pos {
//	return pos{lat * math.Pi / 180, lon * math.Pi / 180}
//}

//func HsDist(p1, p2 pos) float64 {
//	return 2 * RadiusEarth * math.Asin(math.Sqrt(haversine(p2.φ-p1.φ)+
//		math.Cos(p1.φ)*math.Cos(p2.φ)*haversine(p2.ψ-p1.ψ)))
//}
