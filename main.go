package main

import (
	"embed"
	"log"
	"os"
	"os/exec"

	"github.com/CalebQ42/LinuxPA/internal/apps"
	"github.com/CalebQ42/LinuxPA/internal/prefs"
)

const (
	commonSh  = "PortableApps/LinuxPACom/common.sh"
	version   = "3.0.0.0"
	wineDlURL = "https://www.playonlinux.com/wine/binaries/phoenicis/staging-linux-ARCH/PlayOnLinux-wine-VERSION-staging-linux-ARCH.tar.gz"
)

var (
	p *prefs.Prefs
	a []*apps.App
	//go:embed embed
	embedFS embed.FS

	winePath = "PortableApps/LinuxPACom/Wine/bin" //if == "", no wine found
)

func main() {
	initialize()
	//TODO: show UI
	go checkUpdate()
}

func initialize() {
	//detect v2 via old preference format
	var err error
	p, err = prefs.LoadOldPrefs()
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
	a, err = apps.ProcessAllApps(p)
	if err != nil {
		if err != nil {
			log.Fatal("Can't read/process apps:", err)
		}
	}
	_, err = os.Open(winePath)
	if err != nil {
		winePath, _ = exec.LookPath("wine")
	}
}
