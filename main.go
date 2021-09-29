package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/CalebQ42/squashfs"
)

var (
	fyneApp fyne.App
)

func main() {
	testing := flag.Bool("test", false, "Download (if needed) and use test files. Uses the folder LinuxPATest")
	verb := flag.Bool("v", false, "Verbose")
	flag.Parse()
	verbose = *verb
	if *testing {
		_, err := os.Open("LinuxPATest")
		if os.IsNotExist(err) {
			setupTestDir()
			_, err = os.Open("LinuxPATest")
			if err != nil {
				log.Fatal(err)
			}
		}
		err = os.Chdir("LinuxPATest")
		if err != nil {
			log.Fatal(err)
		}
	}
	err := processApps()
	if err != nil {
		log.Println(err)
		return
	}
	buildAndStartFyneUI()
}

func setupTestDir() error {
	resp, err := http.DefaultClient.Get("https://darkstorm.tech/LinuxPATest.sfs")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	rdr, err := squashfs.NewSquashfsReaderFromReader(resp.Body)
	if err != nil {
		return err
	}
	err = rdr.ExtractTo("LinuxPATest")
	if err != nil {
		return err
	}
	return nil
}

func buildAndStartFyneUI() {
	fyneApp = app.New()
	win := fyneApp.NewWindow("LinuxPA")
	tree := buildFyneAppList()
	win.SetContent(tree)
	// tree.OpenAllBranches()
	win.Resize(fyne.NewSize(512, 512))
	win.ShowAndRun()
}

func buildFyneAppList() *widget.Tree {
	tree := widget.NewTree(
		func(id string) (out []string) { //TreeIDs\
			if id == "" {
				out = make([]string, len(apps))
				for i := range apps {
					out[i] = apps[i].name
				}
				return
			} else if strings.HasPrefix(id, "/") {
				return nil
			}
			for i := range apps {
				if id == apps[i].name {
					return apps[i].TreeIDs()
				}
			}
			return nil
		},
		func(id string) bool { //isBranch
			return !strings.HasPrefix(id, "/")
		},
		func(_ bool) fyne.CanvasObject { //CreateObject
			return widget.NewButton("", func() {})
		},
		func(id string, _ bool, obj fyne.CanvasObject) { //setupObject
			but := obj.(*widget.Button)
			but.Text = id
			but.OnTapped = func() {
				fmt.Println("YOO")
			}
		},
	)
	return tree
}
