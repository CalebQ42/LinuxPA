package main

import (
	"os"
	"os/exec"
	"path"

	"github.com/nelsam/gxui"
	"github.com/nelsam/gxui/themes/dark"
)

var (
	dr gxui.Driver
)

func uiMain(dri gxui.Driver) {
	dr = dri
	catAdap := &StrList{}
	catAdap.SetStrings(cats)
	appAdap := &prtapAdap{}
	th := dark.CreateTheme(dr)
	win := th.CreateWindow(500, 500, "LinuxPA")
	top := th.CreateSplitterLayout()
	top.SetOrientation(gxui.Horizontal)
	catlist := th.CreateList()
	catlist.SetAdapter(catAdap)
	catlist.OnItemClicked(func(_ gxui.MouseEvent, it gxui.AdapterItem) {
		str := it.(string)
		appAdap.SetApps(appMaster[str])
	})
	applist := th.CreateList()
	applist.SetAdapter(appAdap)
	applist.OnItemClicked(func(_ gxui.MouseEvent, it gxui.AdapterItem) {
		app := it.(prtap)
		dir, fi := path.Split(app.ex)
		cmd := exec.Command("/bin/sh", "-c", "cd \""+dir+"\"; \"./"+fi+"\"")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Start()
	})
	top.AddChild(catlist)
	top.AddChild(applist)
	win.AddChild(top)
	win.OnClose(dr.Terminate)
}
