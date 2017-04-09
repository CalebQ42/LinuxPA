package main

import (
	"fmt"

	"github.com/CalebQ42/LinuxPA/appimg"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func ui(win *gtk.Window) {
	ls := getCatRows()
	var treeApps []*gtk.TreeIter
	header, _ := gtk.HeaderBarNew()
	header.SetShowCloseButton(true)
	header.SetTitle("LinuxPA")
	header.SetSubtitle("PortableApps.com type launcher")
	settings, _ := gtk.ButtonNewFromIconName("applications-system", gtk.ICON_SIZE_SMALL_TOOLBAR)
	settings.Connect("clicked", func() {
		settingsUI()
	})
	settings.SetTooltipText("Settings")
	dnl, _ := gtk.ButtonNewFromIconName("emblem-downloads", gtk.ICON_SIZE_SMALL_TOOLBAR)
	dnl.SetTooltipText("Download Apps")
	header.PackStart(settings)
	header.PackEnd(dnl)
	win.SetTitlebar(header)
	topLvl, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	lrBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	catList, _ := gtk.ListBoxNew()
	catList.SetActivateOnSingleClick(true)
	store, _ := gtk.TreeStoreNew(glib.TYPE_OBJECT, glib.TYPE_STRING)
	appsList, _ := gtk.TreeViewNewWithModel(store)
	render, _ := gtk.CellRendererPixbufNew()
	pixColumn, _ := gtk.TreeViewColumnNewWithAttribute("", render, "pixbuf", 0)
	txtRender, _ := gtk.CellRendererTextNew()
	txtColumn, _ := gtk.TreeViewColumnNewWithAttribute("", txtRender, "text", 1)
	appsList.AppendColumn(pixColumn)
	appsList.AppendColumn(txtColumn)
	appsList.SetHeadersVisible(false)
	catList.SetHExpand(true)
	catList.SetVExpand(true)
	appsList.SetHExpand(true)
	appsList.SetVExpand(true)
	vScrollCat, _ := gtk.AdjustmentNew(0, 0, 0, 0, 0, 0)
	hScrollCat, _ := gtk.AdjustmentNew(0, 0, 0, 0, 0, 0)
	vScrollApp, _ := gtk.AdjustmentNew(0, 0, 0, 0, 0, 0)
	hScrollApp, _ := gtk.AdjustmentNew(0, 0, 0, 0, 0, 0)
	catScrl, _ := gtk.ScrolledWindowNew(hScrollCat, vScrollCat)
	catScrl.Add(catList)
	catScrl.SetSizeRequest(170, 500)
	appScrl, _ := gtk.ScrolledWindowNew(hScrollApp, vScrollApp)
	appScrl.Add(appsList)
	appScrl.SetSizeRequest(300, 500)
	lrBox.Add(catScrl)
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
		for i := range ls {
			fmt.Println(len(ls) - i)
			catList.Remove(catList.GetRowAtIndex(len(ls) - i - 1))
		}
		ls = getCatRows()
		for _, v := range ls {
			catList.Add(v)
		}
		catList.ShowAll()
	})
	botBox.Add(wineCheck)
	topLvl.Add(lrBox)
	topLvl.PackEnd(botBox, false, true, 0)
	win.Add(topLvl)
	for _, v := range ls {
		catList.Add(v)
	}
	catList.Connect("row-selected", func() {
		store.Clear()
		if catList.GetSelectedRow().GetIndex() >= 0 {
			treeApps = make([]*gtk.TreeIter, 0)
			if wine {
				apps := master[cats[catList.GetSelectedRow().GetIndex()]]
				for _, v := range apps {
					treeApps = append(treeApps, v.getTreeIter(store))
				}
			} else {
				apps := linmaster[lin[catList.GetSelectedRow().GetIndex()]]
				for _, v := range apps {
					treeApps = append(treeApps, v.getTreeIter(store))
				}
			}
		}
	})
	appsList.Connect("row-activated", func() {
		selec, _ := appsList.GetSelection()
		_, it, ok := selec.GetSelected()
		if ok {
			pth, _ := store.GetPath(it)
			ind := pth.GetIndices()
			if len(ind) == 1 {
				if wine {
					app := master[cats[catList.GetSelectedRow().GetIndex()]][ind[0]]
					app.launch()
				} else {
					app := linmaster[lin[catList.GetSelectedRow().GetIndex()]][ind[0]]
					app.launch()
				}
			} else if len(ind) == 2 {
				if wine {
					app := master[cats[catList.GetSelectedRow().GetIndex()]][ind[0]]
					app.launchSub(ind[1])
				} else {
					app := linmaster[lin[catList.GetSelectedRow().GetIndex()]][ind[0]]
					app.launchSub(ind[1])
				}
			}
		}
	})
	dnl.Connect("clicked", func() {
		appimg.ShowUI(func() {
			store.Clear()
			ls = getCatRows()
			for _, v := range ls {
				catList.Add(v)
			}
		})
	})
}

func getCatRows() (out []*gtk.Label) {
	if wine {
		for _, v := range cats {
			txt, _ := gtk.LabelNew(v)
			out = append(out, txt)
		}
	} else {
		for _, v := range lin {
			txt, _ := gtk.LabelNew(v)
			out = append(out, txt)
		}
	}
	return
}
