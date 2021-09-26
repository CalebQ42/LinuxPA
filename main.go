package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

var (
	fyneApp fyne.App
)

func main() {
	err := processApps()
	if err != nil {
		log.Println(err)
		return
	}
	buildAndStartFyneUI()
}

func buildAndStartFyneUI() {
	fyneApp = app.New()
	win := fyneApp.NewWindow("LinuxPA")
	tree := buildFyneAppList()
	win.SetContent(tree)
	tree.OpenAllBranches()
	win.Resize(fyne.NewSize(512, 512))
	win.ShowAndRun()
}

func buildFyneAppList() *widget.Tree {
	m := make(map[string][]string)
	name := make([]string, 0)
	for _, a := range apps {
		name = append(name, a.name)
	}
	m[""] = name
	for _, a := range apps {
		m[a.name] = a.fyneTreeStrings()
	}
	return widget.NewTreeWithStrings(m)
}
