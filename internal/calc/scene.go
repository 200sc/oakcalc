package calc

import (
	"image"
	"image/color"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/200sc/oakcalc/internal/arith"
	"github.com/200sc/oakcalc/internal/components/titlebar"
	"github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/entities/x/btn"
	"github.com/oakmound/oak/v3/entities/x/mods"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
	"golang.org/x/image/colornames"
	mkey "golang.org/x/mobile/event/key"
)

const SceneName = "calc"

type tokenWithShortcut struct {
	arith.Token
	shortcutRune rune
	shortcutKey  mkey.Code // for keys without runes
}

func Scene() scene.Scene {
	return scene.Scene{
		Start: func(ctx *scene.Context) {
			titlebar.New(ctx,
				titlebar.WithColor(colornames.Darkgreen),
				titlebar.WithHeight(32),
				titlebar.WithLayers([]int{10}),
				titlebar.WithTitle("OakCalc"),
			)

			var disp arithmeticDisplay
			disp.ctx = ctx
			disp.fnt = render.DefaultFont()
			disp.fnt.Fallbacks = loadFallbackFonts(10)
			disp.current = disp.fnt.NewText("", 400, 430)
			ctx.DrawStack.Draw(disp.current, 9)
			ctx.Window.(*oak.Window).SetColorBackground(image.NewUniform(color.RGBA{0, 20, 0, 255}))

			tokens := [][]tokenWithShortcut{
				{
					{
						Token:        arith.Token{Number: i64p(7)},
						shortcutRune: '7',
					},
					{
						Token:        arith.Token{Number: i64p(8)},
						shortcutRune: '8',
					},
					{
						Token:        arith.Token{Number: i64p(9)},
						shortcutRune: '9',
					},
					{
						Token:        arith.Token{Op: opP(arith.OpDivide)},
						shortcutRune: '/',
					},
					{
						Token:        arith.Token{Op: opP(arith.OpOpenParen)},
						shortcutRune: '(',
					},
				}, {
					{
						Token:        arith.Token{Number: i64p(4)},
						shortcutRune: '4',
					},
					{
						Token:        arith.Token{Number: i64p(5)},
						shortcutRune: '5',
					},
					{
						Token:        arith.Token{Number: i64p(6)},
						shortcutRune: '6',
					},
					{
						Token:        arith.Token{Op: opP(arith.OpMultiply)},
						shortcutRune: '*',
					},
					{
						Token:        arith.Token{Op: opP(arith.OpCloseParen)},
						shortcutRune: ')',
					},
				}, {
					{
						Token:        arith.Token{Number: i64p(1)},
						shortcutRune: '1',
					},
					{
						Token:        arith.Token{Number: i64p(2)},
						shortcutRune: '2',
					},
					{
						Token:        arith.Token{Number: i64p(3)},
						shortcutRune: '3',
					},
					{
						Token:        arith.Token{Op: opP(arith.OpMinus)},
						shortcutRune: '-',
					},
					{
						Token:        arith.Token{Op: opP(arith.OpSquareRoot)},
						shortcutRune: 'q',
					},
				}, {
					{
						Token:       arith.Token{Op: opP(arith.OpBackspace)},
						shortcutKey: mkey.CodeDeleteBackspace,
					},
					{
						Token:        arith.Token{Number: i64p(0)},
						shortcutRune: '0',
					},
					{
						Token:       arith.Token{Op: opP(arith.OpEquals)},
						shortcutKey: mkey.CodeReturnEnter,
					},
					{
						Token:        arith.Token{Op: opP(arith.OpPlus)},
						shortcutRune: '+',
					},
				},
			}
			btnColor := colornames.Darkolivegreen
			highlightColor := mods.Lighter(btnColor, .10)
			pressColor := mods.Lighter(btnColor, .20)
			const width = 50
			const height = 50
			const xSpacing = 10
			const ySpacing = 10
			const xStart = 20
			const yStart = 200
			var x float64 = xStart
			var y float64 = yStart
			btnFnt, _ := render.DefaultFont().RegenerateWith(func(fg render.FontGenerator) render.FontGenerator {
				fg.Size = 25
				return fg
			})

			btnFnt.Fallbacks = loadFallbackFonts(25)
			for _, tokenRow := range tokens {
				for _, tokenShortcut := range tokenRow {
					token := tokenShortcut.Token
					shortcutRune := tokenShortcut.shortcutRune
					shortcutKey := tokenShortcut.shortcutKey
					s := ""
					if token.Op != nil {
						s = string(*token.Op)
					} else {
						s = strconv.FormatInt(*token.Number, 10)
					}
					r := render.NewSwitch("nohover", map[string]render.Modifiable{
						"nohover": render.NewColorBox(width, height, btnColor),
						"hover":   render.NewColorBox(width, height, highlightColor),
						"onpress": render.NewColorBox(width, height, pressColor),
					})
					btn.New(
						btn.Text(s),
						btn.TxtOff(12, 8),
						btn.Font(btnFnt),
						btn.Pos(x, y),
						btn.Width(width),
						btn.Height(height),
						btn.Renderable(r),
						btn.Layers(1),
						btn.Click(mouse.Binding(func(c event.CID, e *mouse.Event) int {
							disp.Add(token)
							return 0
						})),
						btn.Binding(mouse.Start, mouse.Binding(func(c event.CID, e *mouse.Event) int {
							b, _ := ctx.CallerMap.GetEntity(c).(btn.Btn)
							if sw, ok := b.GetRenderable().(*render.Switch); ok {
								sw.Set("hover")
							}
							return 0
						})),
						btn.Binding(mouse.Stop, mouse.Binding(func(c event.CID, e *mouse.Event) int {
							b, _ := ctx.CallerMap.GetEntity(c).(btn.Btn)
							if sw, ok := b.GetRenderable().(*render.Switch); ok {
								sw.Set("nohover")
							}
							return 0
						})),
						btn.Binding(mouse.PressOn, mouse.Binding(func(c event.CID, e *mouse.Event) int {
							b, _ := ctx.CallerMap.GetEntity(c).(btn.Btn)
							if sw, ok := b.GetRenderable().(*render.Switch); ok {
								sw.Set("onpress")
							}
							return 0
						})),
						btn.Binding(key.Down, func(c event.CID, i interface{}) int {
							kv, ok := i.(key.Event)
							if !ok {
								return 0
							}
							pressed := false
							if shortcutRune != 0 && kv.Rune == shortcutRune {
								pressed = true
							} else if shortcutKey != 0 && kv.Code == shortcutKey {
								pressed = true
							}
							if pressed {
								b, _ := ctx.CallerMap.GetEntity(c).(btn.Btn)
								if sw, ok := b.GetRenderable().(*render.Switch); ok {
									sw.Set("onpress")
									disp.Add(token)
									ctx.DoAfter(50*time.Millisecond, func() {
										sw.Set("nohover")
									})
								}
							}
							return 0
						}),
					)
					x += width + xSpacing
				}
				x = xStart
				y += height + ySpacing
			}

			bkg := render.NewColorBoxR(395, 480, color.RGBA{50, 75, 50, 255})
			ctx.DrawStack.Draw(bkg, 0)
		},
		// No Loop function, this is the only scene.
		// No End function, this is the only scene.
	}
}

type arithmeticDisplay struct {
	render.LayeredPoint
	ctx     *scene.Context
	fnt     *render.Font
	current *render.Text

	history          []*render.Text
	mu               sync.Mutex
	currentOperation []arith.Token
}

func (disp *arithmeticDisplay) AddToHistory(s string) {
	const textheight = 30
	const textX = 400
	const textY = 400
	for _, h := range disp.history {
		h.ShiftY(-textheight)
	}
	txt := disp.fnt.NewText(s, textX, textY)
	disp.ctx.DrawStack.Draw(txt, 1)
	disp.history = append(disp.history, txt)
}

func (disp *arithmeticDisplay) Add(t arith.Token) {
	disp.mu.Lock()
	defer disp.mu.Unlock()
	// special cases
	if t.Op != nil && *t.Op == arith.OpEquals {
		if len(disp.currentOperation) == 0 {
			disp.currentOperation = append(disp.currentOperation, arith.Token{
				Number: i64p(0),
			})
		}
		tree, err := arith.Parse(disp.currentOperation)
		if err == nil {
			result := arith.Eval(tree)
			pretty := arith.Pretty(tree)
			disp.AddToHistory(pretty)
			disp.AddToHistory(" = " + strconv.FormatInt(result, 10))
		}
		disp.currentOperation = []arith.Token{}
		disp.current.SetString("")
		return
	}
	defer func() {
		strs := make([]string, len(disp.currentOperation))
		for i, t := range disp.currentOperation {
			if t.Op != nil {
				strs[i] = string(*t.Op)
			} else {
				strs[i] = strconv.FormatInt(*t.Number, 10)
			}
		}
		pretty := strings.Join(strs, " ")
		disp.current.SetString(pretty)
	}()
	if t.Op != nil && *t.Op == arith.OpBackspace {
		if len(disp.currentOperation) != 0 {
			disp.currentOperation = disp.currentOperation[:len(disp.currentOperation)-1]
		}
		return
	}
	if t.Number != nil {
		if len(disp.currentOperation) != 0 {
			if disp.currentOperation[len(disp.currentOperation)-1].Number != nil {
				// combine the two numbers
				*disp.currentOperation[len(disp.currentOperation)-1].Number *= 10
				*disp.currentOperation[len(disp.currentOperation)-1].Number += *t.Number
				return
			}
		}
	}
	disp.currentOperation = append(disp.currentOperation, t.Copy())
}

func i64p(i int64) *int64 {
	return &i
}

func opP(o arith.Op) *arith.Op {
	return &o
}
