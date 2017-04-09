package appimg

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

func downloadApp(parent *gtk.Window, ap appimg) {
	parent.SetSensitive(false)
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.SetTransientFor(parent)
	win.Connect("destroy", func() {
		parent.SetSensitive(true)
	})
	spn, _ := gtk.SpinnerNew()
	spn.Start()
	lbl, _ := gtk.LabelNew("Downloading " + ap.name + "...")
	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	box.SetMarginStart(10)
	box.SetMarginEnd(10)
	box.SetMarginTop(10)
	box.SetMarginBottom(10)
	box.Add(spn)
	box.Add(lbl)
	win.Add(box)
	win.SetPosition(gtk.WIN_POS_CENTER_ON_PARENT)
	win.ShowAll()
	win.Show()
	go func(win *gtk.Window, ap appimg) {
		defer win.Close()
		check := http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				r.URL.Opaque = r.URL.Path
				return nil
			},
		}
		resp, err := check.Get(urlBase + ap.name)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		name := strings.Split(ap.name, "-")[0]
		var foldName string
		if _, err = os.Open("PortableApps/" + name + "Portable"); err == nil {
			foldName = "PortableApps/" + name + "Portable"
		} else if _, err = os.Open("PortableApps/" + name); err == nil {
			foldName = "PortableApps/" + name
		} else {
			os.Mkdir("PortableApps/"+name, 0777)
			foldName = "PortableApps/" + name
		}
		fil, err := os.Create(foldName + "/" + ap.name)
		if err != nil {
			fmt.Println(err)
			return
		}
		io.Copy(fil, resp.Body)
		_ = fil.Chmod(0777)
	}(win, ap)
}
