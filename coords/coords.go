package coords

import (
	"math"
)

// RadiusEarth is the radius of the Earth
const RadiusEarth = 6372.8

// Coordinate on the Earth
type Coordinate struct {
	Lat float64
	Lon float64
}

// Pos is a coordinate in radians
type Pos struct {
	φ float64 // latitude, radians
	ψ float64 // longitude, radians
}

// Great Circle
// http://rosettacode.org/wiki/Haversine_formula#Go
func haversine(θ float64) float64 {
	return .5 * (1 - math.Cos(θ))
}

// DegPos returns (radians) from lat and lon
func DegPos(lat, lon float64) Pos {
	return Pos{lat * math.Pi / 180, lon * math.Pi / 180}
}

// HsDist is the  distance from two positions using the great circle formula
func HsDist(p1, p2 Pos) float64 {
	return 2 * RadiusEarth * math.Asin(math.Sqrt(haversine(p2.φ-p1.φ)+
		math.Cos(p1.φ)*math.Cos(p2.φ)*haversine(p2.ψ-p1.ψ)))
}
