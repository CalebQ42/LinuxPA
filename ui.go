package main

import (
	"os/exec"

	"github.com/nelsam/gxui"
)

func ui() {
	catListAdap := &StrList{}
	appListAdap := &catAdap{}
	catListAdap.SetStrings(lin)
	win := th.CreateWindow(500, 500, "LinuxPA")
	top := th.CreateLinearLayout()
	top.SetDirection(gxui.BottomToTop)
	splBox := th.CreateLinearLayout()
	spl := th.CreateSplitterLayout()
	spl.SetOrientation(gxui.Horizontal)
	catList := th.CreateList()
	catList.SetAdapter(catListAdap)
	catList.OnSelectionChanged(func(it gxui.AdapterItem) {
		appListAdap.setCat(it.(string))
	})
	appList := th.CreateTree()
	appList.SetAdapter(appListAdap)
	spl.AddChild(catList)
	spl.AddChild(appList)
	splBox.AddChild(spl)
	butBox := th.CreateLinearLayout()
	butBox.SetDirection(gxui.LeftToRight)
	if _, err := exec.LookPath("wine"); err == nil {
		wineBut := th.CreateButton()
		wineBut.SetType(gxui.ToggleButton)
		wineBut.SetChecked(wine)
		wineBut.SetText("Show Windows Apps")
		wineBut.OnClick(func(gxui.MouseEvent) {
			wine = wineBut.IsChecked()
			appListAdap.refresh()
			if wineBut.IsChecked() {
				catListAdap.SetStrings(cats)
				wineBut.SetText("Hide Windows Apps")
			} else {
				catListAdap.SetStrings(lin)
				wineBut.SetText("Show Windows Apps")
			}
		})
		butBox.AddChild(wineBut)
	}
	top.AddChild(butBox)
	top.AddChild(splBox)
	win.AddChild(top)
	win.OnClose(dr.Terminate)
}
