package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

var (
	fyneApp fyne.App
)

func main() {
	fyneApp = app.New()
	win := fyneApp.NewWindow("LinuxPA")
	apps = []*portableApp{
		{
			name: "FIRE",
			execs: []string{
				"yellow.sh",
				"Potato",
			},
		},
		{
			name: "SUPER FIRE",
			execs: []string{
				"super yellow.sh",
				"super Potato",
			},
		},
	}
	tree := buildAppList()
	win.SetContent(tree)
	tree.(*widget.Tree).OpenAllBranches()
	win.Resize(fyne.NewSize(512, 512))
	win.ShowAndRun()
}

func buildAppList() fyne.CanvasObject {
	m := make(map[string][]string)
	name := make([]string, 0)
	for _, a := range apps {
		name = append(name, a.name)
	}
	m[""] = name
	for _, a := range apps {
		m[a.name] = a.execs
	}
	return widget.NewTreeWithStrings(m)
}
