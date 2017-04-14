//Package appimg is for downloading new AppImages for LinuxPA
package appimg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	urlBase = "https://dl.bintray.com/probono/AppImages/"
)

//ShowUI shows the list of possible AppImages to be downloaded in a gtk.Window
func ShowUI(newestVersionOnly bool, clsFunc func()) {
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.Connect("destroy", func() {
		clsFunc()
	})
	apps := make([]appimg, 0)
	win.SetSizeRequest(400, 400)
	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	appList, _ := gtk.ListBoxNew()
	apch := make(chan appimg)
	appList.SetHExpand(true)
	appList.SetVExpand(true)
	vScrollCat, _ := gtk.AdjustmentNew(0, 0, 0, 0, 0, 0)
	hScrollCat, _ := gtk.AdjustmentNew(0, 0, 0, 0, 0, 0)
	lst, _ := gtk.ScrolledWindowNew(hScrollCat, vScrollCat)
	lst.SetSizeRequest(170, 500)
	lst.Add(appList)
	box.Add(lst)
	win.Add(box)
	appList.Connect("row-activated", func() {
		if appList.GetSelectedRow().GetIndex() >= 0 {
			downloadApp(win, apps[appList.GetSelectedRow().GetIndex()])
		}
	})
	win.SetPosition(gtk.WIN_POS_CENTER_ON_PARENT)
	win.ShowAll()
	win.Show()
	getList(win, apch)
	go func(win *gtk.Window, apch chan appimg, list *gtk.ListBox) {
		if newestVersionOnly {
			imgs := make([]appimg, 0)
			a := make(map[string][]appimg)
			names := make([]string, 0)
			for i := range apch {
				imgs = append(imgs, i)
			}
			for i, v := range imgs {
				sp := strings.Split(v.full, "-")
				if len(sp) >= 2 {
					vers := sp[1]
					removeLetters(vers)
					imgs[i].version = vers
					imgs[i].name = sp[0]
					if _, ok := a[imgs[i].name]; !ok {
						names = append(names, imgs[i].name)
					}
					a[imgs[i].name] = append(a[imgs[i].name], imgs[i])
				}
			}
			sort.Strings(names)
			for _, name := range names {
				glib.IdleAdd(func(name string, list *gtk.ListBox, i appimg) {
					lbl, _ := gtk.LabelNew(name)
					list.Add(lbl)
					apps = append(apps, i)
					lbl.Show()
				}, name, list, a[name][compareVersions(a[name])])
			}
		} else {
			for i := range apch {
				glib.IdleAdd(func(list *gtk.ListBox, i appimg) {
					lbl, _ := gtk.LabelNew(i.full)
					list.Add(lbl)
					apps = append(apps, i)
					lbl.Show()
				}, list, i)
			}
		}
	}(win, apch, appList)
}

func getList(parent *gtk.Window, apch chan appimg) {
	win, _ := gtk.WindowNew(gtk.WINDOW_POPUP)
	win.SetTransientFor(parent)
	win.SetDestroyWithParent(true)
	win.Connect("destroy", func() {
		parent.SetSensitive(true)
	})
	parent.SetSensitive(false)
	spin, _ := gtk.SpinnerNew()
	spin.Start()
	txt, _ := gtk.LabelNew("Getting List...")
	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	box.SetMarginBottom(10)
	box.SetMarginEnd(10)
	box.SetMarginStart(10)
	box.SetMarginTop(10)
	box.Add(spin)
	box.Add(txt)
	win.Add(box)
	win.SetPosition(gtk.WIN_POS_CENTER_ON_PARENT)
	win.ShowAll()
	win.Show()
	go func(win *gtk.Window, apch chan appimg) {
		check := http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				r.URL.Opaque = r.URL.Path
				return nil
			},
		}
		resp, err := check.Get(urlBase)
		if err != nil {
			fmt.Println(err)
			close(apch)
			win.Close()
			return
		}
		defer resp.Body.Close()
		btys, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			close(apch)
			win.Close()
			return
		}
		tgs := convert(string(btys))
		for _, v := range tgs {
			if strings.HasSuffix(strings.ToLower(v.Meat), ".appimage") {
				apch <- newApp(v.Meat)
			}
		}
		close(apch)
		win.Close()
	}(win, apch)
}
