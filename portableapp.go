package main

import (
	"errors"
	"io/fs"
	"log"
	"os"
)

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
		return &pa, nil
	}
	pa.name = root.Name()
	dirs, err := os.ReadDir(pa.root)
	if err != nil {
		return nil, err
	}
	for i := range dirs {
		var e exe
		e, err = ParseExe(pa.name + "/" + dirs[i].Name())
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

func (pa *portableApp) fyneTreeStrings() (out []string) {
	for i := range pa.execs {
		out = append(out, pa.execs[i].String())
	}
	return
}
