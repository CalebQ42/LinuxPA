package main

import (
	"bufio"
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
func versionDL() (bool, error) {
	versionFile, err := os.Create("PortableApps/LinuxPACom/Version")
	if err != nil {
		return false, err
	}
	versionFile.Chmod(0777)
	check := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	response, err := check.Get(versionURL)
	if err != nil {
		return false, err
	}
	_, err = io.Copy(versionFile, response.Body)
	if err != nil {
		return false, err
	}
	return true, nil
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

func checkForUpdate(new string) (bool, error) {
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
		} else {
			return false, err
		}
		if newNums[i] > curNums[i] {
			return true, nil
		}
	}
	return false, nil
}

func downloadUpdate(newVersion string) (bool, error) {
	url := strings.Replace(downloadURL, "XXX", newVersion, -1)
	err := os.Rename("LinuxPA", ".LinuxPA.old")
	if err != nil {
		return false, err
	}
	fil, err := os.Create("LinuxPA")
	fil.Chmod(0777)
	defer fil.Close()
	if err != nil {
		os.Rename(".LinuxPA.old", "LinuxPA")
		return false, err
	}
	check := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	re, err := check.Get(url)
	if err != nil {
		return false, err
	}
	defer re.Body.Close()
	_, err = io.Copy(fil, re.Body)
	if err != nil {
		return false, err
	}
	return true, nil
}
