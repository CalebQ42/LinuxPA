package main

import (
	"io/ioutil"
	"os"

	"github.com/gotk3/gotk3/gtk"
)

func settingsUI() {
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.SetDefaultSize(600, 300)
	win.SetPosition(gtk.WIN_POS_CENTER)
	comTagTbl, _ := gtk.TextTagTableNew()
	comBuf, _ := gtk.TextBufferNew(comTagTbl)
	ntbk, _ := gtk.NotebookNew()
	gnrl, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	gnrl.SetMarginStart(10)
	gnrl.SetMarginEnd(10)
	gnrl.SetMarginTop(10)
	gnrl.SetMarginBottom(10)
	dlWine, _ := gtk.ButtonNewWithLabel("Download Wine")
	wineLbl, _ := gtk.LabelNew("PortableApps/LinuxPACom/Wine present")
	gnrl.Add(wineLbl)
	dlWine.Connect("clicked", func() {
		cb := make(chan bool)
		downloadWine(win, cb)
		go func() {
			v := <-cb
			if v {
				setupTxt(comBuf)
				wineLbl.Show()
			}
		}()
	})
	if !comEnbld {
		dlWine.SetSensitive(false)
		dlWine.SetTooltipText("common.sh needed")
	}
	gnrl.Add(dlWine)
	ntbk.AppendPage(gnrl, getLabel("General"))
	com, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	com.SetMarginStart(10)
	com.SetMarginEnd(10)
	com.SetMarginTop(10)
	com.SetMarginBottom(10)
	comEdit, _ := gtk.TextViewNewWithBuffer(comBuf)
	comEdit.SetVExpand(true)
	comEdit.SetHExpand(true)
	vScroll, _ := gtk.AdjustmentNew(0, 0, 0, 0, 0, 0)
	hScroll, _ := gtk.AdjustmentNew(0, 0, 0, 0, 0, 0)
	comScrl, _ := gtk.ScrolledWindowNew(hScroll, vScroll)
	comScrl.Add(comEdit)
	com.Add(comScrl)
	svBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	sv, _ := gtk.ButtonNewWithLabel("Save")
	sv.Connect("clicked", func() {
		beg, end := comBuf.GetBounds()
		txt, _ := comBuf.GetText(beg, end, true)
		ioutil.WriteFile("PortableApps/LinuxPACom/common.sh", []byte(txt), 0777)
	})
	cnl, _ := gtk.ButtonNewWithLabel("Cancel")
	cnl.Connect("clicked", func() {
		setupTxt(comBuf)
	})
	svBox.Add(sv)
	svBox.Add(cnl)
	com.Add(svBox)
	ntbk.AppendPage(com, getLabel("common.sh"))
	win.Add(ntbk)
	win.ShowAll()
	if !comEnbld {
		comScrl.Hide()
		svBox.Hide()
		mkCom, _ := gtk.ButtonNewWithLabel("Create common.sh")
		mkCom.Connect("clicked", func() {
			err := ioutil.WriteFile("PortableApps/LinuxPACom/common.sh", []byte("export HOME=$PWD/PortableApps/LinuxPACom/Home"), 0777)
			if err == nil {
				mkCom.Hide()
				comScrl.Show()
				svBox.Show()
				setupTxt(comBuf)
				comEnbld = true
				dlWine.SetSensitive(true)
				dlWine.SetTooltipText("")
			}
		})
		mkCom.Show()
		com.Add(mkCom)
	} else {
		setupTxt(comBuf)
	}
	if _, err := os.Open("PortableApps/LinuxPACom/Wine"); err != nil && os.IsNotExist(err) {
		wineLbl.Hide()
	}
	win.Show()
}

func setupTxt(buf *gtk.TextBuffer) {
	fil, _ := os.Open("PortableApps/LinuxPACom/common.sh")
	btys, _ := ioutil.ReadAll(fil)
	buf.SetText(string(btys))
}

func getLabel(name string) *gtk.Label {
	lbl, _ := gtk.LabelNew(name)
	return lbl
}
