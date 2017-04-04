package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/nelsam/gxui"
	"github.com/nelsam/gxui/drivers/gl"
	"github.com/nelsam/gxui/themes/dark"
	"github.com/nelsam/gxui/themes/light"
)

const (
	version = "1.1.0.0"
	defIni  = "[basic]\ntheme=dk"
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
		master = make(map[string][]app)
		linmaster = make(map[string][]app)
		gl.StartDriver(appMain)
	}
}

func appMain(dri gxui.Driver) {
	dr = dri
	setup()
	if darkTheme {
		th = dark.CreateTheme(dr)
	} else {
		th = light.CreateTheme(dr)
	}
	th = dark.CreateTheme(dr)
	ui()
}

func contains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}
