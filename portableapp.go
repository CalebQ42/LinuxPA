package main

import (
	"io"
	"io/fs"

	"github.com/probonopd/go-appimage/src/goappimage"
)

type portableApp struct {
	root          string
	name          string
	icon          []byte
	execs         []string
	appimageExecs []bool
	wineExecs     []string
	appimage      bool
}

func processPortableApp(root fs.DirEntry) (*portableApp, error) {
	var pa portableApp
	pa.root = root.Name()
	if !root.IsDir() {
		pa.appimage = true
		appimage, err := goappimage.NewAppImage("PortableApps/" + pa.root)
		if err != nil {
			return nil, nil
		}
		pa.name = appimage.Name
		var iconRdr io.ReadCloser
		iconRdr, _, err = appimage.Icon()
		if err == nil {
			pa.icon, _ = io.ReadAll(iconRdr)
			iconRdr.Close()
		}
		return &pa, nil
	}
	pa.name = root.Name()
	return &pa, nil
}
