package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/probonopd/go-appimage/src/goappimage"
	"gopkg.in/ini.v1"
)

type app struct {
	desktop    *ini.File
	name       string
	path       string
	categories []string
	icon       []byte
	execs      []exe
	appimage   bool
}

func ProcessApp(dir string) ([]*app, error) {
	var a app
	dirFil, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	a.path = dir
	stat, _ := dirFil.Stat()
	if !stat.IsDir() {
		var e exe
		e, err = ProcessExe(dir)
		if err != nil {
			return nil, err
		}
		if !e.appimage {
			return nil, errors.New("file not in application folder and not AppImage")
		}
		a.appimage = true
		a.execs = append(a.execs, e)
		ai, _ := goappimage.NewAppImage(dir)
		a.name = ai.Name
		a.categories = strings.Split(ai.Desktop.Section("Desktop Entry").Key("Categories").String(), "；")
		var iconRdr io.ReadCloser
		iconRdr, _, err = ai.Icon()
		if err != nil {
			iconRdr, err = ai.Thumbnail()
			if err != nil {
				if verbose {
					log.Println("Can't get icon for", a.name, ":", err)
				}
				return []*app{&a}, nil
			}
		}
		defer iconRdr.Close()
		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, iconRdr)
		if err != nil {
			if verbose {
				log.Println("Can't get icon for", a.name, ":", err)
			}
			return []*app{&a}, nil
		}
		a.icon = buf.Bytes()
		return []*app{&a}, nil
	}
	dirs, err := dirFil.ReadDir(0)
	if err != nil {
		return nil, err
	}
	var desktopFiles []string
	var e exe
	for i := range dirs {
		if dirs[i].IsDir() {
			continue
		}
		if strings.HasSuffix(dirs[i].Name(), ".desktop") {
			desktopFiles = append(desktopFiles, dirs[i].Name())
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
		return nil, errors.New("application folder contains no executables")
	}
	a.orderExecs()
	if len(desktopFiles) == 0 {
		return []*app{&a}, nil
	}
	var out []*app
	for _, d := range desktopFiles {
		var deskFil *os.File
		deskFil, err = os.Open(path.Join(dir, d))
		if err != nil {
			if verbose {
				log.Println("Error while opening desktop file", path.Join(dir, d), ":", err)
				log.Println("Ignoring...")
			}
			continue
		}
		var desktopApp app
		desktopApp.desktop, err = ini.Load(deskFil)
		if err != nil {
			if verbose {
				log.Println("Error while processing desktop file", path.Join(dir, d), ":", err)
				log.Println("Ignoring...")
			}
			continue
		}
		exec := desktopApp.desktop.Section("DesktopEntry").Key("Exec").String()
		if exec == "" {
			if verbose {
				log.Println("Desktop file", path.Join(dir, d), "does not have an Exec key. Ignoring...")
			}
			continue
		}
		var args string
		execSplit := strings.Split(exec, "")
		exec = execSplit[0]
		for i := 1; i < len(execSplit); i++ {
			if strings.HasPrefix(execSplit[i], "$") {
				continue
			}
			args += " " + execSplit[i]
		}
		args = strings.TrimSpace(args)
		//TODO
		desktopApp.name = desktopApp.desktop.Section("Desktop Entry").Key("Name").String()
		desktopApp.categories = strings.Split(desktopApp.desktop.Section("Desktop Entry").Key("Categories").String(), "；")

	}
	if len(out) == 0 {
		return []*app{&a}, nil
	}
	//TODO: Process information from PortableApps.com format and, if not there, any AppImages present.
	//TODO: Process any .desktop files as a new app (if there are multiple)
	return out, nil
}

func (a *app) orderExecs() {
	var scripts, appimages, linExecs, wine []exe
	for _, e := range a.execs {
		if e.script {
			scripts = append(scripts, e)
		} else if e.appimage {
			appimages = append(appimages, e)
		} else if e.IsWine() {
			wine = append(wine, e)
		} else {
			linExecs = append(linExecs, e)
		}
	}
	a.execs = append(append(append(scripts, appimages...), linExecs...), wine...)
}

func ProcessAllApps() {
	dirs, err := os.ReadDir("PortableApps")
	if err != nil {
		log.Fatal(err)
	}
	for i := range dirs {
		var a []*app
		a, err = ProcessApp("PortableApps/" + dirs[i].Name())
		if err != nil {
			if verbose {
				log.Println("Error while processing", dirs[i].Name(), ":", err)
				log.Println("Ignoring...")
			}
			continue
		}
		apps = append(apps, a...)
	}
	sort.Slice(apps, func(i, j int) bool {
		return apps[i].name < apps[j].name
	})
}
