package main

import (
	"fmt"
	"os"

	"github.com/gotk3/gotk3/gtk"
)

const (
	version = "2.0.0.1"
	defIni  = ""
)

var (
	master    map[string][]app
	linmaster map[string][]app
	cats      []string
	lin       []string
	wine      bool
	comEnbld  bool
	darkTheme = true
)

func main() {
	os.MkdirAll("PortableApps/LinuxPACom", 0777)
	master = make(map[string][]app)
	linmaster = make(map[string][]app)
	uiStart()
}

func uiStart() {
	gtk.Init(nil)
	setup()
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		fmt.Println("Window not created", err)
	}
	win.SetTitle("LinuxPA")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	win.SetDefaultSize(500, 500)
	win.SetPosition(gtk.WIN_POS_CENTER)
	ui(win)
	win.ShowAll()
	win.Show()
	update(win)
	gtk.Main()
}

func contains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}
