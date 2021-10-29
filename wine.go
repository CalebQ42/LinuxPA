package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/mholt/archiver/v3"
)

func getWine() error {
	resp, err := http.DefaultClient.Get("https://darkstorm.tech/LinuxPA/wine")
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}
	ver, _ := buf.ReadString('\n')
	if ver == "" {
		return errors.New("couldn't get wine version")
	}
	wineURL := strings.ReplaceAll(wineDlURL, "ARCH", "amd64")
	wineURL = strings.ReplaceAll(wineURL, "VERSION", ver)
	resp, err = http.DefaultClient.Get(wineURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	wineTmp, err := os.Create("PortableApps/LinuxPACom/wine.tar.gz")
	if err != nil {
		return err
	}
	defer func() {
		wineTmp.Close()
		os.Remove("PortableApps/LinuxPACom/wine.tar.gz")
	}()
	_, err = io.Copy(wineTmp, resp.Body)
	if err != nil {
		return err
	}
	err = archiver.Unarchive("PortableApps/LinuxPACom/wine.tar.gz", "PortableApps/LinuxPACom/Wine")
	if err != nil {
		return err
	}
	winePath = "PortableApps/LinuxPACom/Wine/bin"
	return nil
}
