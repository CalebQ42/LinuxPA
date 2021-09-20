package main

import (
	"log"
	"os"
)

var (
	apps []*portableApp
)

func processApps() error {
	pa, err := os.Open("PortableApps")
	if os.IsNotExist(err) {
		os.Mkdir("PortableApps", 0777)
		pa, err = os.Open("PortableApps")
	}
	if err != nil {
		return err
	}
	dirs, err := pa.ReadDir(0)
	if err != nil {
		return err
	}
	var portable *portableApp
	for i := range dirs {
		portable, err = processPortableApp(dirs[i])
		if err != nil {
			if verbose {
				log.Println("Error while processing", dirs[i].Name(), ":", err)
				log.Println("Ignoring")
			}
			continue
		}
		if portable != nil {
			apps = append(apps, portable)
		}
	}
	return nil
}
