package coords

import (
	"math"
)

const RadiusEarth = 6372.8

type Coordinate struct {
	Lat float64
	Lon float64
}

type pos struct {
    φ float64 // latitude, radians
    ψ float64 // longitude, radians
}

// Great Circle
// http://rosettacode.org/wiki/Haversine_formula#Go
func haversine(θ float64) float64 {
    return .5 * (1 - math.Cos(θ))
}

func DegPos(lat, lon float64) pos {
    return pos{lat * math.Pi / 180, lon * math.Pi / 180}
}

func HsDist(p1, p2 pos) float64 {
    return 2 * RadiusEarth * math.Asin(math.Sqrt(haversine(p2.φ-p1.φ)+
        math.Cos(p1.φ)*math.Cos(p2.φ)*haversine(p2.ψ-p1.ψ)))
}
