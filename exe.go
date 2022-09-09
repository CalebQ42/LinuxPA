package main

import (
	"bufio"
	"debug/elf"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/probonopd/go-appimage/src/goappimage"
)

var (
	errorNotExecutable = errors.New("file not executable")
)

type executable struct {
	ai     *goappimage.AppImage
	path   string
	script bool
	wine   bool
}

func parseFile(path string) (exe executable, err error) {
	//First, make sure it's actually a file
	file, err := os.Open(path)
	if err != nil {
		return
	}
	stat, _ := file.Stat()
	if stat.IsDir() {
		return exe, errorNotExecutable
	}
	//First check if it's a win executable file (.exe or .bat).
	//TODO: Check if actually executable
	exe.path = path
	if strings.HasSuffix(path, ".exe") || strings.HasSuffix(path, ".bat") {
		exe.wine = true
		return
	}
	//Check for a shebang to see if file is a script
	scriptChecker := bufio.NewReader(file)
	var line string
	for {
		line, err = scriptChecker.ReadString('\n')
		line = strings.TrimSpace(strings.TrimSuffix(line, "\n"))
		if err == io.EOF && len(line) > 0 {
			err = nil
		}
		if err != nil {
			break
		}
		if strings.HasPrefix(line, "#!") {
			exe.script = true
			return
		} else if !strings.HasPrefix(line, "#") {
			break
		}
	}
	//Check for an ELF header (linux executable)
	e, err := elf.Open(path)
	if err == nil {
		//TODO: Check current architecture and check e.Machine accordingly.
		if e.Type == elf.ET_EXEC && (e.Machine == elf.EM_X86_64 || e.Machine == elf.EM_386) {
			//Check if it's an AppImage
			exe.ai, err = goappimage.NewAppImage(path)
			if err != nil {
				exe.ai = nil
				err = nil
			}
			return
		}
	}
	return exe, errorNotExecutable
}
