package main

import (
	"fmt"

	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
)

const (
	version = "3.0.0-alpha1"
)

func main() {
	b := core.NewBody("LinuxPA")
	core.NewButton(b).
		SetText("This is a test").
		OnClick(func(_ events.Event) {
			fmt.Println("This is a test")
		})
	b.RunMainWindow()
}
