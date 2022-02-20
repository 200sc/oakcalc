package titlebar

import (
	"image"
	"image/color"
	"strconv"
	"time"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/entities/x/btn"
	"github.com/oakmound/oak/v3/entities/x/mods"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
	"github.com/oakmound/oak/v3/shape"
)

type TitleBar struct {
	lastPressAt        time.Time
	draggingStartPos   floatgeom.Point2
	draggingWindow     bool
	buttons            map[Button]btn.Btn
	startingDimensions intgeom.Point2
	maximized          bool
}

type Constructor struct {
	Color          color.Color
	HighlightColor color.Color
	MouseDownColor color.Color
	Height         float64
	Layers         []int

	Title          string
	TitleFontSize  int
	TitleXOffset   int
	TitleTextColor color.Color

	Buttons              []Button
	ButtonWidth          float64
	DoubleClickThreshold time.Duration
}

type Button uint8

// Buttons to show on the title bar
const (
	ButtonMinimize Button = iota
	ButtonMaximize Button = iota
	ButtonClose    Button = iota
)

var DefaultConstructor = Constructor{
	Color:  color.RGBA{128, 128, 128, 255},
	Height: 32,
	Layers: []int{},
	Buttons: []Button{
		ButtonMinimize,
		ButtonMaximize,
		ButtonClose,
	},
	ButtonWidth:          32,
	TitleFontSize:        17,
	TitleXOffset:         10,
	TitleTextColor:       color.RGBA{255, 255, 255, 255},
	DoubleClickThreshold: 200 * time.Millisecond,
}

// New constructs a new TitleBar
func New(ctx *scene.Context, opts ...Option) *TitleBar {

	construct := DefaultConstructor
	for _, opt := range opts {
		construct = opt(construct)
	}
	if construct.HighlightColor == nil {
		construct.HighlightColor = mods.Lighter(construct.Color, .10)
	}
	if construct.MouseDownColor == nil {
		construct.MouseDownColor = mods.Lighter(construct.Color, .20)
	}

	screenHeight := ctx.Window.Height()
	screenWidth := ctx.Window.Width()

	font, _ := render.DefaultFont().RegenerateWith(func(fg render.FontGenerator) render.FontGenerator {
		fg.Size = float64(construct.TitleFontSize)
		fg.Color = image.NewUniform(construct.TitleTextColor)
		return fg
	})

	dragBarWidth := float64(screenWidth)

	totalButtonsSize := construct.ButtonWidth * float64(len(construct.Buttons))
	dragBarWidth -= totalButtonsSize

	hdr := &TitleBar{
		lastPressAt:        time.Now(),
		buttons:            make(map[Button]btn.Btn),
		startingDimensions: intgeom.Point2{ctx.Window.Width(), ctx.Window.Height()},
	}

	for i, button := range construct.Buttons {
		i := i
		button := button
		var r render.Modifiable = render.NewColorBox(int(construct.ButtonWidth), int(construct.Height), construct.Color)
		txt := strconv.Itoa(i)
		var clickBinding = func(c event.CID, e *mouse.Event) int {
			return 0
		}

		switch button {
		case ButtonMinimize:
			r = render.NewSwitch("nohover", map[string]render.Modifiable{
				"nohover": render.SpriteFromShape(minimizeIcon, int(construct.ButtonWidth), int(construct.Height), color.RGBA{255, 255, 255, 255}, construct.Color),
				"hover":   render.SpriteFromShape(minimizeIcon, int(construct.ButtonWidth), int(construct.Height), color.RGBA{255, 255, 255, 255}, construct.HighlightColor),
				"onpress": render.SpriteFromShape(minimizeIcon, int(construct.ButtonWidth), int(construct.Height), color.RGBA{255, 255, 255, 255}, construct.MouseDownColor),
			})
			txt = ""
			clickBinding = func(c event.CID, e *mouse.Event) int {
				ctx.Window.(*oak.Window).Minimize()
				return 0
			}
		case ButtonClose:
			r = render.NewSwitch("nohover", map[string]render.Modifiable{
				"nohover": render.SpriteFromShape(closeIcon, int(construct.ButtonWidth), int(construct.Height), color.RGBA{255, 255, 255, 255}, construct.Color),
				"hover":   render.SpriteFromShape(closeIcon, int(construct.ButtonWidth), int(construct.Height), color.RGBA{255, 255, 255, 255}, construct.HighlightColor),
				"onpress": render.SpriteFromShape(closeIcon, int(construct.ButtonWidth), int(construct.Height), color.RGBA{255, 255, 255, 255}, construct.MouseDownColor),
			})
			txt = ""
			clickBinding = func(c event.CID, e *mouse.Event) int {
				ctx.Window.Quit()
				return 0
			}
		case ButtonMaximize:
			r = render.NewSwitch("nohover", map[string]render.Modifiable{
				"nohover":        render.SpriteFromShape(maximizeIcon, int(construct.ButtonWidth), int(construct.Height), color.RGBA{255, 255, 255, 255}, construct.Color),
				"hover":          render.SpriteFromShape(maximizeIcon, int(construct.ButtonWidth), int(construct.Height), color.RGBA{255, 255, 255, 255}, construct.HighlightColor),
				"onpress":        render.SpriteFromShape(maximizeIcon, int(construct.ButtonWidth), int(construct.Height), color.RGBA{255, 255, 255, 255}, construct.MouseDownColor),
				"nohover-revert": render.SpriteFromShape(normalizeIcon, int(construct.ButtonWidth), int(construct.Height), color.RGBA{255, 255, 255, 255}, construct.Color),
				"hover-revert":   render.SpriteFromShape(normalizeIcon, int(construct.ButtonWidth), int(construct.Height), color.RGBA{255, 255, 255, 255}, construct.HighlightColor),
				"onpress-revert": render.SpriteFromShape(normalizeIcon, int(construct.ButtonWidth), int(construct.Height), color.RGBA{255, 255, 255, 255}, construct.MouseDownColor),
			})
			txt = ""

			clickBinding = func(c event.CID, e *mouse.Event) int {
				b, _ := ctx.CallerMap.GetEntity(c).(btn.Btn)
				hdr.maximized = toggleMaximize(ctx, b)
				return 0
			}
		}
		hdr.buttons[button] = btn.New(
			btn.Text(txt),
			btn.Pos(dragBarWidth+float64(i)*construct.ButtonWidth, 0),
			btn.Renderable(r),
			btn.Height(construct.Height),
			btn.Width(construct.ButtonWidth),
			btn.Layers(construct.Layers...),
			btn.Binding(mouse.Start, mouse.Binding(func(c event.CID, e *mouse.Event) int {
				b, _ := ctx.CallerMap.GetEntity(c).(btn.Btn)
				if sw, ok := b.GetRenderable().(*render.Switch); ok {
					suffix, _ := b.Metadata("switch-suffix")
					sw.Set("hover" + suffix)
				}
				return 0
			})),
			btn.Binding(mouse.Stop, mouse.Binding(func(c event.CID, e *mouse.Event) int {
				b, _ := ctx.CallerMap.GetEntity(c).(btn.Btn)
				if sw, ok := b.GetRenderable().(*render.Switch); ok {
					suffix, _ := b.Metadata("switch-suffix")
					sw.Set("nohover" + suffix)
				}
				return 0
			})),
			btn.Binding(mouse.PressOn, mouse.Binding(func(c event.CID, e *mouse.Event) int {
				b, _ := ctx.CallerMap.GetEntity(c).(btn.Btn)
				if sw, ok := b.GetRenderable().(*render.Switch); ok {
					suffix, _ := b.Metadata("switch-suffix")
					sw.Set("onpress" + suffix)
				}
				return 0
			})),
			btn.Click(mouse.Binding(clickBinding)),
			btn.Binding(oak.WindowSizeChange, oak.SizeChangeEvent(func(c event.CID, pt intgeom.Point2) int {
				b, _ := ctx.CallerMap.GetEntity(c).(*btn.Box)
				b.SetPos(float64(pt.X())-totalButtonsSize+float64(i)*construct.ButtonWidth, 0)
				return 0
			})),
		)
	}

	btn.New(
		btn.Font(font),
		btn.Text(construct.Title),
		btn.TxtOff(10, construct.Height/2-float64(construct.TitleFontSize)/2),
		btn.Layers(construct.Layers...),
		btn.Width(dragBarWidth),
		btn.Height(construct.Height),
		btn.Color(construct.Color),
		btn.Binding(mouse.PressOn, func(_ event.CID, ev interface{}) int {
			if time.Since(hdr.lastPressAt) < construct.DoubleClickThreshold {
				if mxbtn, ok := hdr.buttons[ButtonMaximize]; ok {
					hdr.maximized = toggleMaximize(ctx, mxbtn)
				}
				// if this is not set, dragging can persist after the window shrinks
				hdr.draggingWindow = false
				return 0
			}
			hdr.lastPressAt = time.Now()
			hdr.draggingWindow = true
			x, y, _ := oak.GetCursorPosition()
			hdr.draggingStartPos = floatgeom.Point2{float64(x), float64(y)}
			return 0
		}),
		// Q: Why not mouse.Drag?
		// A: mouse.Drag is only triggered for on-screen mouse events. If the mouse
		//    falls out of the window, as it likely will if you drag the window up,
		//    the window will freeze until you bring the mouse cursor back into the window.
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
				if hdr.maximized {
					if mxbtn, ok := hdr.buttons[ButtonMaximize]; ok {
						hdr.maximized = toggleMaximize(ctx, mxbtn)
					}
				}
				ctx.Window.MoveWindow(int(newX), int(newY), screenWidth, screenHeight)
				if !floatgeom.NewRect2WH(0, 0, float64(screenWidth), float64(screenHeight)).Contains(hdr.draggingStartPos) {
					hdr.draggingStartPos = floatgeom.Point2{
						float64(screenWidth) / 2, 16,
					}
				}
			}
			return 0
		}),
		btn.Binding(mouse.Release, func(_ event.CID, ev interface{}) int {
			if hdr.draggingWindow {
				hdr.draggingWindow = false
			}
			return 0
		}),
		btn.Binding(oak.WindowSizeChange, oak.SizeChangeEvent(func(c event.CID, pt intgeom.Point2) int {
			ctx.Window.(*oak.Window).UpdateViewSize(pt.X(), pt.Y())
			b, _ := ctx.CallerMap.GetEntity(c).(*btn.TextBox)
			newW := float64(pt.X()) - totalButtonsSize
			b.Box.R.Undraw()
			b.Box.R = render.NewColorBox(int(newW), int(construct.Height), construct.Color)
			ctx.DrawStack.Draw(b.Box.R, construct.Layers...)
			ctx.MouseTree.UpdateSpace(0, 0, newW, construct.Height, b.Box.Space)
			return 0
		})),
	)
	return hdr
}

var closeIcon = shape.JustIn(shape.AndIn(
	shape.XRange(.35, .65),
	func(x, y int, sizes ...int) bool {
		size := sizes[0]
		return x == y || y == (size-x)
	},
))

var minimizeIcon = shape.JustIn(shape.AndIn(
	shape.XRange(.35, .65),
	func(x, y int, sizes ...int) bool {
		return y == sizes[0]/2
	},
))

var maximizeIcon = shape.JustIn(squarePercent(.35, .65))

var normalizeIcon = shape.JustIn(shape.OrIn(
	squarePercent(.35, .65),
	squarePercent(.45, .55),
))

func squarePercent(minPerc, maxPerc float64) shape.In {
	return shape.AndIn(
		shape.XRange(minPerc-.03, maxPerc),
		func(x, y int, sizes ...int) bool {
			yf := float64(y)
			sf := float64(sizes[0])
			return (yf >= sf*(minPerc-.03)) && (yf <= sf*maxPerc)
		},
		func(x, y int, sizes ...int) bool {
			size := sizes[0]
			return x == int(float64(size)*minPerc) ||
				x == int(float64(size)*maxPerc) ||
				y == int(float64(size)*minPerc) ||
				y == int(float64(size)*maxPerc)
		},
	)
}

func toggleMaximize(ctx *scene.Context, b btn.Btn) bool {
	if sfx, _ := b.Metadata("switch-suffix"); sfx != "" {
		ctx.Window.(*oak.Window).Normalize()
		b.SetMetadata("switch-suffix", "")
		if sw, ok := b.GetRenderable().(*render.Switch); ok {
			sw.Set("nohover")
		}
		return false
	}
	ctx.Window.(*oak.Window).Maximize()
	b.SetMetadata("switch-suffix", "-revert")
	if sw, ok := b.GetRenderable().(*render.Switch); ok {
		sw.Set("nohover-revert")
	}
	return true
}
