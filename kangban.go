package main

import tl "github.com/JoelOtter/termloop"

func main() {
	game := tl.NewGame()
	game.Start()

	level := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorGreen,
		Fg: tl.ColorBlack,
		Ch: 'v',
	})

	level.AddEntity(tl.NewRectangle(10, 10, 50, 20, tl.ColorBlue))
}
