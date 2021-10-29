package prefs

import (
	"encoding/gob"
	"encoding/json"
	"os"
)

const (
	prefsLoc = "PortableApps/LinuxPACom/Prefs.json"
	oldPrefs = "PortableApps/LinuxPACom/Prefs.gob"
)

type Prefs struct {
	Verbose         bool
	HideWine        bool
	ShowPortable    bool
	FromRoot        bool
	DisableCommonSh bool
	Beta            bool
	AppImageDirs    bool
}

func LoadPrefs() (p *Prefs, err error) {
	jsonFil, err := os.Open(prefsLoc + ".bak")
	if err != nil {
		jsonFil, err = os.Open(prefsLoc)
		if err != nil {
			return
		}
	}
	p = new(Prefs)
	defer jsonFil.Close()
	d := json.NewDecoder(jsonFil)
	err = d.Decode(p)
	return
}

func (p *Prefs) SavePrefs() error {
	os.Rename(prefsLoc, prefsLoc+".bak")
	jsonFil, err := os.Create(prefsLoc)
	if err != nil {
		return err
	}
	defer jsonFil.Close()
	e := json.NewEncoder(jsonFil)
	err = e.Encode(&p)
	if err != nil {
		os.Remove(prefsLoc)
		os.Rename(prefsLoc+".bak", prefsLoc)
	} else {
		os.Remove(prefsLoc + ".bak")
	}
	return err
}

func LoadOldPrefs() (p *Prefs, err error) {
	prefsFil, err := os.Open(oldPrefs)
	if err != nil {
		return
	}
	p = new(Prefs)
	defer prefsFil.Close()
	d := gob.NewDecoder(prefsFil)
	tmpBool := new(bool)
	err = d.Decode(tmpBool)
	if err != nil {
		return
	}
	p.HideWine = !*tmpBool
	err = d.Decode(tmpBool)
	if err != nil {
		return
	}
	p.ShowPortable = !*tmpBool
	err = d.Decode(tmpBool)
	if err != nil {
		return
	}
	//versionNewest
	err = d.Decode(tmpBool)
	if err != nil {
		return
	}
	p.AppImageDirs = *tmpBool
	err = d.Decode(tmpBool)
	if err != nil {
		return
	}
	p.Beta = !*tmpBool
	return
}
