package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

const (
	versionURL       = "https://www.dropbox.com/s/a0xizzo0a4vsfqt/Version?dl=1"
	downloadURL      = "https://github.com/CalebQ42/LinuxPA/releases/download/vXXX/LinuxPA"
	changelogURL     = "https://www.dropbox.com/s/rk8ec9p14imkh03/Changelog?dl=1"
	changelogBetaURL = "https://www.dropbox.com/s/h2u34g5s8qr8sef/ChangelogBeta?dl=1"
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

func getVersionFileInfo() (stable string, beta string) {
	fil, err := os.Open("PortableApps/LinuxPACom/Version")
	if err != nil {
		return "Error!", ""
	}
	rdr := bufio.NewReader(fil)
	out, _, _ := rdr.ReadLine()
	stable = string(out)
	out, _, _ = rdr.ReadLine()
	beta = string(out)
	return
}

func changelogDL() (bool, error) {
	changelogFile, err := os.Create("PortableApps/LinuxPACom/Changelog")
	if err != nil {
		return false, err
	}
	changelogFile.Chmod(0777)
	check := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	var response *http.Response
	if betaUpdate {
		response, err = check.Get(changelogBetaURL)
	} else {
		response, err = check.Get(changelogURL)
	}
	if err != nil {
		return false, err
	}
	_, err = io.Copy(changelogFile, response.Body)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getChangelog() string {
	fil, err := os.Open("PortableApps/LinuxPACom/Changelog")
	if err != nil {
		return "Error!"
	}
	out, _ := ioutil.ReadAll(fil)
	return string(out)
}

func checkForUpdate(stable, beta string) (bool, error) {
	new := stable
	if betaUpdate {
		new = beta
	}
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

func update(win *gtk.Window, forced bool) {
	stat, err := versionDL()
	if stat {
		stable, beta := getVersionFileInfo()
		if stable != "Error!" {
			stat, err = checkForUpdate(stable, beta)
			if stat || forced {
				stat, err = changelogDL()
				if stat {
					updateWin, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
					updateWin.SetTransientFor(win)
					updateWin.SetPosition(gtk.WIN_POS_CENTER)
					updateWin.SetDefaultSize(600, 300)
					topLvl, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
					lbl, _ := gtk.LabelNew("There's a new update! Here's the changelog:")
					tagTbl, _ := gtk.TextTagTableNew()
					buf, _ := gtk.TextBufferNew(tagTbl)
					tv, _ := gtk.TextViewNewWithBuffer(buf)
					tv.SetWrapMode(gtk.WRAP_WORD)
					tv.SetEditable(false)
					buf.SetText(getChangelog())
					butBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
					upBut, _ := gtk.ButtonNewWithLabel("Update")
					upBut.Connect("clicked", func() {
						updateWin.Close()
						actuallyUpdate(win, forced)
					})
					cnlBut, _ := gtk.ButtonNewWithLabel("Cancel")
					cnlBut.Connect("clicked", func() {
						updateWin.Close()
					})
					butBox.Add(upBut)
					butBox.Add(cnlBut)
					topLvl.Add(lbl)
					topLvl.Add(tv)
					topLvl.Add(butBox)
					topLvl.SetMarginBottom(10)
					topLvl.SetMarginEnd(10)
					topLvl.SetMarginStart(10)
					topLvl.SetMarginTop(10)
					updateWin.Add(topLvl)
					updateWin.ShowAll()
					updateWin.Show()
				} else {
					fmt.Println(err)
				}
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Failed Version File Info")
		}
	} else {
		fmt.Println(err)
	}
}

func actuallyUpdate(win *gtk.Window, forced bool) {
	updateWin, _ := gtk.WindowNew(gtk.WINDOW_POPUP)
	updateWin.SetTransientFor(win)
	updateWin.SetSizeRequest(150, 50)
	topLvl, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	spin, _ := gtk.SpinnerNew()
	spin.Start()
	lbl, _ := gtk.LabelNew("Updating")
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
		defer updateWin.Close()
		stat, err := versionDL()
		if stat {
			stable, beta := getVersionFileInfo()
			if stable != "Error!" {
				stat, err = checkForUpdate(stable, beta)
				if stat || forced {
					if betaUpdate {
						downloadUpdate(beta)
					} else {
						downloadUpdate(stable)
					}
					win.Close()
					cmd := exec.Command("./LinuxPA")
					cmd.Stdin = os.Stdin
					cmd.Stdout = os.Stdout
					cmd.Start()
				} else {
					fmt.Println(err)
				}
			} else {
				fmt.Println("Failed Version File Info")
			}
		} else {
			fmt.Println(err)
		}
	}(win, updateWin)
}
