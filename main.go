package main

import (
	"embed"
	"log"
	"os"
)

var (
	version = "3.0.0.0"
	apps    []*app
	//go:embed embed
	embedFS embed.FS
)

func main() {
	//detect v2 via old preference format
	_, err := os.Open(oldPrefs)
	if err == nil {
		err = upgrade()
		if err != nil {
			log.Fatal("Can't upgrade from v2 to v3. Aborting:", err)
		}
	}
	err = loadPrefs()
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Can't load preferences:", err)
	}
}
