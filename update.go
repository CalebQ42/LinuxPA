package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	versionURL  = "https://www.dropbox.com/s/a0xizzo0a4vsfqt/Version?dl=1"
	downloadURL = "https://github.com/CalebQ42/LinuxPA/releases/download/vXXX/LinuxPA"
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
	defer versionFile.Close()
	check := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	response, err := check.Get(versionURL)
	defer response.Body.Close()
	if err != nil {
		return false
	}
	_, err = io.Copy(versionFile, response.Body)
	if err != nil {
		return false
	}
	return true
}

func getVersionFileInfo() string {
	fil, err := os.Open("PortableApps/LinuxPACom/Version")
	if err != nil {
		return "Error!"
	}
	rdr := bufio.NewReader(fil)
	out, _, _ := rdr.ReadLine()
	return string(out)
}

func checkForUpdate(new string) bool {
	curSlice := strings.Split(version, ".")
	newSlice := strings.Split(new, ".")
	curNums := make([]int, 4)
	newNums := make([]int, 4)
	for i, v := range curSlice {
		num, err := strconv.Atoi(v)
		if err == nil {
			curNums[i] = num
		}
		num, err = strconv.Atoi(newSlice[i])
		if err == nil {
			newNums[i] = num
		}
		if newNums[i] > curNums[i] {
			return true
		}
	}
	return false
}

func downloadUpdate(newVersion string) bool {
	url := strings.Replace(downloadURL, "XXX", newVersion, -1)
	err := os.Rename("LinuxPA", ".LinuxPA.old")
	if err != nil {
		fmt.Println(err)
		return false
	}
	fil, err := os.Create("LinuxPA")
	defer fil.Close()
	if err != nil {
		os.Rename(".LinuxPA.old", "LinuxPA")
		return false
	}
	check := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	re, err := check.Get(url)
	defer re.Body.Close()
	if err != nil {
		return false
	}
	_, err = io.Copy(fil, re.Body)
	if err != nil {
		return false
	}
	return true
}
