package main

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/CalebQ42/LinuxPA/appimg"

	"github.com/gotk3/gotk3/gtk"
)

const (
	version = "2.1.4.1"
)

var (
	master        map[string][]app
	linmaster     map[string][]app
	cats          []string
	lin           []string
	wine          bool
	comEnbld      bool
	wineAvail     bool
	portableHide  bool
	betaUpdate    bool
	versionNewest = true
	paDirs        = true
)

func main() {
	_, err := appimg.GetDesktopFile("gimp.AppImage")
	if err == nil {
		os.Rename("/tmp/my.desktop", "PortableApps/stuff.desktop")
	} else {
		println("error:", err)
	}
	// forced := flag.Bool("force-update", false, "Force the update dialog to be shown")
	// flag.Parse()
	// os.MkdirAll("PortableApps/LinuxPACom", 0777)
	// master = make(map[string][]app)
	// linmaster = make(map[string][]app)
	// uiStart(*forced)
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

func savePrefs() {
	os.Remove("PortableApps/LinuxPACom/Prefs.gob")
	fil, err := os.Create("PortableApps/LinuxPACom/Prefs.gob")
	if err != nil {
		return
	}
	enc := gob.NewEncoder(fil)
	err = enc.Encode(wine)
	if err != nil {
		return
	}
	err = enc.Encode(portableHide)
	if err != nil {
		return
	}
	err = enc.Encode(versionNewest)
	if err != nil {
		return
	}
	err = enc.Encode(paDirs)
	if err != nil {
		return
	}
	err = enc.Encode(betaUpdate)
	if err != nil {
		return
	}
}

func loadPrefs() {
	fil, err := os.Open("PortableApps/LinuxPACom/Prefs.gob")
	if err != nil {
		return
	}
	dec := gob.NewDecoder(fil)
	err = dec.Decode(&wine)
	if err != nil {
		return
	}
	err = dec.Decode(&portableHide)
	if err != nil {
		return
	}
	err = dec.Decode(&versionNewest)
	if err != nil {
		return
	}
	err = dec.Decode(&paDirs)
	if err != nil {
		return
	}
	err = dec.Decode(&betaUpdate)
	if err != nil {
		return
	}
}

func contains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}
