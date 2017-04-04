package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gotk3/gotk3/gtk"
	"github.com/nelsam/gxui"
)

const (
	version = "2.0.0.0"
)

var (
	dr        gxui.Driver
	th        gxui.Theme
	master    map[string][]app
	linmaster map[string][]app
	cats      []string
	lin       []string
	wine      bool
	comEnbld  bool
	darkTheme = true
)

func main() {
	updated := false
	os.MkdirAll("PortableApps/LinuxPACom", 0777)
	stat, err := versionDL()
	if stat {
		res := getVersionFileInfo()
		if res != "Error!" {
			stat, err = checkForUpdate(res)
			if stat {
				downloadUpdate(res)
				updated = true
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Failed Version File Info")
		}
	} else {
		fmt.Println(err)
	}
	if updated {
		cmd := exec.Command("./LinuxPA")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Start()
	} else {
		// master = make(map[string][]app)
		// linmaster = make(map[string][]app)
		uiStart()
	}
}

func uiStart() {
	gtk.Init(nil)
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		fmt.Println("Window not created", err)
	}
	win.SetTitle("LinuxPA")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	win.SetDefaultSize(500, 500)
	ui(win)
	gtk.Main()
}

// func appMain(dri gxui.Driver) {
// 	dr = dri
// 	setup()
// 	if darkTheme {
// 		th = dark.CreateTheme(dr)
// 	} else {
// 		th = light.CreateTheme(dr)
// 	}
// 	th = dark.CreateTheme(dr)
// 	ui()
// }

func contains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}
