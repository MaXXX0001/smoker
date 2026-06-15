package infra

import (
	"math"
	"strconv"
)

// ftoa форматує координату для URL.
func ftoa(f float64) string {
	return strconv.FormatFloat(f, 'f', 5, 64)
}

// haversineKm — відстань по поверхні Землі між двома точками, км.
func haversineKm(lat1, lon1, lat2, lon2 float64) float64 {
	const earthR = 6371.0
	rad := math.Pi / 180
	dlat := (lat2 - lat1) * rad
	dlon := (lon2 - lon1) * rad
	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1*rad)*math.Cos(lat2*rad)*math.Sin(dlon/2)*math.Sin(dlon/2)
	return earthR * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}
