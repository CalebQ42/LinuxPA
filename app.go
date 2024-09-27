package main

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/probonopd/go-appimage/src/goappimage"
)

type ExeType byte

const (
	Script = ExeType(iota + 1)
	AppImage
	ELF
	Win
)

type Exe struct {
	Filename string
	Type     ExeType
}

type App struct {
	Name  string
	Execs []Exe
	// Image []byte TODO
}

func ProcessDir(dir string) (*App, error) {
	fils, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("Unable to list files in %v: %v\n", dir, err)
		return nil, err
	}
	app := App{
		Name: filepath.Base(dir),
	}
	for _, f := range fils {
		if f.IsDir() {
			continue
		}
		var t ExeType
		if strings.HasSuffix(strings.ToLower(f.Name()), ".exe") {
			t = Win
		} else {
			fil, err := os.Open(filepath.Join(dir, f.Name()))
			if err != nil {
				log.Printf("Error opening %v: %v\n", filepath.Join(dir, f.Name()), err)
				continue
			}
			startByts := make([]byte, 4)
			_, err = fil.Read(startByts)
			if err != nil {
				log.Printf("Error reading starting bytes of %v: %v\n", filepath.Join(dir, f.Name()), err)
				continue
			}
			if bytes.HasPrefix(startByts, []byte("#!")) {
				t = Script
			} else if bytes.Contains(bytes.ToLower(startByts), []byte("elf")) {
				if strings.HasSuffix(strings.ToLower(f.Name()), ".so") || strings.Contains(strings.ToLower(f.Name()), ".so.") {
					continue
				}
				if goappimage.IsAppImage(filepath.Join(dir, f.Name())) {
					t = AppImage
				} else {
					t = ELF
				}
			}
		}
		if t != 0 {
			app.Execs = append(app.Execs, Exe{
				Filename: f.Name(),
				Type:     t,
			})
		}
	}
	if len(app.Execs) == 0 {
		return nil, nil
	}
	slices.SortFunc(app.Execs, func(a, b Exe) int {
		if a.Type > b.Type {
			return 1
		} else if a.Type < b.Type {
			return -1
		}
		return strings.Compare(a.Filename, b.Filename)
	})
	//TODO: get "proper" name & icon, either from AppImage or PortableApps spec
	return &app, nil
}
