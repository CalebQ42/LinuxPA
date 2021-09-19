package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	app := app.New()
	win := app.NewWindow("LinuxPA")
	win.SetContent(container.NewVBox(widget.NewLabel("HII"), widget.NewLabel("PPOOOP")))
	win.Resize(fyne.NewSize(512, 512))
	win.ShowAndRun()
}
