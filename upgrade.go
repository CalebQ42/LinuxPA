package main

import (
	"io"
	"os"
)

//upgrades from v2 to v3 format
func upgrade() error {
	com, err := os.Open(commonSh)
	if os.IsNotExist(err) {
		err = createCommonSh()
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		_, err = com.WriteString("\n\n#Run the app. DON'T TOUCH (unless you know what your doing :P)\n$@")
		if err != nil {
			return err
		}
		com.Sync()
	}
	err = p.SavePrefs()
	if err != nil {
		return err
	}
	os.Remove("PortableApps/LinuxPACom/Prefs.gob")
	return nil
}

func createCommonSh() error {
	os.MkdirAll("PortableApps.com/LinuxPACom", 0775)
	comEmbed, err := embedFS.Open("common.sh")
	if err != nil {
		return err
	}
	defer comEmbed.Close()
	var comFil *os.File
	comFil, err = os.Create(commonSh)
	if err != nil {
		return err
	}
	defer comFil.Close()
	_, err = io.Copy(comFil, comEmbed)
	if err != nil {
		return err
	}
	comFil.Sync()
	return nil
}
