package utils

import "math"

const EarthRadius = 6371.0

func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	radLat1 := toRadians(lat1)
	radLat2 := toRadians(lat2)
	deltaLat := toRadians(lat2 - lat1)
	deltaLon := toRadians(lon2 - lon1)

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(radLat1)*math.Cos(radLat2)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadius * c
}

func toRadians(degree float64) float64 {
	return degree * math.Pi / 180
}
