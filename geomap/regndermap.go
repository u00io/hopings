package geomap

import (
	"image"
	"strings"

	"github.com/fogleman/gg"
)

func RenderMap(settings *Settings) (image.Image, error) {
	width := settings.Width
	height := settings.Height

	dc := gg.NewContext(width, height)
	// background
	dc.SetRGB(0.06, 0.07, 0.09)
	dc.Clear()
	// land
	err := DrawLand(dc, float64(width), float64(height))
	if err != nil {
		return nil, err
	}

	// all capitals
	if settings.ShowAllCapitals {
		for _, p := range countries_capitals {
			x, y := project(p.Lat, p.Lon, float64(width), float64(height))
			drawMarker(dc, x, y, settings.AllCapitals)
		}
	}

	// highlight points
	for _, highlightPoint := range settings.Highlight {
		for _, p := range countries_capitals {
			if p.ISO2 != strings.ToUpper(highlightPoint.CountryCode) {
				continue
			}
			x, y := project(p.Lat, p.Lon, float64(width), float64(height))

			drawMarker(dc, x, y, highlightPoint.Style)
		}
	}

	// highlight paths
	var startPoint *Point
	for _, highlightPathItem := range settings.HighlightPath {
		for _, p := range countries_capitals {
			if p.ISO2 != strings.ToUpper(highlightPathItem.CountryCode) {
				continue
			}
			x, y := project(p.Lat, p.Lon, float64(width), float64(height))
			drawMarker(dc, x, y, highlightPathItem.Style)
			if startPoint != nil {
				x1, y1 := project(startPoint.Lat, startPoint.Lon, float64(width), float64(height))
				dc.SetColor(settings.HighlightPathColor)
				dc.SetLineWidth(settings.HighlightPathLineWidth)
				dc.MoveTo(x1, y1)
				dc.LineTo(x, y)
				dc.Stroke()
				startPoint = &p
			} else {
				startPoint = &p
			}
		}
	}

	return dc.Image(), nil
}
