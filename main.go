package main

import (
	"fmt"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/widget/material"
)

const (
	version = "3.0.0-alpha1"
)

func main() {
	go func() {
		w := &app.Window{}
		w.Option(app.Size(500, 500), app.Title("LinuxPA"))
		if err := ui(w); err != nil {
			fmt.Println(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func ui(w *app.Window) error {
	th := material.NewTheme()
	op := &op.Ops{}
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			ctx := app.NewContext(op, e)
			lbl := material.Body1(th, "Hello there!")
			lbl.Layout(ctx)
			e.Frame(op)
		}
	}
}
