package main

import (
	"fmt"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var catIcons = []string{"accessories", "development", "engineering", "games", "graphics", "internet", "multimedia", "office", "other", "science", "system", "utilities"}

func ui(win *gtk.Window) {
	header, _ := gtk.HeaderBarNew()
	header.SetShowCloseButton(true)
	header.SetTitle("LinuxPA")
	header.SetSubtitle("PortableApps.com type launcher")
	settings, _ := gtk.ButtonNewFromIconName("applications-system", gtk.ICON_SIZE_SMALL_TOOLBAR)
	settings.Connect("clicked", func() {
		settingsUI()
	})
	settings.SetTooltipText("Settings (Coming Soon!)")
	header.PackStart(settings)
	win.SetTitlebar(header)
	topLvl, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	lrBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	// catList, _ := gtk.ListBoxNew()
	// catList.SetActivateOnSingleClick(true)
	store, _ := gtk.TreeStoreNew(glib.TYPE_OBJECT, glib.TYPE_STRING)
	appsList, _ := gtk.TreeViewNewWithModel(store)
	render, _ := gtk.CellRendererPixbufNew()
	pixColumn, _ := gtk.TreeViewColumnNewWithAttribute("", render, "pixbuf", 0)
	txtRender, _ := gtk.CellRendererTextNew()
	txtColumn, _ := gtk.TreeViewColumnNewWithAttribute("", txtRender, "text", 1)
	appsList.AppendColumn(pixColumn)
	appsList.AppendColumn(txtColumn)
	appsList.SetHeadersVisible(false)
	// catList.SetHExpand(true)
	// catList.SetVExpand(true)
	appsList.SetHExpand(true)
	appsList.SetVExpand(true)
	// vScrollCat, _ := gtk.AdjustmentNew(0, 0, 0, 0, 0, 0)
	// hScrollCat, _ := gtk.AdjustmentNew(0, 0, 0, 0, 0, 0)
	vScrollApp, _ := gtk.AdjustmentNew(0, 0, 0, 0, 0, 0)
	hScrollApp, _ := gtk.AdjustmentNew(0, 0, 0, 0, 0, 0)
	// catScrl, _ := gtk.ScrolledWindowNew(hScrollCat, vScrollCat)
	// catScrl.Add(catList)
	// catScrl.SetSizeRequest(170, 500)
	appScrl, _ := gtk.ScrolledWindowNew(hScrollApp, vScrollApp)
	appScrl.Add(appsList)
	appScrl.SetSizeRequest(300, 500)
	// lrBox.Add(catScrl)
	lrBox.Add(appScrl)
	botBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 2)
	wineCheck, _ := gtk.CheckButtonNewWithLabel("Show Windows apps (Wine)")
	if !wineAvail {
		wineCheck.SetSensitive(false)
		wineCheck.SetTooltipText("Download wine to run windows apps")
	}
	wineCheck.SetActive(wine)
	wineCheck.Connect("toggled", func() {
		wine = wineCheck.GetActive()
		store.Clear()
		getTreeIters(store)
	})
	botBox.Add(wineCheck)
	topLvl.Add(lrBox)
	topLvl.PackEnd(botBox, false, true, 0)
	win.Add(topLvl)
	appsList.Connect("row-activated", func() {
		selec, _ := appsList.GetSelection()
		_, it, ok := selec.GetSelected()
		if ok {
			pth, _ := store.GetPath(it)
			ind := pth.GetIndices()
			if len(ind) == 2 {
				if wine {
					app := master[cats[ind[0]]][ind[1]]
					app.launch()
				} else {
					app := linmaster[lin[ind[0]]][ind[1]]
					app.launch()
				}
			} else if len(ind) == 3 {
				if wine {
					app := master[cats[ind[0]]][ind[1]]
					app.launchSub(ind[2])
				} else {
					app := linmaster[lin[ind[0]]][ind[1]]
					app.launchSub(ind[2])
				}
			}
		}
	})
}

func getTreeIters(store *gtk.TreeStore) (out []*gtk.TreeIter) {
	if wine {
		for _, v := range cats {
			it := store.Append(nil)
			if contains(catIcons, strings.ToLower(v)) {
				img, _ := gtk.ImageNewFromIconName("applications-"+strings.ToLower(v), gtk.ICON_SIZE_BUTTON)
				buf, _ := img.GetPixbuf().ScaleSimple(32, 32, gdk.INTERP_BILINEAR)
				store.SetValue(it, 0, buf)
			} else {
				img, _ := gtk.ImageNewFromIconName("applications-other", gtk.ICON_SIZE_BUTTON)
				buf, _ := img.GetPixbuf().ScaleSimple(32, 32, gdk.INTERP_BILINEAR)
				store.SetValue(it, 0, buf)
			}
			store.SetValue(it, 1, v)
			for _, v := range master[v] {
				v.getTreeIter(store, it)
			}
			out = append(out, it)
		}
	} else {
		for _, v := range lin {
			fmt.Println(v)
			it := store.Append(nil)
			// if contains(catIcons, strings.ToLower(v)) {
			// 	img, _ := gtk.ImageNewFromIconName("applications-"+strings.ToLower(v), gtk.ICON_SIZE_BUTTON)
			// 	buf := img.GetPixbuf()
			// 	store.SetValue(it, 0, buf)
			// } else {
			// 	img, _ := gtk.ImageNewFromIconName("applications-other", gtk.ICON_SIZE_BUTTON)
			// 	buf, _ := img.GetPixbuf().ScaleSimple(32, 32, gdk.INTERP_BILINEAR)
			// 	store.SetValue(it, 0, buf)
			// }
			store.SetValue(it, 1, v)
			for _, v := range linmaster[v] {
				v.getTreeIter(store, it)
			}
			out = append(out, it)
		}
	}
	return
}
