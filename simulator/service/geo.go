package service

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
)

var sfPolygon orb.Polygon

func GenerateCoordinateInCity(city string) (float64, float64) {
	if city == "San Francisco" {
		bounds := sfPolygon.Bound()
		for {
			lat := bounds.Min.Lat() + rand.Float64()*(bounds.Max.Lat()-bounds.Min.Lat())
			lng := bounds.Min.Lon() + rand.Float64()*(bounds.Max.Lon()-bounds.Min.Lon())
			point := orb.Point{lng, lat}
			if planar.PolygonContains(sfPolygon, point) {
				return lat, lng
			}
		}
	}
	return 0, 0
}

func GenerateCoordinateWithinDistanceInCity(city string, lat, lng, distance float64) (float64, float64) {
	if city == "San Francisco" {
		for {
			angle := rand.Float64() * 2 * math.Pi
			// Assume distance is in meters, convert to degrees (approximate)
			distanceDegrees := distance / 111000 // 1 degree ≈ 111km
			newLat := lat + distanceDegrees*math.Cos(angle)
			newLng := lng + distanceDegrees*math.Sin(angle)/math.Cos(lat*math.Pi/180)
			point := orb.Point{newLng, newLat}
			if planar.PolygonContains(sfPolygon, point) {
				return newLat, newLng
			}
		}
	}
	return 0, 0
}

func GetDistanceBetweenCoordinates(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371000
	lat1Rad := lat1 * math.Pi / 180
	lng1Rad := lng1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lng2Rad := lng2 * math.Pi / 180
	dlat := lat2Rad - lat1Rad
	dlng := lng2Rad - lng1Rad
	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dlng/2)*math.Sin(dlng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c
	return distance
}

func GenerateMiddleCoordinates(startLat, startLng, endLat, endLng, intervalDistance float64) [][2]float64 {
	totalDistance := GetDistanceBetweenCoordinates(startLat, startLng, endLat, endLng)
	numPoints := int(math.Floor(totalDistance / intervalDistance))
	result := make([][2]float64, numPoints)

	for i := 0; i < numPoints; i++ {
		t := float64(i+1) * intervalDistance / totalDistance
		interpolatedLat := startLat + t*(endLat-startLat)
		interpolatedLng := startLng + t*(endLng-startLng)
		result[i] = [2]float64{interpolatedLat, interpolatedLng}
	}

	return result
}

func LoadGeoData() {
	var err error
	sfPolygon, err = loadSFPolygon()
	if err != nil {
		panic(err)
	}
}

func loadSFPolygon() (orb.Polygon, error) {
	data, err := os.ReadFile(filepath.Join("data", "san-francisco.geojson"))
	if err != nil {
		return nil, err
	}
	fc, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		return nil, err
	}
	if len(fc.Features) == 0 {
		return nil, fmt.Errorf("no features found in GeoJSON")
	}
	polygon, ok := fc.Features[0].Geometry.(orb.Polygon)
	if !ok {
		return nil, fmt.Errorf("first feature is not a polygon")
	}
	return polygon, nil
}
