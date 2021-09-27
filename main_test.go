package main

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/CalebQ42/squashfs"
)

func TestUI(t *testing.T) {
	_, err := os.Open("LinuxPATest")
	if os.IsNotExist(err) {
		err = setupTestDir()
		if err != nil {
			t.Fatal(err)
		}
		_, err = os.Open("LinuxPATest")
	}
	if err != nil {
		t.Fatal(err)
	}
	err = os.Chdir("LinuxPATest")
	if err != nil {
		t.Fatal(err)
	}
	main()
}

func setupTestDir() error {
	resp, err := http.DefaultClient.Get("https://darkstorm.tech/LinuxPATest.sfs")
	if err != nil {
		return err
	}
	os.Remove("LinuxPATest.sfs")
	sfsFil, err := os.Create("LinuxPATest.sfs")
	if err != nil {
		return err
	}
	_, err = io.Copy(sfsFil, resp.Body)
	if err != nil {
		return err
	}
	rdr, err := squashfs.NewSquashfsReader(sfsFil)
	if err != nil {
		return err
	}
	err = rdr.ExtractTo("LinuxPATest")
	if err != nil {
		return err
	}
	return nil
}
