package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	fyneApp fyne.App
)

func main() {
	fyneApp = app.New()
	win := fyneApp.NewWindow("LinuxPA")
	win.SetContent(container.NewVBox(widget.NewLabel("HII"), widget.NewLabel("PPOOOP")))
	win.Resize(fyne.NewSize(512, 512))
	win.ShowAndRun()
}
