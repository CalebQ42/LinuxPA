package main

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	commonHelp = "The common.sh is run before every app is launched and allows you to set variables such as $HOME. For directories, ALWAYS start the directory with $PWD which points to the directory where LinuxPA is. To allow for greater customization and isolation, you can use the $PANAME variable which is the filename of the executable you're using."
)

func settingsUI(parent *gtk.Window, onExit func()) {
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.SetTransientFor(parent)
	parent.SetSensitive(false)
	win.SetDefaultSize(600, 300)
	win.SetPosition(gtk.WIN_POS_CENTER_ON_PARENT)
	win.Connect("destroy", func() {
		parent.SetSensitive(true)
		onExit()
	})
	comTagTbl, _ := gtk.TextTagTableNew()
	comBuf, _ := gtk.TextBufferNew(comTagTbl)
	ntbk, _ := gtk.NotebookNew()
	gnrl, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	gnrl.SetMarginStart(10)
	gnrl.SetMarginEnd(10)
	gnrl.SetMarginTop(10)
	gnrl.SetMarginBottom(10)
	dlWine, _ := gtk.ButtonNewWithLabel("Download Wine")
	wineCheck, _ := gtk.CheckButtonNewWithLabel("Show Windows apps (Wine)")
	wineLbl, _ := gtk.LabelNew("PortableApps/LinuxPACom/Wine present")
	dlWine.Connect("clicked", func() {
		cb := make(chan bool)
		downloadWine(win, cb)
		go func() {
			v := <-cb
			if v {
				setupTxt(comBuf)
				wineLbl.Show()
			}
			if _, err := os.Open("PortableApps/LinuxPACom/Wine"); os.IsNotExist(err) {
				if _, errd := exec.LookPath("wine"); errd == nil {
					wineAvail = true
				}
			} else if err == nil {
				wineAvail = true
			}
			glib.IdleAdd(func() {
				if !wineAvail {
					wineCheck.SetSensitive(false)
					wineCheck.SetTooltipText("Download wine to run windows apps")
				} else {
					wineCheck.SetSensitive(true)
					wineCheck.SetTooltipText("")
				}
			})
		}()
	})
	if !comEnbld {
		dlWine.SetSensitive(false)
		dlWine.SetTooltipText("common.sh needed")
	}
	pthdCheck, _ := gtk.CheckButtonNewWithLabel("Hide \"Portable\" from app name")
	pthdCheck.Connect("toggled", func() {
		portableHide = pthdCheck.GetActive()
		master = make(map[string][]app)
		linmaster = make(map[string][]app)
		cats = make([]string, 0)
		lin = make([]string, 0)
		setup()
	})
	pthdCheck.SetActive(portableHide)
	if !wineAvail {
		wineCheck.SetSensitive(false)
		wineCheck.SetTooltipText("Download wine to run windows apps")
	}
	wineCheck.SetActive(wine)
	wineCheck.Connect("toggled", func() {
		wine = wineCheck.GetActive()
	})
	versCheck, _ := gtk.CheckButtonNewWithLabel("Only show newest app version in downloads (A bit iffy ATM)")
	versCheck.SetActive(versionNewest)
	versCheck.Connect("toggled", func() {
		versionNewest = versCheck.GetActive()
	})
	paDirsCheck, _ := gtk.CheckButtonNewWithLabel("Create .home and .config directories for AppImages")
	paDirsCheck.SetActive(paDirs)
	paDirsCheck.Connect("toggled", func() {
		paDirs = paDirsCheck.GetActive()
	})
	gnrl.Add(wineLbl)
	gnrl.Add(dlWine)
	gnrl.Add(pthdCheck)
	gnrl.Add(wineCheck)
	gnrl.Add(versCheck)
	gnrl.Add(paDirsCheck)
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
	info, _ := gtk.ButtonNewWithLabel("Info")
	info.Connect("clicked", func() {
		infoBox, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
		infoBox.SetTransientFor(parent)
		infoBox.SetDefaultSize(300, 80)
		infoBox.SetName("common.sh info")
		infoBox.SetPosition(gtk.WIN_POS_CENTER_ON_PARENT)
		box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
		infolbl, _ := gtk.LabelNew(commonHelp)
		infolbl.SetLineWrap(true)
		infolbl.SetSizeRequest(200, 50)
		box.Add(infolbl)
		infoBox.Add(box)
		infoBox.ShowAll()
		infoBox.Show()
	})
	svBox.Add(sv)
	svBox.Add(cnl)
	svBox.Add(info)
	com.Add(comScrl)
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
		in, _ := gtk.ButtonNewWithLabel("Info")
		in.Connect("clicked", func() {
			infoBox, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
			infoBox.SetTransientFor(parent)
			infoBox.SetDefaultSize(300, 80)
			infoBox.SetName("common.sh info")
			infoBox.SetPosition(gtk.WIN_POS_CENTER_ON_PARENT)
			box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
			infolbl, _ := gtk.LabelNew(commonHelp)
			infolbl.SetLineWrap(true)
			infolbl.SetSizeRequest(200, 50)
			box.Add(infolbl)
			infoBox.Add(box)
			infoBox.ShowAll()
			infoBox.Show()
		})
		mkCom.Show()
		com.Add(mkCom)
		com.Add(in)
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
