package calc

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/200sc/oakcalc/internal/arith"
	"github.com/200sc/oakcalc/internal/components/titlebar"
	"github.com/oakmound/oak/v3/entities/x/btn"
	"github.com/oakmound/oak/v3/entities/x/mods"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/render"
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

			var disp arithmeticDisplay
			disp.ctx = ctx
			disp.fnt = render.DefaultFont()
			disp.current = disp.fnt.NewText("", 400, 430)
			ctx.DrawStack.Draw(disp.current, 9)

			tokens := [][]arith.Token{
				{
					{Number: i64p(7)},
					{Number: i64p(8)},
					{Number: i64p(9)},
					{Op: opP(arith.OpDivide)},
				}, {
					{Number: i64p(4)},
					{Number: i64p(5)},
					{Number: i64p(6)},
					{Op: opP(arith.OpMultiply)},
				}, {
					{Number: i64p(1)},
					{Number: i64p(2)},
					{Number: i64p(3)},
					{Op: opP(arith.OpMinus)},
				}, {
					{Op: opP(arith.OpBackspace)},
					{Number: i64p(0)},
					{Op: opP(arith.OpEquals)},
					{Op: opP(arith.OpPlus)},
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
			for _, tokenRow := range tokens {
				for _, token := range tokenRow {
					token := token
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
						btn.Pos(x, y),
						btn.Width(width),
						btn.Height(height),
						btn.Renderable(r),
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
					)
					x += width + xSpacing
				}
				x = xStart
				y += height + ySpacing
			}

			// Result display
			// 123 456 789
			// + - * /
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
	disp.ctx.DrawStack.Draw(txt)
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
			pretty := tree.Pretty()
			result := tree.Eval()
			fmt.Println(pretty)
			fmt.Println(result)
			disp.AddToHistory(pretty)
			disp.AddToHistory(" = " + strconv.FormatInt(result, 10))
		} else {
			fmt.Println(err)
		}
		disp.currentOperation = []arith.Token{}
		disp.current.SetString("")
		return
	}
	defer func() {
		tree, err := arith.Parse(disp.currentOperation)
		if err == nil {
			pretty := tree.Pretty()
			disp.current.SetString(pretty)
		} else {
			fmt.Println(err)
		}
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
