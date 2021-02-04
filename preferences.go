package main

import (
	"encoding/gob"
	"encoding/json"
	"log"
	"os"
)

var (
	wine         bool
	wineAvail    bool
	portableHide bool
	betaUpdate   bool
	paDirs       = true
)

const prefLocation = "PortableApps/LinuxPA/Preferences.json"

func savePrefs() {
	os.Rename(prefLocation, prefLocation+".bak")
	defer os.Remove(prefLocation + ".bak")
	fil, err := os.Create(prefLocation)
	if err != nil {
		log.Println("Cannot create", prefLocation)
		log.Println(err)
		os.Exit(-1)
	}
	enc := json.NewEncoder(fil)
	err = enc.Encode(map[string]interface{}{
		"show wine":               wine,
		"hide portable":           portableHide,
		"home and config sandbox": paDirs,
		"beta":                    betaUpdate,
	})
}

func loadPrefs() {
	fil, err := os.Open("PortableApps/LinuxPA")
	if os.IsNotExist(err) {
		fil, err = os.Open("PortableApps/LinuxPACom")
		if err == nil {
			loadPrefsLegacy()
			os.RemoveAll("PortableApps/LinuxPACom")
		}
		err = os.MkdirAll("PortableApps/LinuxPA", os.ModePerm)
		if err != nil {
			log.Println("Error while creating PortableApps/LinuxPA:")
			log.Println(err)
			os.Exit(-1)
		}
	} else if err != nil {
		log.Println("Cannot open PortableApps/LinuxPA")
		log.Println(err)
		os.Exit(-1)
	}
	fil, err = os.Open("PortableApps/LinuxPA/Preferences.json")
	if os.IsNotExist(err) {
		//TODO: load compatible preferences from PortableApps.com Launcher
		return
	}
	dec := json.NewDecoder(fil)
	var prefs map[string]interface{}
	err = dec.Decode(&prefs)
	if err != nil {
		log.Println("Error while reading preferences")
		log.Println(err)
		os.Exit(-1)
	}
	var ok bool
	wine, ok = prefs["show wine"].(bool)
	if !ok {
		log.Println("show wine preference is of the incorrect type (should be bool)")
		wine = false
	}
	portableHide, ok = prefs["hide portable"].(bool)
	if !ok {
		log.Println("hide portable preference is of the incorrect type (should be bool)")
		portableHide = false
	}
	paDirs, ok = prefs["home and config sandbox"].(bool)
	if !ok {
		log.Println("home and config sadbox preference is of the incorrect type (should be bool)")
		paDirs = true
	}
	betaUpdate, ok = prefs["beta"].(bool)
	if !ok {
		log.Println("beta preference is of the incorrect type (should be bool)")
		betaUpdate = false
	}
}

func loadPrefsLegacy() {
	fil, err := os.Open("PortableApps/LinuxPACom/Prefs.gob")
	if err != nil {
		return
	}
	dec := gob.NewDecoder(fil)
	err = dec.Decode(&wine)
	if err != nil {
		return
	}
	err = dec.Decode(&portableHide)
	if err != nil {
		return
	}
	var unused bool
	err = dec.Decode(&unused)
	if err != nil {
		return
	}
	err = dec.Decode(&paDirs)
	if err != nil {
		return
	}
	err = dec.Decode(&betaUpdate)
	if err != nil {
		return
	}
}
