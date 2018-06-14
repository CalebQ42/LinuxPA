package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gotk3/gotk3/gtk"
)

const (
	version = "2.1.4.1"
)

var (
	master    map[string][]app
	linmaster map[string][]app
	cats      []string
	lin       []string
	comEnbld  bool
	populated bool
)

func main() {
	forced := flag.Bool("force-update", false, "Force the update dialog to be shown")
	flag.Parse()
	os.MkdirAll("PortableApps/LinuxPACom", 0777)
	master = make(map[string][]app)
	linmaster = make(map[string][]app)
	uiStart(*forced)
}

func uiStart(forced bool) {
	gtk.Init(nil)
	setup()
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		fmt.Println("Window not created", err)
	}
	win.SetTitle("LinuxPA")
	win.Connect("destroy", func() {
		savePrefs()
		gtk.MainQuit()
	})
	win.SetDefaultSize(500, 500)
	win.SetPosition(gtk.WIN_POS_CENTER)
	ui(win)
	win.ShowAll()
	win.Show()
	update(win, forced)
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
