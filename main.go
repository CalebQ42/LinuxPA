package main

import (
	"embed"
	"log"
	"os"

	"github.com/CalebQ42/LinuxPA/internal/apps"
	"github.com/CalebQ42/LinuxPA/internal/prefs"
)

const (
	commonSh = "PortableApps/LinuxPACom/common.sh"
	version  = "3.0.0.0"
)

var (
	p *prefs.Prefs
	a []*apps.App
	//go:embed embed
	embedFS embed.FS
)

func main() {
	//detect v2 via old preference format
	p, err := prefs.LoadOldPrefs()
	if err == nil {
		err = upgrade()
		if err != nil {
			log.Fatal("Can't upgrade from v2 to v3. Aboring:", err)
		}
	} else {
		p, err = prefs.LoadPrefs()
		if err != nil && !os.IsNotExist(err) {
			log.Fatal("Can't load preferences:", err)
		}
	}
	a, err := apps.ProcessAllApps(p)
	if err != nil {
		if err != nil {
			log.Fatal("Can't read/process apps:", err)
		}
	}
	_ = a
	//TODO: check for wine. If not installed, disable showWine
	//TODO: load apps
	//TODO: show UI
	//TODO: Check for updates
}
