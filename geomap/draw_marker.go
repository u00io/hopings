package geomap

import "github.com/fogleman/gg"

func drawMarker(dc *gg.Context, x, y float64, s Style) {
	dc.SetColor(s.Col)

	switch s.Shape {

	case "square":
		dc.DrawRectangle(x-s.Size/2, y-s.Size/2, s.Size, s.Size)
		dc.Fill()

	case "diamond":
		dc.Push()
		dc.Translate(x, y)
		dc.Rotate(gg.Radians(45))
		dc.DrawRectangle(-s.Size/2, -s.Size/2, s.Size, s.Size)
		dc.Pop()
		dc.Fill()

	case "ring":
		dc.DrawCircle(x, y, s.Size/2)
		dc.SetLineWidth(s.Size * 0.25)
		dc.Stroke()

	default:
		dc.DrawCircle(x, y, s.Size/2)
		dc.Fill()
	}
}
