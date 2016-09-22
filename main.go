package main

import (
	"github.com/nelsam/gxui"
	"github.com/nelsam/gxui/drivers/gl"
	"github.com/nelsam/gxui/themes/dark"
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
)

func main() {
	master = make(map[string][]app)
	linmaster = make(map[string][]app)
	gl.StartDriver(appMain)
}

func appMain(dri gxui.Driver) {
	dr = dri
	th = dark.CreateTheme(dr)
	setup()
	ui()
}
