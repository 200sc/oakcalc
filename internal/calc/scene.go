package calc

import (
	"github.com/200sc/oakcalc/internal/components/titlebar"
	"github.com/oakmound/oak/v3/scene"
	"golang.org/x/image/colornames"
)

const SceneName = "calc"

func Scene() scene.Scene {
	return scene.Scene{
		Start: func(ctx *scene.Context) {
			titlebar.New(ctx,
				titlebar.WithColor(colornames.Darkgreen),
				titlebar.WithHeight(32),
				titlebar.WithLayers([]int{10}),
				titlebar.WithTitle("OakCalc"),
			)
		},
		// No Loop function, this is the only scene.
		// No End function, this is the only scene.
	}
}
