package main

import (
	"fmt"
	"os"

	"github.com/200sc/oakcalc/internal/calc"
	"github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/render"
)

func main() {

	render.SetDrawStack(render.NewStaticHeap())

	oak.AddScene(calc.SceneName, calc.Scene())
	err := oak.Init(calc.SceneName, func(c oak.Config) (oak.Config, error) {
		c.Title = "OakCalc"
		c.Borderless = true
		return c, nil
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
