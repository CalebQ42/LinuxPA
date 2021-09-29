package main

import (
	"errors"
	"io/fs"
	"log"
	"os"

	"fyne.io/fyne/v2/widget"
)

var (
	apps []*portableApp
)

func processApps() error {
	pa, err := os.Open("PortableApps")
	if os.IsNotExist(err) {
		os.Mkdir("PortableApps", 0777)
		pa, err = os.Open("PortableApps")
	}
	if err != nil {
		return err
	}
	dirs, err := pa.ReadDir(0)
	if err != nil {
		return err
	}
	var portable *portableApp
	for i := range dirs {
		if dirs[i].Name() == "PortableApps.com" || dirs[i].Name() == "LinuxPACom" {
			continue
		}
		portable, err = processPortableApp(dirs[i])
		if err != nil {
			if verbose {
				log.Println("Error while processing", dirs[i].Name(), ":", err)
				log.Println("Ignoring")
			}
			continue
		}
		if portable != nil && (portable.appimage || len(portable.execs) > 0) {
			apps = append(apps, portable)
		}
	}
	return nil
}

type portableApp struct {
	root     string
	name     string
	icon     []byte
	execs    []exe
	appimage bool
}

func processPortableApp(root fs.DirEntry) (*portableApp, error) {
	var pa portableApp
	pa.root = "PortableApps/" + root.Name()
	if !root.IsDir() {
		e, err := ParseExe(pa.root)
		root.Name()
		if err != nil {
			return nil, err
		}
		if !e.ai {
			return nil, errors.New("file in PortableApps folder is not AppImage")
		}
		pa.appimage = e.ai
		pa.name = e.name
		//TODO get icon
		return &pa, nil
	}
	pa.name = root.Name()
	dirs, err := os.ReadDir(pa.root)
	if err != nil {
		return nil, err
	}
	for i := range dirs {
		var e exe
		e, err = ParseExe(pa.root + "/" + dirs[i].Name())
		if err != nil {
			if verbose {
				log.Println(e.path, "is not an executable:", err)
				log.Println("Ignoring")
			}
			continue
		}
		pa.execs = append(pa.execs, e)
	}
	return &pa, nil
}

func (pa portableApp) TreeIDs() (out []widget.TreeNodeID) {
	if !pa.appimage {
		out = make([]string, len(pa.execs))
		for i := range pa.execs {
			out[i] = pa.execs[i].path
		}
	}
	return
}
