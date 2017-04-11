package main

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/gotk3/gotk3/gtk"
)

const (
	version = "2.1.0.5"
	defIni  = ""
)

var (
	master    map[string][]app
	linmaster map[string][]app
	cats      []string
	lin       []string
	wine      bool
	comEnbld  bool
	wineAvail bool
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
		savePrefs()
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

func savePrefs() {
	fil, err := os.Open("PortableApps/LinuxPACom/Prefs.gob")
	if os.IsNotExist(err) {
		fil, err = os.Create("PortableApps/LinuxPACom/Prefs.gob")
	}
	if err != nil {
		return
	}
	enc := gob.NewEncoder(fil)
	enc.Encode(wine)
}

func loadPrefs() {
	fil, err := os.Open("PortableApps/LinuxPACom/Prefs.gob")
	if err != nil {
		return
	}
	dec := gob.NewDecoder(fil)
	dec.Decode(&wine)
}

func contains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}
