package main

import (
	"fmt"
	"os"

	"gioui.org/app"
)

const (
	version = "3.0.0-alpha1"
)

func main() {
	go func() {
		w := &app.Window{}
		if err := ui(w); err != nil {
			fmt.Println(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func ui(w *app.Window) error {
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:

		}
	}
}
