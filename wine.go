package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mholt/archiver/v3"
)

const (
	wineURL = "https://www.playonlinux.com/wine/binaries/phoenicis/staging-linux-amd64/PlayOnLinux-wine-5.20-staging-linux-amd64.tar.gz"
)

func downloadWine(parent *gtk.Window, cb chan bool) {
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.SetTransientFor(parent)
	win.SetDestroyWithParent(true)
	win.Connect("destroy", func() {
		parent.SetSensitive(true)
	})
	parent.SetSensitive(false)
	spin, _ := gtk.SpinnerNew()
	spin.Start()
	txt, _ := gtk.LabelNew("Downloading Wine")
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
	go func(win *gtk.Window, txt *gtk.Label) {
		defer win.Close()
		wineTar, err := os.Create("PortableApps/LinuxPACom/wine5.20.tar.bz2")
		if err != nil {
			fmt.Println(err)
			cb <- false
			return
		}
		defer wineTar.Close()
		check := http.Client{
			CheckRedirect: func(r *http.Request, _ []*http.Request) error {
				r.URL.Opaque = r.URL.Path
				return nil
			},
		}
		resp, err := check.Get(wineURL)
		if err != nil {
			fmt.Println(err)
			cb <- false
			return
		}
		os.RemoveAll("PortableApps/LinuxPACom/Wine")
		defer resp.Body.Close()
		_, err = io.Copy(wineTar, resp.Body)
		if err != nil {
			fmt.Println(err)
			cb <- false
			return
		}
		txt.SetText("Extracting Wine")
		err = archiver.DefaultTarBz2.Unarchive("PortableApps/LinuxPACom/wine2.5.tar.bz2", "PortableApps/LinuxPACom/Wine")
		if err != nil {
			fmt.Println(err)
			cb <- false
			return
		}
		fil, err := os.Open("PortableApps/LinuxPACom/common.sh")
		if err != nil {
			fmt.Println(err)
			cb <- false
			return
		}
		tmp, err := ioutil.ReadAll(fil)
		if err != nil {
			fmt.Println(err)
			cb <- false
			return
		}
		if !strings.Contains(string(tmp), "export PATH=$PWD/PortableApps/LinuxPACom/Wine/wineversion/2.5/bin:$PATH") {
			tmp = append(tmp, []byte("\nexport PATH=$PWD/PortableApps/LinuxPACom/Wine/wineversion/2.5/bin:$PATH")...)
			ioutil.WriteFile("PortableApps/LinuxPACom/common.sh", tmp, 0777)
			if err != nil {
				fmt.Println(err)
				cb <- false
				return
			}
		}
		cb <- true
		return
	}(win, txt)
}
