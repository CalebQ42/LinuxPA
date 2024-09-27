package main

import (
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"testing"
)

const (
	testImage = "https://darkstorm.tech/files/LinuxPATest.sfs"
)

func setupTestEnv() error {
	_, err := exec.LookPath("unsquashfs")
	if err != nil {
		return errors.New("unsquashfs not installed")
	}
	fold, err := os.Stat("testing")
	if os.IsNotExist(err) {
		err = os.Mkdir("testing", 0777)
		if err != nil {
			return err
		}
		fold, err = os.Stat("testing")
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	os.RemoveAll("testing/PortableApps")
	os.RemoveAll("testing/Documents")
	os.RemoveAll("testing/Start.exe")
	if !fold.IsDir() {
		return errors.New("./testing is not a directory!!!")
	}
	img, err := os.Open("testing/LinuxPATest.sfs")
	if os.IsNotExist(err) {
		img, err = os.Create("testing/LinuxPATest.sfs")
		if err != nil {
			return errors.New("Cannot create testing/LinuxPATest.sfs")
		}
		resp, err := http.DefaultClient.Get(testImage)
		if err != nil {
			return err
		}
		_, err = io.Copy(img, resp.Body)
		resp.Body.Close()
		if err != nil {
			return err
		}
	}
	err = exec.Command("unsquashfs", "-d", "./testing", "./testing/LinuxPATest.sfs").Run()
	if err != nil {
		return err
	}
	return nil
}

func TestStuff(t *testing.T) {
	setupTestEnv()
	main()
}
