package main

import (
	"fmt"

	"github.com/nelsam/gxui"
	"github.com/nelsam/gxui/drivers/gl"
	"github.com/nelsam/gxui/themes/dark"
	"github.com/nelsam/gxui/themes/light"
)

const (
	version = "0.1.1.1"
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
	stat := versionDL()
	if stat {
		res := getVersionFileInfo()
		if res != "Error!" {
			stat = checkForUpdate(res)
			if stat {
				downloadUpdate(res)
			} else {
				fmt.Println("Failed DL")
			}
		} else {
			fmt.Println("Failed Version File Info")
		}
	} else {
		fmt.Println("Failed Version DL")
	}
	master = make(map[string][]app)
	linmaster = make(map[string][]app)
	gl.StartDriver(appMain)
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
