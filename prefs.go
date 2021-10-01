package main

import (
	"encoding/gob"
	"encoding/json"
	"os"
)

const (
	prefsLoc = "PortableApps/LinuxPACom/preferences.json"
	oldPrefs = "PortableApps/LinuxPACom/Prefs.gob"
	commonSh = "PortableApps/LinuxPACom/common.sh"
)

var prefs preferences

type preferences struct {
	verbose         bool
	hideWine        bool
	showPortable    bool
	fromRoot        bool
	disableCommonSh bool
	beta            bool
	apimgDirs       bool
}

func loadPrefs() error {
	jsonFil, err := os.Open(prefsLoc)
	if err != nil {
		return err
	}
	defer jsonFil.Close()
	d := json.NewDecoder(jsonFil)
	err = d.Decode(&prefs)
	return err
}

func savePrefs() error {
	os.Rename(prefsLoc, prefsLoc+".bak")
	jsonFil, err := os.Create(prefsLoc)
	if err != nil {
		return err
	}
	defer jsonFil.Close()
	e := json.NewEncoder(jsonFil)
	err = e.Encode(&prefs)
	if err != nil {
		os.Remove(prefsLoc)
		os.Rename(prefsLoc+".bak", prefsLoc)
	} else {
		os.Remove(prefsLoc + ".bak")
	}
	return err
}

func loadOldPrefs() error {
	prefsFil, err := os.Open(oldPrefs)
	if err != nil {
		return err
	}
	defer prefsFil.Close()
	d := gob.NewDecoder(prefsFil)
	tmpBool := new(bool)
	err = d.Decode(tmpBool)
	if err != nil {
		return nil
	}
	prefs.hideWine = !*tmpBool
	err = d.Decode(tmpBool)
	if err != nil {
		return nil
	}
	prefs.showPortable = !*tmpBool
	err = d.Decode(tmpBool)
	if err != nil {
		return nil
	}
	//versionNewest
	err = d.Decode(tmpBool)
	if err != nil {
		return nil
	}
	prefs.apimgDirs = *tmpBool
	err = d.Decode(tmpBool)
	if err != nil {
		return nil
	}
	prefs.beta = !*tmpBool
	return nil
}
