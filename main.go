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
	err := loadPrefs()
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Can't load preferences:", err)
	}
}
