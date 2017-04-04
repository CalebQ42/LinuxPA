package main

import (
	"io"
	"net/http"
	"os"

	"github.com/mholt/archiver"
)

const (
	wineURL = "https://www.playonlinux.com/wine/binaries/linux-amd64/PlayOnLinux-wine-2.5-linux-amd64.pol"
)

func downloadWine() (bool, error) {
	wineTar, err := os.Create("PortableApps/LinuxPACom/wine2.5.tar.bz2")
	if err != nil {
		return false, err
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
		return false, err
	}
	defer resp.Body.Close()
	_, err = io.Copy(wineTar, resp.Body)
	if err != nil {
		return false, err
	}
	err = archiver.TarBz2.Open("PortableApps/LinuxPACom/wine2.5.tar.bz2", "PortableApps/LinuxPACom/Wine")
	if err != nil {
		return false, err
	}
	return true, nil
}
