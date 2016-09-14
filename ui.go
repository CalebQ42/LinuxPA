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
	top := th.CreateLinearLayout()
	top.SetDirection(gxui.BottomToTop)
	top.SetHorizontalAlignment(gxui.AlignRight)
	spl := th.CreateSplitterLayout()
	spl.SetOrientation(gxui.Horizontal)
	catlist := th.CreateList()
	catlist.SetAdapter(catAdap)
	catlist.OnItemClicked(func(_ gxui.MouseEvent, it gxui.AdapterItem) {
		str := it.(string)
		appAdap.SetApps(appMaster[str])
	})
	applist := th.CreateList()
	applist.SetAdapter(appAdap)
	spl.AddChild(catlist)
	spl.AddChild(applist)
	but := th.CreateLinearLayout()
	but.SetDirection(gxui.RightToLeft)
	launch := th.CreateButton()
	launch.SetText("Launch!")
	launch.OnClick(func(gxui.MouseEvent) {
		if appAdap.ItemIndex(applist.Selected()) != -1 {
			app := applist.Selected().(prtap)
			dir, fi := path.Split(app.ex)
			var cmd *exec.Cmd
			if app.wine {
				cmd = exec.Command("/bin/sh", "-c", "cd \""+dir+"\"; wine \""+fi+"\"")
			} else {
				if commEnbl {
					cmd = exec.Command("/bin/sh", "-c", ". "+common+" || exit 1;cd \""+dir+"\"; \"./"+fi+"\"")
				} else {
					cmd = exec.Command("/bin/sh", "-c", "cd \""+dir+"\"; \"./"+fi+"\"")
				}
			}
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Start()
		}
	})
	wine := th.CreateButton()
	wine.SetType(gxui.ToggleButton)
	wine.OnClick(func(gxui.MouseEvent) {
		if wine.IsChecked() {
			appAdap.Wine(true)
		} else {
			appAdap.Wine(false)
		}
	})
	wine.SetText("Show Windows apps")
	wine.SetChecked(appAdap.wine)
	but.AddChild(wine)
	but.AddChild(launch)
	top.AddChild(but)
	top.AddChild(spl)
	win.AddChild(top)
	win.OnClose(func() {
		dr.Terminate()
	})
}
