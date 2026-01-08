package geomap

import (
	"github.com/fogleman/gg"
	"github.com/jonas-p/go-shp"
)

func DrawLand(dc *gg.Context, w float64, h float64) error {
	shape, err := shp.Open(ne_110m_land_shp_path)
	if err != nil {
		return err
	}
	defer shape.Close()

	for shape.Next() {
		_, geom := shape.Shape()

		poly, ok := geom.(*shp.Polygon)
		if !ok {
			continue
		}

		for part := 0; part < len(poly.Parts); part++ {
			start := int(poly.Parts[part])
			end := len(poly.Points)
			if part+1 < len(poly.Parts) {
				end = int(poly.Parts[part+1])
			}

			for i := start; i < end; i++ {
				p := poly.Points[i]
				x, y := project(p.Y, p.X, w, h)

				if i == start {
					dc.MoveTo(x, y)
				} else {
					dc.LineTo(x, y)
				}
			}
			dc.ClosePath()
		}
	}

	// fill land
	dc.SetRGB(0.12, 0.13, 0.15)
	dc.FillPreserve()

	// outline
	dc.SetRGB(0.25, 0.26, 0.28)
	dc.SetLineWidth(1.0)
	dc.Stroke()

	return nil
}
