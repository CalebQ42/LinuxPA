package main

import (
	"bytes"
	"debug/elf"
	"errors"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/probonopd/go-appimage/src/goappimage"
)

type exe struct {
	path     string
	name     string
	appimage bool
}

func ProcessFile(file string) (e exe, err error) {
	if !path.IsAbs(file) {
		var wd string
		wd, err = os.Getwd()
		if err != nil {
			return
		}
		file = path.Join(wd, file)
	}
	file = path.Clean(file)
	e.path = file
	e.name = path.Base(file)
	if strings.Contains(file, ".") {
		e.name = e.name[:strings.LastIndex(file, ".")]
	}
	//Check if file exists/is openable
	fil, err := os.Open(file)
	if err != nil {
		return
	}
	//If it's a windows executable, we don't need to process further
	if e.Wine() {
		return
	}
	//Then Check for a shebang (script file)
	check := make([]byte, 2)
	_, err = fil.Read(check)
	if err != nil {
		return
	}
	if bytes.Equal(check, []byte("#!")) {
		return
	}
	//Then we check if the file is an ELF file that's executable
	elfFile, err := elf.Open(file)
	if err != nil {
		return
	}
	if elfFile.Type != elf.ET_EXEC {
		return e, errors.New("file is not executable elf")
	}
	//Lastly, check if it's an AppImage so we can potentially enable advanced features
	_, err = goappimage.NewAppImage(file)
	e.appimage = err == nil
	err = nil
	return
}

func (e exe) Wine() bool {
	return strings.HasSuffix(e.path, ".exe")
}

func (e exe) Cmd() (cmd *exec.Cmd) {
	if commonSh != "" {
		cmd = exec.Command(commonSh, e.path)
	} else {
		cmd = exec.Command(e.path)
	}
	if !fromRoot {
		cmd.Dir = path.Dir(e.path)
	}
	return
}
