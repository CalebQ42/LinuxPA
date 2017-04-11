package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type app struct {
	name   string
	cat    string
	appimg []string
	lin    []string
	ex     []string
	icon   *gdk.Pixbuf
	dir    string
	ini    *os.File
}

func (a *app) getTreeIter(store *gtk.TreeStore) *gtk.TreeIter {
	it := store.Append(nil)
	store.SetValue(it, 0, a.icon)
	store.SetValue(it, 1, a.name)
	if len(a.ex) > 1 {
		if wine {
			for _, v := range a.ex {
				i := store.Append(it)
				store.SetValue(i, 1, v)
			}
		} else {
			for _, v := range a.lin {
				i := store.Append(it)
				store.SetValue(i, 1, v)
			}
		}
	}
	return it
}

func (a *app) launch() {
	if len(a.ex) == 1 {
		if wine {
			var cmd *exec.Cmd
			if !contains(a.lin, a.ex[0]) {
				if comEnbld {
					cmd = exec.Command("/bin/sh", "-c", ". PortableApps/LinuxPACom/common.sh || exit 1;cd \""+a.dir+"\"; wine \""+a.ex[0]+"\"")
				} else {
					cmd = exec.Command("/bin/sh", "-c", "cd \""+a.dir+"\"; wine \""+a.ex[0]+"\"")
				}
			} else {
				if comEnbld {
					cmd = exec.Command("/bin/sh", "-c", ". PortableApps/LinuxPACom/common.sh || exit 1;cd \""+a.dir+"\"; \"./"+a.ex[0]+"\"")
				} else {
					cmd = exec.Command("/bin/sh", "-c", "cd \""+a.dir+"\"; \"./"+a.ex[0]+"\"")
				}
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Start()
		} else {
			var cmd *exec.Cmd
			if comEnbld {
				cmd = exec.Command("/bin/sh", "-c", ". PortableApps/LinuxPACom/common.sh || exit 1;cd \""+a.dir+"\"; \"./"+a.ex[0]+"\"")
			} else {
				cmd = exec.Command("/bin/sh", "-c", "cd \""+a.dir+"\"; \"./"+a.ex[0]+"\"")
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Start()
		}
	} else {
		if wine {
			var cmd *exec.Cmd
			if len(a.lin) == 0 {
				if comEnbld {
					cmd = exec.Command("/bin/sh", "-c", ". PortableApps/LinuxPACom/common.sh || exit 1;cd \""+a.dir+"\"; wine \""+a.ex[0]+"\"")
				} else {
					cmd = exec.Command("/bin/sh", "-c", "cd \""+a.dir+"\"; wine \""+a.ex[0]+"\"")
				}
			} else {
				var ind int
				for i, v := range a.lin {
					if strings.HasSuffix(v, ".sh") {
						ind = i
						break
					}
				}
				if comEnbld {
					cmd = exec.Command("/bin/sh", "-c", ". PortableApps/LinuxPACom/common.sh || exit 1;cd \""+a.dir+"\"; \"./"+a.lin[ind]+"\"")
				} else {
					cmd = exec.Command("/bin/sh", "-c", "cd \""+a.dir+"\"; \"./"+a.lin[ind]+"\"")
				}
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Start()
		} else {
			if len(a.lin) != 0 {
				var ind int
				for i, v := range a.lin {
					if strings.HasSuffix(v, ".sh") {
						ind = i
						break
					}
				}
				var cmd *exec.Cmd
				if comEnbld {
					cmd = exec.Command("/bin/sh", "-c", ". PortableApps/LinuxPACom/common.sh || exit 1;cd \""+a.dir+"\"; \"./"+a.lin[ind]+"\"")
				} else {
					cmd = exec.Command("/bin/sh", "-c", "cd \""+a.dir+"\"; \"./"+a.lin[ind]+"\"")
				}
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Start()
			}
		}
	}
}

func (a *app) launchSub(sub int) {
	if wine {
		var cmd *exec.Cmd
		if !contains(a.lin, a.ex[sub]) {
			if comEnbld {
				cmd = exec.Command("/bin/sh", "-c", ". PortableApps/LinuxPACom/common.sh || exit 1;cd \""+a.dir+"\"; wine \""+a.ex[sub]+"\"")
			} else {
				cmd = exec.Command("/bin/sh", "-c", "cd \""+a.dir+"\"; wine \""+a.ex[sub]+"\"")
			}
		} else {
			if comEnbld {
				cmd = exec.Command("/bin/sh", "-c", ". PortableApps/LinuxPACom/common.sh || exit 1;cd \""+a.dir+"\"; \"./"+a.ex[sub]+"\"")
			} else {
				cmd = exec.Command("/bin/sh", "-c", "cd \""+a.dir+"\"; \"./"+a.ex[sub]+"\"")
			}
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
	} else {
		var cmd *exec.Cmd
		if comEnbld {
			cmd = exec.Command("/bin/sh", "-c", ". PortableApps/LinuxPACom/common.sh || exit 1;cd \""+a.dir+"\"; \"./"+a.ex[sub]+"\"")
		} else {
			cmd = exec.Command("/bin/sh", "-c", "cd \""+a.dir+"\"; \"./"+a.ex[sub]+"\"")
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
	}
}

func (a *app) edit(parent *gtk.Window, reload func()) {
	tmp := *a
	parent.SetSensitive(false)
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.Connect("destroy", func() {
		master = make(map[string][]app)
		linmaster = make(map[string][]app)
		cats = make([]string, 0)
		lin = make([]string, 0)
		setup()
		reload()
		parent.SetSensitive(true)
	})
	win.SetDefaultSize(400, 135)
	topLvl, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	topLvl.SetMarginStart(10)
	topLvl.SetMarginEnd(10)
	topLvl.SetMarginTop(10)
	topLvl.SetMarginBottom(10)
	top, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	img, _ := gtk.ImageNewFromPixbuf(a.icon)
	imgBut, _ := gtk.ButtonNew()
	imgBut.SetImage(img)
	imgBut.SetSizeRequest(100, 100)
	imgBut.Connect("clicked", func() {
		fil, _ := gtk.FileChooserDialogNewWith2Buttons("Select Icon", win, gtk.FILE_CHOOSER_ACTION_OPEN, "Cancel", gtk.RESPONSE_CANCEL, "Open", gtk.RESPONSE_ACCEPT)
		filter, _ := gtk.FileFilterNew()
		filter.AddPixbufFormats()
		filter.SetName("Supported Pictures")
		fil.AddFilter(filter)
		resp := fil.Run()
		if resp == int(gtk.RESPONSE_ACCEPT) {
			filename := fil.GetFilename()
			_, err := os.Open(filename)
			if err != nil {
				fmt.Println(err)
				return
			}
			pix, _ := gdk.PixbufNewFromFileAtSize(filename, 32, 32)
			tmp.icon = pix
			img.SetFromPixbuf(pix)
			fil.Close()
		}
	})
	topRt, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	nameLbl, _ := gtk.LabelNew("Name:")
	nameLbl.SetHAlign(gtk.ALIGN_START)
	txtgtbl, _ := gtk.TextTagTableNew()
	txtBuf, _ := gtk.TextBufferNew(txtgtbl)
	nameTxt, _ := gtk.TextViewNewWithBuffer(txtBuf)
	nameTxt.SetAcceptsTab(false)
	nameTxt.SetWrapMode(gtk.WRAP_CHAR)
	nameTxt.SetPixelsBelowLines(5)
	nameTxt.SetHExpand(true)
	nameTxt.SetVExpand(false)
	nameTxt.SetBorderWindowSize(gtk.TEXT_WINDOW_BOTTOM, 5)
	txtBuf.SetText(tmp.name)
	vScrollName, _ := gtk.AdjustmentNew(0, 0, 0, 0, 0, 0)
	hScrollName, _ := gtk.AdjustmentNew(0, 0, 0, 0, 0, 0)
	nameScr, _ := gtk.ScrolledWindowNew(hScrollName, vScrollName)
	nameScr.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_NEVER)
	nameScr.SetSizeRequest(300, 25)
	nameScr.SetVExpand(false)
	nameScr.Add(nameTxt)
	catLbl, _ := gtk.LabelNew("Category:")
	catLbl.SetHAlign(gtk.ALIGN_START)
	catTbl, _ := gtk.TextTagTableNew()
	catBuf, _ := gtk.TextBufferNew(catTbl)
	catTxt, _ := gtk.TextViewNewWithBuffer(catBuf)
	catBuf.SetText(tmp.cat)
	catTxt.SetAcceptsTab(false)
	catTxt.SetWrapMode(gtk.WRAP_CHAR)
	catTxt.SetPixelsBelowLines(5)
	catTxt.SetHExpand(true)
	catTxt.SetVExpand(false)
	catTxt.SetBorderWindowSize(gtk.TEXT_WINDOW_BOTTOM, 5)
	vScrollCat, _ := gtk.AdjustmentNew(0, 0, 0, 0, 0, 0)
	hScrollCat, _ := gtk.AdjustmentNew(0, 0, 0, 0, 0, 0)
	catScr, _ := gtk.ScrolledWindowNew(hScrollCat, vScrollCat)
	catScr.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_NEVER)
	catScr.SetSizeRequest(300, 25)
	catScr.SetVExpand(false)
	catScr.Add(catTxt)
	topRt.Add(nameLbl)
	topRt.Add(nameScr)
	topRt.Add(catLbl)
	topRt.Add(catScr)
	top.Add(imgBut)
	top.Add(topRt)
	bot, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	sv, _ := gtk.ButtonNewWithLabel("Save")
	sv.Connect("clicked", func() {
		tmp.name, _ = txtBuf.GetText(txtBuf.GetStartIter(), txtBuf.GetEndIter(), true)
		tmp.cat, _ = catBuf.GetText(catBuf.GetStartIter(), catBuf.GetEndIter(), true)
		tmp.makeIni()
		os.Remove(a.dir + "/appicon.png")
		tmp.icon.SavePNG(a.dir+"/appicon.png", 0)
		win.Close()
	})
	cnl, _ := gtk.ButtonNewWithLabel("Cancel")
	cnl.Connect("clicked", func() {
		win.Close()
	})
	bot.PackEnd(sv, false, false, 0)
	bot.PackEnd(cnl, false, false, 0)
	topLvl.Add(top)
	topLvl.Add(bot)
	win.Add(topLvl)
	win.ShowAll()
	win.Show()
}

func (a *app) makeIni() {
	os.Remove(a.dir + "/appinfo.ini")
	fil, err := os.Create(a.dir + "/appinfo.ini")
	if err != nil {
		return
	}
	ini := "[General]\n"
	ini += "Category=" + a.cat + "\n"
	ini += "Name=" + a.name + "\n"
	wrt := bufio.NewWriter(fil)
	wrt.WriteString(ini)
	wrt.Flush()
}
