package main

import (
	"os"
	"path/filepath"

	"github.com/probonopd/go-appimage/src/goappimage"
)

type App struct {
	icon          string
	name          string
	preferredExec string
	scripts       []string
	appImages     []*goappimage.AppImage
	otherExecs    []string
	winExecs      []string
}

func ProcessApp(fold *os.File) (a *App, err error) {
	a = new(App)
	dirs, err := fold.Readdirnames(-1)
	if err != nil {
		return
	}
	var tmpFil *os.File
	for _, d := range dirs {
		tmpFil, err = os.Open(d)
		if err != nil {
			return
		}
		s, _ := tmpFil.Stat()
		if s.IsDir() {
			continue
		}

	}
	//TODO: load previously set preferences.
	a.name = filepath.Base(fold.Name())
	return
}
