package geomap

import (
	"bytes"
	"encoding/csv"
	"image/color"
	"strconv"
)

type PointShape string

const (
	PointShapeRing    PointShape = "ring"
	PointShapeCircle  PointShape = "circle"
	PointShapeSquare  PointShape = "square"
	PointShapeDiamond PointShape = "diamond"
)

type Point struct {
	ISO2 string
	ISO3 string
	Lat  float64
	Lon  float64
}

type Style struct {
	Size  float64
	Col   color.Color
	Shape PointShape
}

func LoadCSV(bs []byte) ([]Point, error) {
	r := csv.NewReader(bytes.NewReader(bs))

	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var pts []Point
	for i, row := range rows {
		if i == 0 {
			continue
		}
		lat, _ := strconv.ParseFloat(row[2], 64)
		lon, _ := strconv.ParseFloat(row[3], 64)

		pts = append(pts, Point{
			ISO2: row[0],
			ISO3: row[1],
			Lat:  lat,
			Lon:  lon,
		})
	}
	return pts, nil
}
