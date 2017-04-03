package main

import (
	"io"
	"net/http"
	"os"
)

var (
	updateNeeded = false
)

const (
	versionURL = "https://www.dropbox.com/s/a0xizzo0a4vsfqt/Version?dl=1"
)

//Thanks to https://www.socketloop.com/tutorials/golang-download-file-example
//For some of the code

//Returns if success
func versionDL() bool {
	versionFile, err := os.Open("PortableApps/LinuxPACom/Version")
	if err != nil {
		versionFile, err = os.Create("PortableApps/LinuxPACom/Version")
		if err != nil {
			return false
		}
	}
	check := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	response, err := check.Get(versionURL)
	if err != nil {
		return false
	}
	_, err = io.Copy(versionFile, response.Body)
	if err != nil {
		return false
	}
	return true
}
