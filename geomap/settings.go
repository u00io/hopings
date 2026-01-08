package geomap

import "image/color"

type HighlightPoint struct {
	CountryCode string
	Style       Style
}

type Settings struct {
	Width  int
	Height int

	ShowAllCapitals bool
	AllCapitals     Style

	Highlight     []HighlightPoint
	HighlightPath []HighlightPoint

	HighlightPathLineWidth float64
	HighlightPathColor     color.Color
}

func NewSettings() *Settings {
	var c Settings
	c.Width = 1000
	c.Height = 500
	c.ShowAllCapitals = true
	c.AllCapitals = Style{
		Size: 2,
		Col:  ColorFromHex("#888888"),
	}
	c.HighlightPathLineWidth = 1
	c.HighlightPathColor = ColorFromHex("#2fe7bf")
	return &c
}
