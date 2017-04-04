package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
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
		} else if curNums[i] > newNums[i] {
			return false, nil
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

func update(win *gtk.Window) {
	updateWin, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	updateWin.SetTransientFor(win)
	topLvl, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	spin, _ := gtk.SpinnerNew()
	spin.Start()
	lbl, _ := gtk.LabelNew("Checking for updates")
	topLvl.Add(spin)
	topLvl.Add(lbl)
	topLvl.SetMarginBottom(10)
	topLvl.SetMarginEnd(10)
	topLvl.SetMarginStart(10)
	topLvl.SetMarginTop(10)
	updateWin.SetPosition(gtk.WIN_POS_CENTER_ON_PARENT)
	updateWin.Add(topLvl)
	updateWin.ShowAll()
	updateWin.Show()
	go func(win, updateWin *gtk.Window) {
		stat, err := versionDL()
		if stat {
			res := getVersionFileInfo()
			if res != "Error!" {
				stat, err = checkForUpdate(res)
				if stat {
					lbl.SetText("Updating!")
					downloadUpdate(res)
					updateWin.Close()
					win.Close()
					cmd := exec.Command("./LinuxPA")
					cmd.Stdin = os.Stdin
					cmd.Stdout = os.Stdout
					cmd.Start()
				} else {
					fmt.Println(err)
					updateWin.Close()
				}
			} else {
				fmt.Println("Failed Version File Info")
				updateWin.Close()
			}
		} else {
			fmt.Println(err)
			updateWin.Close()
		}
	}(win, updateWin)
}
