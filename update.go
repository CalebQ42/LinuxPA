package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func checkUpdate() {
	resp, err := http.DefaultClient.Get("https://darkstorm.tech/LinuxPA/Version")
	if err != nil {
		if p.Verbose {
			log.Println("Error while checking update information:", err)
		}
		return
	}
	var ver string
	r := bufio.NewReader(resp.Body)
	if p.Beta {
		r.ReadString('\n')
		ver, err = r.ReadString('\n')
	} else {
		ver, err = r.ReadString('\n')
	}
	resp.Body.Close()
	if err == io.EOF {
		err = nil
	} else if err != nil {
		if p.Verbose {
			log.Println("Error while checking update information:", err)
		}
		return
	}
	var newVer string
	ver = strings.Trim(ver, "\n ")
	newSplit := strings.Split(ver, ".")
	curSplit := strings.Split(version, ".")
	var new, cur int
	for i := range newSplit {
		new, err = strconv.Atoi(newSplit[i])
		if err != nil {
			if p.Verbose {
				log.Println("Error while checking update information:", err)
			}
			return
		}
		cur, err = strconv.Atoi(curSplit[i])
		if err != nil {
			if p.Verbose {
				log.Println("Error while checking update information:", err)
			}
			return
		}
		if new > cur {
			newVer = ver
			break
		}
	}
	if newVer == "" {
		if p.Verbose {
			log.Println("You are on the newest version")
		}
		return
	}
	resp, err = http.DefaultClient.Get("https://darkstorm.tech/LinuxPA/changelog")
	if err != nil {
		if p.Verbose {
			log.Println("Error while getting changelog:", err)
		}
		return
	}
	dat, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		if p.Verbose {
			log.Println("Error while getting changelog:", err)
		}
		return
	}
	showChangelogUI(newVer, string(dat))
}

//TODO: Add args to make UI
func showChangelogUI(ver string, log string) {
	_ = log
	//showChangelogUI(log)
	//onAccept(func (){
	update(ver)
	//})
}

//TOOD: Add args to close all UI
func update(ver string) error {
	resp, err := http.DefaultClient.Get("https://github.com/CalebQ42/LinuxPA/releases/download/v" + ver + "/LinuxPA")
	if err != nil {
		if p.Verbose {
			log.Println("Can't download update:", err)
		}
		return err
	}
	defer resp.Body.Close()
	err = os.Rename("LinuxPA", ".LinuxPA.bak")
	if err != nil {
		if p.Verbose {
			log.Println("Can't move old version, aborting update:", err)
		}
		return err
	}
	linPa, err := os.Create("LinuxPA")
	if err != nil {
		if p.Verbose {
			log.Println("Can't create new LinuxPA, aborting update:", err)
		}
		os.Remove("LinuxPA")
		os.Rename(".LinuxPA.bak", "LinuxPA")
		return err
	}
	_, err = io.Copy(linPa, resp.Body)
	if err != nil {
		if p.Verbose {
			log.Println("Can't download new LinuxPA, aborting update:", err)
		}
		os.Remove("LinuxPA")
		os.Rename(".LinuxPA.bak", "LinuxPA")
		return err
	}
	os.Remove(".LinuxPA.bak")
	//TODO: stop all UI
	e := exec.Command("LinuxPA")
	return e.Start()
}
