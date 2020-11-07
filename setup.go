package main

import (
	"bufio"
	_ "image/png"
	"os"
	"os/exec"
	"reflect"
	"sort"
	"strings"

	goappimage "github.com/CalebQ42/GoAppImage"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func setup() {
	loadPrefs()
	if _, err := os.Open("PortableApps/LinuxPACom/Wine"); os.IsNotExist(err) {
		if _, errd := exec.LookPath("wine"); errd == nil {
			wineAvail = true
		}
	} else if err == nil {
		wineAvail = true
	}
	if !wineAvail {
		wine = false
	}
	PortableAppsFold, err := os.Open("PortableApps")
	if PAStat, _ := PortableAppsFold.Stat(); err != nil || !PAStat.IsDir() {
		os.Mkdir("PortableApps", 0777)
		PortableAppsFold, err = os.Open("PortableApps")
		if err != nil {
			panic("Can't find PortableApps folder and can't create one!")
		}
	}
	if _, err = os.Open("PortableApps/LinuxPACom"); err != nil {
		os.Mkdir("PortableApps/LinuxPACom", 0777)
	}
	_, err = os.Open("PortableApps/LinuxPACom/common.sh")
	if err == nil {
		comEnbld = true
	}
	PAFolds, _ := PortableAppsFold.Readdirnames(-1)
	sort.Strings(PAFolds)
	for _, v := range PAFolds {
		fold, _ := os.Open("PortableApps/" + v)
		if stat, _ := fold.Stat(); stat.IsDir() && stat.Name() != "PortableApps.com" && stat.Name() != "LinuxPACom" {
			ap := processApp("PortableApps/" + v)
			if !reflect.DeepEqual(ap, app{}) {
				if _, ok := master[ap.cat]; !ok {
					cats = append(cats, ap.cat)
					sort.Strings(cats)
				}
				if len(ap.lin) != 0 {
					if _, ok := linmaster[ap.cat]; !ok {
						lin = append(lin, ap.cat)
						sort.Strings(lin)
					}
				}
				master[ap.cat] = append(master[ap.cat], ap)
				if len(ap.lin) != 0 {
					linmaster[ap.cat] = append(linmaster[ap.cat], ap)
				}
			}
		}
	}
	populated = true
}

func processApp(fold string) (out app) {
	wd, _ := os.Getwd()
	out.dir = wd + "/" + fold
	folder, _ := os.Open(fold)
	fis, _ := folder.Readdirnames(-1)
	for _, v := range fis {
		tmp, _ := os.Open(fold + "/" + v)
		if stat, _ := tmp.Stat(); stat.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToLower(v), ".appimage") {
			out.appimg = append(out.appimg, v)
			out.ex = append(out.ex, v)
			out.lin = append(out.lin, v)
		} else if strings.HasSuffix(strings.ToLower(v), ".exe") {
			out.ex = append(out.ex, v)
		} else {
			btys := make([]byte, 4)
			rdr := bufio.NewReader(tmp)
			rdr.Read(btys)
			if (strings.Contains(strings.ToLower(string(btys)), "elf") && !strings.HasSuffix(strings.ToLower(v), ".so") && !strings.Contains(v, ".so.")) || strings.HasPrefix(strings.ToLower(string(btys)), "#!") {
				out.ex = append(out.ex, v)
				out.lin = append(out.lin, v)
			}
		}
	}
	if len(out.ex) == 0 {
		return app{}
	}
	if len(out.lin) == 0 {
		out.name += " (Wine)"
		out.wine = true
	}
	out.icon = getIcon(fold)
	out.ini = findInfo(fold)
	if out.ini != nil {
		out.name = getName(out.ini)
		out.cat = getCat(out.ini)
	}
	if len(out.appimg) > 0 && out.name == "" && out.cat == "" && out.icon == nil {
		os.Mkdir(out.dir+"/.appimageconfig", 0777)
		ai := goappimage.NewAppImage(out.dir + "/" + out.appimg[0])
		fil, err := os.Open(out.dir + "/.appimageconfig/the.md5")
		if os.IsNotExist(err) {
			ai.ExtractFile("*.desktop", out.dir+"/.appimageconfig/", false)
			appimageconfig, _ := os.Open(out.dir + "/.appimageconfig")
			appdirs, _ := appimageconfig.Readdirnames(-1)
			for _, dirs := range appdirs {
				desktopFil, _ := os.Open(out.dir + "/.appimageconfig/" + dirs)
				if stat, _ := desktopFil.Stat(); strings.HasSuffix(dirs, ".desktop") && !stat.IsDir() {
					os.Rename(out.dir+"/.appimageconfig/"+dirs, out.dir+"/.appimageconfig/the.desktop")
					break
				}
			}
			desk, _ := os.Open(out.dir + "/.appimageconfig/the.desktop")
			name, cat, icon := extractDesktopInfo(desk)
			if out.name == "" {
				out.name = name
			}
			if out.cat == "" {
				out.cat = cat
			}
			if out.icon == nil {
				it, _ := gtk.IconThemeGetDefault()
				out.icon, err = it.LoadIcon(icon, 64, gtk.ICON_LOOKUP_GENERIC_FALLBACK)
			}
			fil, _ = os.Create(out.dir + "/.appimageconfig/the.md5")
			wrtr := bufio.NewWriter(fil)
			wrtr.WriteString(ai.Md5)
			wrtr.Flush()
		} else {
			rdr := bufio.NewReader(fil)
			filMd, _, _ := rdr.ReadLine()
			oldMd := string(filMd)
			if oldMd != ai.Md5 {
				ai.ExtractFile("*.desktop", out.dir+"/.appimageconfig/", false)
				appimageconfig, _ := os.Open(out.dir + "/.appimageconfig")
				appdirs, _ := appimageconfig.Readdirnames(-1)
				for _, dirs := range appdirs {
					desktopFil, _ := os.Open(out.dir + "/.appimageconfig/" + dirs)
					if stat, _ := desktopFil.Stat(); strings.HasSuffix(dirs, ".desktop") && !stat.IsDir() {
						os.Rename(out.dir+"/.appimageconfig/"+dirs, out.dir+"/.appimageconfig/the.desktop")
						break
					}
				}
				os.Remove(out.dir + "/.appimageconfig/the.md5")
				fil, _ = os.Create(out.dir + "/.appimageconfig/the.md5")
				wrtr := bufio.NewWriter(fil)
				wrtr.WriteString(ai.Md5)
				wrtr.Flush()
			}
		}
		desk, _ := os.Open(out.dir + "/.appimageconfig/the.desktop")
		name, cat, icon := extractDesktopInfo(desk)
		if out.name == "" {
			out.name = name
		}
		if out.cat == "" {
			out.cat = cat
		}
		if out.icon == nil {
			it, _ := gtk.IconThemeGetDefault()
			out.icon, err = it.LoadIcon(icon, 32, gtk.ICON_LOOKUP_GENERIC_FALLBACK)
		}
	}
	if out.name == "" {
		out.name = strings.TrimPrefix(fold, "PortableApps/")
	}
	if out.cat == "" {
		out.cat = "Other"
	}
	if portableHide {
		out.name = strings.TrimSuffix(out.name, "Portable")
	}
	out.name = strings.TrimSpace(out.name)
	return
}

func getCat(ini *os.File) string {
	rdr := bufio.NewReader(ini)
	var ret string
	for line, _, err := rdr.ReadLine(); err == nil; line, _, err = rdr.ReadLine() {
		if strings.HasPrefix(string(line), "Category=") {
			ret = strings.TrimPrefix(string(line), "Category=")
			break
		} else if strings.HasPrefix(string(line), "category=") {
			ret = strings.TrimPrefix(string(line), "category=")
		}
	}
	rdr.Reset(ini)
	return ret
}

func extractDesktopInfo(desk *os.File) (name, category, iconName string) {
	rdr := bufio.NewReader(desk)
	var nameGot, catGot, iconGot bool
	for line, _, err := rdr.ReadLine(); err == nil; line, _, err = rdr.ReadLine() {
		ln := string(line)
		if !nameGot && strings.HasPrefix(ln, "Name=") {
			name = strings.TrimPrefix(ln, "Name=")
			nameGot = true
		} else if !catGot && strings.HasPrefix(ln, "Categories=") {
			cats := strings.Split(strings.TrimPrefix(ln, "Categories="), ";")
			if len(cats) > 0 {
				category = cats[0]
			}
			catGot = true
		} else if !iconGot && strings.HasPrefix(ln, "Icon=") {
			iconName = strings.TrimPrefix(ln, "Icon=")
			iconGot = true
		}
		if nameGot && catGot && iconGot {
			break
		}
	}
	return
}

func getName(ini *os.File) string {
	rdr := bufio.NewReader(ini)
	var ret string
	for line, _, err := rdr.ReadLine(); err == nil; line, _, err = rdr.ReadLine() {
		if strings.HasPrefix(string(line), "Name=") {
			ret = strings.TrimPrefix(string(line), "Name=")
			break
		} else if strings.HasPrefix(string(line), "name=") {
			ret = strings.TrimPrefix(string(line), "name=")
			break
		}
	}
	rdr.Reset(ini)
	return ret
}

func getIcon(fold string) *gdk.Pixbuf {
	var pic string
	if _, err := os.Open(fold + "/appicon.png"); err == nil {
		pic = fold + "/appicon.png"
	} else if folder, err := os.Open(fold + "/App/AppInfo"); err == nil {
		fis, _ := folder.Readdir(-1)
		var pics []string
		for _, v := range fis {
			if !v.IsDir() && strings.HasSuffix(strings.ToLower(v.Name()), ".png") && strings.HasPrefix(strings.ToLower(v.Name()), "appicon_") {
				pics = append(pics, v.Name())
			}
		}
		sort.Strings(pics)
		if len(pics) > 1 {
			var ind int
			if !contains(pics, "appicon_32.png") {
				ind = len(pics) - 1
			} else {
				ind = sort.SearchStrings(pics, "appicon_32.png")
			}
			pic = fold + "/App/AppInfo/" + pics[ind]
		}
	} else {
		img, _ := gtk.ImageNewFromIconName("application-x-executable", gtk.ICON_SIZE_BUTTON)
		buf := img.GetPixbuf()
		return buf
	}
	img, _ := gtk.ImageNewFromFile(pic)
	buf, _ := img.GetPixbuf().ScaleSimple(32, 32, gdk.INTERP_BILINEAR)
	return buf
}

func findInfo(fold string) *os.File {
	if fi, err := os.Open(fold + "/appinfo.ini"); err == nil {
		return fi
	}
	tmp, err := os.Open(fold + "/App/AppInfo")
	if err == nil {
		fis, _ := tmp.Readdirnames(-1)
		for _, v := range fis {
			if strings.ToLower(v) == "appinfo.ini" {
				tmp, _ := os.Open(fold + "/App/AppInfo/" + v)
				return tmp
			}
		}
	}
	return nil
}
