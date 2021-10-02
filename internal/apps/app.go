package apps

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/CalebQ42/LinuxPA/internal/prefs"
	"github.com/probonopd/go-appimage/src/goappimage"
	"gopkg.in/ini.v1"
)

type App struct {
	Desktop    *ini.File
	Name       string
	Path       string
	Categories []string
	Icon       []byte
	Execs      []Exe
	AppImage   bool
}

func ProcessApp(p *prefs.Prefs, dir string) ([]*App, error) {
	var a App
	dirFil, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	a.Path = dir
	stat, _ := dirFil.Stat()
	if !stat.IsDir() {
		var e Exe
		e, err = ProcessExe(dir)
		if err != nil {
			return nil, err
		}
		if !e.appimage {
			return nil, errors.New("file not in application folder and not AppImage")
		}
		a.AppImage = true
		a.Execs = append(a.Execs, e)
		ai, _ := goappimage.NewAppImage(dir)
		a.Name = ai.Name
		a.Categories = strings.Split(ai.Desktop.Section("Desktop Entry").Key("Categories").String(), "；")
		var iconRdr io.ReadCloser
		iconRdr, _, err = ai.Icon()
		if err != nil {
			iconRdr, err = ai.Thumbnail()
			if err != nil {
				if p.Verbose {
					log.Println("Can't get icon for", a.Name, ":", err)
				}
				return []*App{&a}, nil
			}
		}
		defer iconRdr.Close()
		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, iconRdr)
		if err != nil {
			if p.Verbose {
				log.Println("Can't get icon for", a.Name, ":", err)
			}
			return []*App{&a}, nil
		}
		a.Icon = buf.Bytes()
		return []*App{&a}, nil
	}
	dirs, err := dirFil.ReadDir(0)
	if err != nil {
		return nil, err
	}
	var desktopFiles []string
	var e Exe
	for i := range dirs {
		if dirs[i].IsDir() {
			continue
		}
		if strings.HasSuffix(dirs[i].Name(), ".desktop") {
			desktopFiles = append(desktopFiles, dirs[i].Name())
			continue
		}
		e, err = ProcessExe(path.Join(dir, dirs[i].Name()))
		if err != nil {
			if p.Verbose {
				log.Println(path.Join(dir, dirs[i].Name()), "is not an executable. Ignoring...")
			}
			continue
		}
		a.Execs = append(a.Execs, e)
	}
	if len(a.Execs) == 0 {
		return nil, errors.New("application folder contains no executables")
	}
	a.orderExecs()
	if len(desktopFiles) == 0 {
		return []*App{&a}, nil
	}
	var out []*App
	for _, d := range desktopFiles {
		var deskFil *os.File
		deskFil, err = os.Open(path.Join(dir, d))
		if err != nil {
			if p.Verbose {
				log.Println("Error while opening desktop file", path.Join(dir, d), ":", err)
				log.Println("Ignoring...")
			}
			continue
		}
		var desktopApp App
		desktopApp.Desktop, err = ini.Load(deskFil)
		if err != nil {
			if p.Verbose {
				log.Println("Error while processing desktop file", path.Join(dir, d), ":", err)
				log.Println("Ignoring...")
			}
			continue
		}
		exec := desktopApp.Desktop.Section("DesktopEntry").Key("Exec").String()
		if exec == "" {
			if p.Verbose {
				log.Println("Desktop file", path.Join(dir, d), "does not have an Exec key. Ignoring...")
			}
			continue
		}
		var args string
		execSplit := strings.Split(exec, "")
		exec = execSplit[0]
		for i := 1; i < len(execSplit); i++ {
			if strings.HasPrefix(execSplit[i], "$") {
				continue
			}
			args += " " + execSplit[i]
		}
		args = strings.TrimSpace(args)
		//TODO
		desktopApp.Name = desktopApp.Desktop.Section("Desktop Entry").Key("Name").String()
		desktopApp.Categories = strings.Split(desktopApp.Desktop.Section("Desktop Entry").Key("Categories").String(), "；")

	}
	if len(out) == 0 {
		return []*App{&a}, nil
	}
	//TODO: Process information from PortableApps.com format and, if not there, any AppImages present.
	//TODO: Process any .desktop files as a new app (if there are multiple)
	return out, nil
}

func (a *App) orderExecs() {
	var scripts, appimages, linExecs, wine []Exe
	for _, e := range a.Execs {
		if e.script {
			scripts = append(scripts, e)
		} else if e.appimage {
			appimages = append(appimages, e)
		} else if e.IsWine() {
			wine = append(wine, e)
		} else {
			linExecs = append(linExecs, e)
		}
	}
	a.Execs = append(append(append(scripts, appimages...), linExecs...), wine...)
}

func ProcessAllApps(p *prefs.Prefs) (a []*App, err error) {
	dirs, err := os.ReadDir("PortableApps")
	if err != nil {
		return
	}
	for i := range dirs {
		var tmp []*App
		tmp, err = ProcessApp(p, "PortableApps/"+dirs[i].Name())
		if err != nil {
			if p.Verbose {
				log.Println("Error while processing", dirs[i].Name(), ":", err)
				log.Println("Ignoring...")
			}
			continue
		}
		a = append(a, tmp...)
	}
	sort.Slice(a, func(i, j int) bool {
		return a[i].Name < a[j].Name
	})
	return
}
