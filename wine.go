package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mholt/archiver"
)

const (
	wineURL = "https://www.playonlinux.com/wine/binaries/linux-amd64/PlayOnLinux-wine-2.5-linux-amd64.pol"
)

func downloadWine(parent *gtk.Window, cb chan bool) {
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.SetTransientFor(parent)
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
	go func(win *gtk.Window) {
		defer win.Close()
		wineTar, err := os.Create("PortableApps/LinuxPACom/wine2.5.tar.bz2")
		if err != nil {
			fmt.Println(err)
			cb <- false
			return
		}
		wineTar.Chmod(0777)
		defer wineTar.Close()
		check := http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
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
		defer resp.Body.Close()
		_, err = io.Copy(wineTar, resp.Body)
		if err != nil {
			fmt.Println(err)
			cb <- false
			return
		}
		txt.SetText("Extracting Wine")
		os.RemoveAll("PortableApps/LinuxPACom/Wine")
		err = archiver.TarBz2.Open("PortableApps/LinuxPACom/wine2.5.tar.bz2", "PortableApps/LinuxPACom/Wine")
		if err != nil {
			fmt.Println(err)
			cb <- false
			return
		}
		os.Remove("PortableApps/LinuxPACom/wine2.5.tar.bz2")
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
		fmt.Println("Hello")
		if !strings.Contains(string(tmp), "export PATH=$PWD/PortableApps/LinuxPACom/Wine/wineversion/2.5/bin:$PATH") {
			tmp = append(tmp, []byte("\nexport PATH=$PWD/PortableApps/LinuxPACom/Wine/wineversion/2.5/bin:$PATH")...)
			ioutil.WriteFile("PortableApps/LinuxPACom/common.sh", tmp, 0777)
			fmt.Println("Hello2")
			if err != nil {
				fmt.Println(err)
				cb <- false
				return
			}
		}
		fmt.Println("HelloT")
		cb <- true
		return
	}(win)
}
