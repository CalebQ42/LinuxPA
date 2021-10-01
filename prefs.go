package main

const (
	prefs    = "PortableApps/LinuxPACom/preferences.json"
	oldPrefs = "PortableApps/LinuxPACom/Prefs.gob"
)

var (
	verbose      bool
	allowWine    bool
	hidePortable bool
	fromRoot     bool
	commonSh     = "PortableApps/LinuxPACom/common.sh" //If empty, don't use common.sh
)
