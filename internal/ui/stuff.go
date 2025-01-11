package ui

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Window struct{}

func Test() {
	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagWindowHighdpi | rl.FlagVsyncHint | rl.FlagMsaa4xHint)
	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	noto := rl.LoadFont("./noto-sans.ttf")
	txtImg := rl.ImageTextEx(noto, "Hello World!", 50, 0, rl.Black)
	txt := rl.LoadTextureFromImage(txtImg)

	tmpRect := NewRect(10, 10, 250, 500)
	tmpRect.SetBorderRadius(25)

	go func() {
		for {
			time.Sleep(5 * time.Second)
			// tmpRect.SetBorderRadius(0)
			// tmpRect.SetPosition(100, 100)
			tmpRect.SetSize(400, 700)
			time.Sleep(5 * time.Second)
			// tmpRect.SetBorderRadius(25)
			// tmpRect.SetPosition(10, 10)
			tmpRect.SetSize(250, 500)
		}
	}()
	// tmpRect.SetFillColor(rl.Green)
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.DrawFPS(0, 0)

		rl.ClearBackground(rl.DarkGray)
		// rl.DrawText("Congrats! You created your first window!", 190, 100, 20, rl.LightGray)
		// // rl.DrawTextEx(noto, "Congrats! You created your first window!", rl.NewVector2(190, 200), 48, 0, rl.Black)
		tmpRect.Draw()
		rl.DrawTexture(txt, 50, 50, rl.Black)
		rl.EndDrawing()
	}
}
