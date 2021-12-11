package titlebar

import (
	"image"
	"image/color"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/entities/x/btn"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

type TitleBar struct {
	draggingWindow   bool
	draggingStartPos floatgeom.Point2
}

type Constructor struct {
	Color   color.Color
	Height  float64
	Layers  []int
	Buttons []Button

	Title          string
	TitleFontSize  int
	TitleXOffset   int
	TitleTextColor color.Color
}

type Button uint8

// Buttons to show on the title bar
const (
	// skip 0
	_                     = iota
	MinimizeButton Button = iota
	MaximizeButton Button = iota
	CloseButton    Button = iota
)

var DefaultConstructor = Constructor{
	Color:  color.RGBA{128, 128, 128, 255},
	Height: 32,
	Layers: []int{},
	Buttons: []Button{
		MinimizeButton,
		MaximizeButton,
		CloseButton,
	},
	TitleFontSize:  17,
	TitleXOffset:   10,
	TitleTextColor: color.RGBA{255, 255, 255, 255},
}

// New constructs a new
func New(ctx *scene.Context, opts ...Option) *TitleBar {

	c := DefaultConstructor
	for _, opt := range opts {
		c = opt(c)
	}

	screenHeight := ctx.Window.Height()
	screenWidth := ctx.Window.Width()

	font, _ := render.DefaultFont().RegenerateWith(func(fg render.FontGenerator) render.FontGenerator {
		fg.Size = float64(c.TitleFontSize)
		fg.Color = image.NewUniform(c.TitleTextColor)
		return fg
	})

	hdr := &TitleBar{}
	btn.New(
		btn.Font(font),
		btn.Text(c.Title),
		btn.TxtOff(10, c.Height/2-float64(c.TitleFontSize)/2),
		btn.Layers(c.Layers...),
		btn.Width(float64(screenWidth)),
		btn.Height(c.Height),
		btn.Color(c.Color),
		btn.Binding(mouse.PressOn, func(_ event.CID, ev interface{}) int {
			hdr.draggingWindow = true
			x, y, _ := oak.GetCursorPosition()
			hdr.draggingStartPos = floatgeom.Point2{float64(x), float64(y)}
			return 0
		}),
		btn.Binding(event.Enter, func(_ event.CID, ev interface{}) int {
			if hdr.draggingWindow {
				x, y, err := oak.GetCursorPosition()
				if err != nil {
					return 0
				}
				pt := floatgeom.Point2{float64(x), float64(y)}
				delta := pt.Sub(hdr.draggingStartPos)
				if delta == (floatgeom.Point2{}) {
					return 0
				}
				newX, newY, _ := ctx.Window.(*oak.Window).GetDesktopPosition()
				newX += delta.X()
				newY += delta.Y()
				ctx.Window.MoveWindow(int(newX), int(newY), screenWidth, screenHeight)
			}
			return 0
		}),
		btn.Binding(mouse.Release, func(_ event.CID, ev interface{}) int {
			if hdr.draggingWindow {
				hdr.draggingWindow = false
			}
			return 0
		}),
	)
	return hdr
}
