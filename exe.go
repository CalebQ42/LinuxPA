package main

import (
	"debug/elf"
	"errors"
	"os"
	"path"
	"strings"

	"github.com/probonopd/go-appimage/src/goappimage"
)

type exe struct {
	path string
	name string
	wine bool
	ai   bool
}

func ParseExe(file string) (out exe, err error) {
	if !strings.HasPrefix(file, "/") {
		var wd string
		wd, err = os.Getwd()
		if err != nil {
			return
		}
		file = path.Join(wd, file)
	}
	file = path.Clean(file)
	out.path = file
	out.wine = strings.HasSuffix(file, ".exe")
	out.name = path.Base(file)
	extInd := strings.LastIndex(out.name, ".")
	if extInd != -1 {
		out.name = out.name[:extInd]
	}
	if out.wine {
		//If it's a wine executable, we don't need to do more processing
		return
	}
	exeFil, err := os.Open(file)
	if err != nil {
		return
	}
	//Check for a shebang (#!). If present, file is a script.
	shebangCheck := make([]byte, 2)
	_, err = exeFil.Read(shebangCheck)
	if err != nil {
		return
	} else if string(shebangCheck) == "#!" {
		return
	}
	elfFil, err := elf.Open(file)
	if err != nil {
		return
	}
	//and is an executable
	if elfFil.Type != elf.ET_EXEC {
		return out, errors.New("file not executable elf")
	}
	//and is of the proper platform
	//TODO: Replace this with values based on the current platform.
	//This would allow LinuxPA to work (in theory) on all Linux systems (minus Android)
	if elfFil.Machine != elf.EM_386 && elfFil.Machine != elf.EM_X86_64 {
		return out, errors.New("file not executable for the current platform")
	}
	//and if it's an appimage
	ai, err := goappimage.NewAppImage(file)
	out.ai = (err == nil)
	if out.ai {
		out.name = ai.Name
	}
	err = nil
	return
}

func (e exe) isOwnerExecutable() bool {
	stat, err := os.Stat(e.path)
	if err != nil {
		return false
	}
	perm := stat.Mode() & os.ModePerm
	return perm&0100 == 0100
}

func (e exe) setOwnerExecutable() error {
	stat, err := os.Stat(e.path)
	if err != nil {
		return err
	}
	perm := stat.Mode() & os.ModePerm
	if perm&0100 == 0100 {
		return nil
	}
	return os.Chmod(e.path, perm|0100)
}

func (e exe) isGroupExecutable() bool {
	stat, err := os.Stat(e.path)
	if err != nil {
		return false
	}
	perm := stat.Mode() & os.ModePerm
	return perm&0010 == 0010
}

func (e exe) setGroupExecutable() error {
	stat, err := os.Stat(e.path)
	if err != nil {
		return err
	}
	perm := stat.Mode() & os.ModePerm
	if perm&0010 == 0010 {
		return nil
	}
	return os.Chmod(e.path, perm|0010)
}

func (e exe) isOtherExecutable() bool {
	stat, err := os.Stat(e.path)
	if err != nil {
		return false
	}
	perm := stat.Mode() & os.ModePerm
	return perm&0001 == 0001
}

func (e exe) setOtherExecutable() error {
	stat, err := os.Stat(e.path)
	if err != nil {
		return err
	}
	perm := stat.Mode() & os.ModePerm
	if perm&0001 == 0001 {
		return nil
	}
	return os.Chmod(e.path, perm|0001)
}

func (e exe) displayName() (out string) {
	out = e.name
	if e.wine {
		out += " (Wine)"
	}
	return
}

func (e exe) String() string {
	return e.displayName() + ": " + e.path
}
