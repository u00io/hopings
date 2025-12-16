package flags

import (
	"embed"
	"fmt"
	"image"
	"image/png"
	"strings"
)

//go:embed *.png
var flagsFS embed.FS

var flagsCache = make(map[string][]byte)
var flagsImageCache = make(map[string]image.Image)

func init() {
	items, _ := flagsFS.ReadDir(".")
	for _, item := range items {
		if !item.IsDir() {
			name := item.Name()
			data, err := flagsFS.ReadFile(name)
			if err == nil {
				flagsCache[name] = data
			}
		}
	}

	for iso, data := range flagsCache {
		img, err := png.Decode(strings.NewReader(string(data)))
		if err == nil {
			isoKey := strings.TrimSuffix(iso, ".png")
			flagsImageCache[isoKey] = img
		}
	}
}

func GetFlagBytes(iso string) ([]byte, error) {
	name := fmt.Sprintf("%s.png", strings.ToLower(iso))
	data, ok := flagsCache[name]
	if !ok {
		return nil, fmt.Errorf("flag not found: %s", iso)
	}
	return data, nil
}

func GetFlagImage(iso string) (image.Image, error) {
	isoKey := strings.ToLower(iso)
	img, ok := flagsImageCache[isoKey]
	if !ok {
		return nil, fmt.Errorf("flag image not found: %s", iso)
	}
	return img, nil
}
