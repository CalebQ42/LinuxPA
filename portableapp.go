package main

import (
	"debug/elf"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/probonopd/go-appimage/src/goappimage"
)

type portableApp struct {
	root          string
	name          string
	icon          []byte
	execs         []string
	appimageExecs []bool
	wineExecs     []string
	appimage      bool
}

func processPortableApp(root fs.DirEntry) (*portableApp, error) {
	var pa portableApp
	pa.root = root.Name()
	if !root.IsDir() {
		pa.appimage = true
		appimage, err := goappimage.NewAppImage("PortableApps/" + pa.root)
		if err != nil {
			return nil, nil
		}
		err = setExecPermissions("PortableApps/" + pa.root)
		if err != nil {
			if verbose {
				log.Println("Can't set exec permissions for", root.Name(), ":", err)
				log.Println("Ignoring")
			}
			return nil, err
		}
		pa.name = appimage.Name
		var iconRdr io.ReadCloser
		iconRdr, _, err = appimage.Icon()
		if err == nil {
			pa.icon, _ = io.ReadAll(iconRdr)
			iconRdr.Close()
		} else {
			iconRdr, err = appimage.Thumbnail()
			if err == nil {
				pa.icon, _ = io.ReadAll(iconRdr)
				iconRdr.Close()
			}
		}
		return &pa, nil
	}
	pa.name = root.Name()
	dirs, err := os.ReadDir("PortableApps/" + root.Name())
	if err != nil {
		return nil, err
	}
	for i := range dirs {
		if strings.HasSuffix(dirs[i].Name(), ".exe") {
			pa.wineExecs = append(pa.wineExecs, dirs[i].Name())
		}
		filename := "PortableApps/" + root.Name() + "/" + dirs[i].Name()
		if !dirs[i].Type().IsRegular() {
			continue
		}
		var fil *os.File
		fil, err = os.Open(filename)
		if err != nil {
			continue
		}
		defer fil.Close()
		//check if file is a script
		maybeShebang := make([]byte, 2)
		_, err = fil.Read(maybeShebang)
		if err != nil {
			//If it can't eaven be read this much, then it's probably not execuable anyway
			continue
		}
		if string(maybeShebang) == "#!" {
			err = setExecPermissions(filename)
			if err != nil {
				if verbose {
					log.Println("Can't set exec permissions for", dirs[i].Name(), ":", err)
					log.Println("Ignoring")
				}
				continue
			}
			pa.execs = append(pa.execs, dirs[i].Name())
			pa.appimageExecs = append(pa.appimageExecs, false)
			continue
		}
		//check if file has an elf header
		var e *elf.File
		e, err = elf.NewFile(fil)
		if err != nil {
			continue
		}
		//and is an executable
		if e.Type != elf.ET_EXEC {
			continue
		}
		//and is of the proper platform
		//TODO: Replace this with values based on the current platform.
		//This would allow LinuxPA to work (in theory) on all Go and Fyne supported platforms (minus Android)
		if e.Machine != elf.EM_386 && e.Machine != elf.EM_X86_64 {
			continue
		}
		err = setExecPermissions(filename)
		if err != nil {
			if verbose {
				log.Println("Can't set exec permissions for", dirs[i].Name(), ":", err)
				log.Println("Ignoring")
			}
			continue
		}
		pa.execs = append(pa.execs, dirs[i].Name())
		//and if it's an appimage
		_, err = goappimage.NewAppImage("PortableApps/" + root.Name() + "/" + dirs[i].Name())
		pa.appimageExecs = append(pa.appimageExecs, err == nil)
	}
	return &pa, nil
}

func (pa *portableApp) buildUI() fyne.CanvasObject {
	item := make([]*widget.AccordionItem, 0)
	for _, e := range pa.execs {
		item = append(item, widget.NewAccordionItem(e, nil))
	}
	return widget.NewAccordion(item...)
}

func setExecPermissions(file string) error {
	stat, err := os.Stat(file)
	if err != nil {
		return err
	}
	perm := stat.Mode() & os.ModePerm
	if perm&0400 == 0400 {
		//we have owner exec permissions. Hopefully that's enough
		return nil
	}
	if perm&0040 == 0040 {
		//we have group exec permissions. Hopefully that's enough
		return nil
	}
	//Otherwise we try to set both owner and group exec permissions, just to be sure.
	//If this doesn't work, then the runner of LinuxPA isn't an owner or group owner of the exec file.
	//There isn't much we can do about that. We could possibly check the file's owners, but
	//that would be more complexity and is honestly a problem for the LinuxPA runner.
	return os.Chmod(file, perm|0440)
}
