package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/CalebQ42/LinuxPA/internal/apps"
	"github.com/CalebQ42/LinuxPA/internal/prefs"
	"github.com/CalebQ42/squashfs"
)

func TestProcessing(t *testing.T) {
	err := os.Chdir("LinuxPATest")
	if os.IsNotExist(err) {
		var resp *http.Response
		resp, err = http.DefaultClient.Get("https://darkstorm.tech/LinuxPATest.sfs")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		var rdr *squashfs.Reader
		rdr, err = squashfs.NewSquashfsReaderFromReader(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		err = rdr.ExtractTo("LinuxPATest")
		if err != nil {
			t.Fatal(err)
		}
		err = os.Chdir("LinuxPATest")
		if err != nil {
			log.Fatal(err)
		}
	}
	ap, err := apps.ProcessAllApps(&prefs.Prefs{})
	for _, a := range ap {
		fmt.Println(a.String())
	}
	t.Fatal(err)
}
