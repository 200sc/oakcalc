package calc

import (
	"image"
	"image/color"

	findfont "github.com/flopp/go-findfont"
	"github.com/oakmound/oak/v3/render"
)

func loadFallbackFonts(size int) []*render.Font {
	fallbackFonts := []string{
		"Arial.ttf",
	}
	fallbacks := []*render.Font{}
	for _, fontname := range fallbackFonts {
		fontPath, err := findfont.Find(fontname)
		if err != nil {
			continue
		}
		fg := render.FontGenerator{
			File:  fontPath,
			Color: image.NewUniform(color.RGBA{255, 255, 255, 255}),
			FontOptions: render.FontOptions{
				Size: float64(size),
			},
		}
		fallbackFont, err := fg.Generate()
		if err != nil {
			panic(err)
		}
		fallbacks = append(fallbacks, fallbackFont)
	}
	return fallbacks
}
