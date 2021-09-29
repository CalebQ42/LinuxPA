package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"path"
	"sort"

	"github.com/probonopd/go-appimage/src/goappimage"
)

type app struct {
	name     string
	path     string
	category string
	icon     []byte
	execs    []exe
	appimage bool
}

func ProcessApp(dir string) (a app, err error) {
	dirFil, err := os.Open(dir)
	if err != nil {
		return
	}
	a.path = dir
	stat, _ := dirFil.Stat()
	if !stat.IsDir() {
		var e exe
		e, err = ProcessExe(dir)
		if err != nil {
			return
		}
		if !e.appimage {
			return a, errors.New("file not in application folder and not AppImage")
		}
		a.appimage = true
		a.execs = append(a.execs, e)
		ai, _ := goappimage.NewAppImage(dir)
		a.name = ai.Name
		a.category = ai.Desktop.Section("Desktop Entry").Key("Categories").String()
		var iconRdr io.ReadCloser
		iconRdr, _, err = ai.Icon()
		if err != nil {
			iconRdr, err = ai.Thumbnail()
			if err != nil {
				err = nil
				return
			}
		}
		defer iconRdr.Close()
		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, iconRdr)
		if err != nil {
			err = nil
			if verbose {
				log.Println("Can't get icon for", a.name)
			}
			return
		}
		a.icon = buf.Bytes()
		return
	}
	dirs, err := dirFil.ReadDir(0)
	if err != nil {
		return
	}
	var e exe
	for i := range dirs {
		if dirs[i].IsDir() {
			continue
		}
		e, err = ProcessExe(path.Join(dir, dirs[i].Name()))
		if err != nil {
			if verbose {
				log.Println(path.Join(dir, dirs[i].Name()), "is not an executable. Ignoring...")
			}
			continue
		}
		if !allowWine && e.IsWine() {
			a.execs = append(a.execs, e)
		}
	}
	if len(a.execs) == 0 {
		return a, errors.New("application folder contains no executables")
	}
	//TODO: Order execs in order of importance.
	//TODO: Process information from PortableApps.com format and, if not there, any AppImages present.
	return
}

func ProcessAllApps() {
	dirs, err := os.ReadDir("PortableApps")
	if err != nil {
		log.Fatal(err)
	}
	for i := range dirs {
		var a app
		a, err = ProcessApp("PortableApps/" + dirs[i].Name())
		if err != nil {
			if verbose {
				log.Println("Error while processing", dirs[i].Name(), ":", err)
				log.Println("Ignoring...")
			}
			continue
		}
		apps = append(apps, a)
	}
	sort.Slice(apps, func(i, j int) bool {
		return apps[i].name < apps[j].name
	})
}
