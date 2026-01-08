package geomap

import (
	_ "embed"
	"fmt"
	"image/color"
	"os"
)

// /////////////////////////////////////////////////////////
// ---------------- Embedded Data ----------------
//
//go:embed ne_110m_land.shp
var ne_110m_land_shp []byte
var ne_110m_land_shp_path = ""

//go:embed countries_capitals.csv
var countries_capitals_csv []byte
var countries_capitals []Point

type ProjectionType string

func init() {
	// Get Temp folder of user
	ne_110m_land_shp_path = writeTempFileIfNotExists("ne_110m_land.shp", ne_110m_land_shp)
	fmt.Println("Using temp shapefile path:", ne_110m_land_shp_path)
	countries_capitals, _ = LoadCSV(countries_capitals_csv)
}
func writeTempFileIfNotExists(filename string, data []byte) string {
	tmpDir := os.TempDir()
	fullPath := tmpDir + string(os.PathSeparator) + filename
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		os.WriteFile(fullPath, data, 0644)
	}
	return fullPath
}

////////////////////////////////////////////////////////////

func DefaultMarkerStyle() Style {
	return Style{
		Size:  8,
		Col:   ColorFromHex("#1579b3"),
		Shape: PointShapeCircle,
	}
}

// ///////////////////////////////////////////////////////
// ---------------- Projection --------------------------
func project(lat, lon, w, h float64) (x, y float64) {
	return projectEquirectangular(lat, lon, w, h)
}

func projectEquirectangular(lat, lon, w, h float64) (x, y float64) {
	x = (lon + 180) / 360 * w
	y = (90 - lat) / 180 * h
	return
}

/*func projectRobinson(lat, lon, w, h float64) (x, y float64) {
	x = (lon + 180) / 360 * w
	y = h/2 - lat/180*h*0.85
	return
}

func projectMercator(lat, lon, w, h float64) (x, y float64) {
	x = (lon + 180) / 360 * w

	lat = math.Max(math.Min(lat, 85), -85)
	rad := lat * math.Pi / 180
	y = (1 - math.Log(math.Tan(rad)+1/math.Cos(rad))/math.Pi) / 2 * h
	return
}*/

//////////////////////////////////////////////////////////

func ColorFromHex(hex string) (col color.Color) {
	var r, g, b, a uint8 = 0, 0, 0, 255
	if len(hex) == 7 {
		fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	} else if len(hex) == 9 {
		fmt.Sscanf(hex, "#%02x%02x%02x%02x", &r, &g, &b, &a)
	}
	return color.RGBA{R: r, G: g, B: b, A: a}
}
