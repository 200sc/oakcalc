package calc

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/200sc/oakcalc/internal/arith"
	"github.com/200sc/oakcalc/internal/components/titlebar"
	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/entities/x/btn"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/mouse"
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

			tokens := [][]arith.Token{
				{
					{Number: i64p(9)},
					{Number: i64p(8)},
					{Number: i64p(7)},
					{Op: opP(arith.OpDivide)},
				}, {
					{Number: i64p(6)},
					{Number: i64p(5)},
					{Number: i64p(4)},
					{Op: opP(arith.OpMultiply)},
				}, {
					{Number: i64p(3)},
					{Number: i64p(2)},
					{Number: i64p(1)},
					{Op: opP(arith.OpMinus)},
				}, {
					{Op: opP(arith.OpBackspace)},
					{Number: i64p(0)},
					{Op: opP(arith.OpEquals)},
					{Op: opP(arith.OpPlus)},
				},
			}
			var btnColor = colornames.Darkolivegreen
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
					btn.New(
						btn.Text(s),
						btn.Pos(x, y),
						btn.Width(width),
						btn.Height(height),
						btn.Color(btnColor),
						btn.Click(mouse.Binding(func(c event.CID, e *mouse.Event) int {
							fmt.Println("Would add ", s, " to stack")
							disp.Add(token)
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
	history          []string
	mu               sync.Mutex
	currentOperation []arith.Token
}

func (disp *arithmeticDisplay) Add(t arith.Token) {
	disp.mu.Lock()
	defer disp.mu.Unlock()
	// special cases
	if t.Op != nil && *t.Op == arith.OpEquals {
		tree, err := arith.Parse(disp.currentOperation)
		if err == nil {
			fmt.Println(tree.Eval())
		} else {
			fmt.Println(err)
		}
		disp.currentOperation = []arith.Token{}
		return
	}
	if t.Op != nil && *t.Op == arith.OpBackspace {
		if len(disp.currentOperation) != 0 {
			disp.currentOperation = disp.currentOperation[:len(disp.currentOperation)-1]
		}
		return
	}
	disp.currentOperation = append(disp.currentOperation, t)

}

type displayPosition struct {
	rect floatgeom.Rect2
	arith.Token
}

var displayPositions = []displayPosition{}

func i64p(i int64) *int64 {
	return &i
}

func opP(o arith.Op) *arith.Op {
	return &o
}
