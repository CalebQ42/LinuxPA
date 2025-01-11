package ui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Window struct{}

func Test() {
	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	rl.SetWindowState(rl.FlagWindowResizable | rl.FlagWindowHighdpi)
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	defer rl.CloseWindow()
	rl.Scalef(2, 2, 2)

	noto := rl.LoadFont("./noto-sans.ttf")
	txtImg := rl.ImageTextEx(noto, "Hello World!", 50, 0, rl.Black)
	// rl.ImageResize(txtImg, 450, 200)
	txt := rl.LoadTextureFromImage(txtImg)

	tmpRect := NewRect(10, 10, 250, 500)
	tmpRect.SetBorderRadius(25)
	// tmpRect.SetFillColor(rl.Green)
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.DarkGray)
		// rl.DrawText("Congrats! You created your first window!", 190, 100, 20, rl.LightGray)
		// // rl.DrawTextEx(noto, "Congrats! You created your first window!", rl.NewVector2(190, 200), 48, 0, rl.Black)
		tmpRect.Draw()
		rl.DrawTexture(txt, 50, 50, rl.Black)
		rl.EndDrawing()
	}
}
