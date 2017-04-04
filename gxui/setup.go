package main

import (
	"bufio"
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/nelsam/gxui"
)

func setup() {
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
	fmt.Println(err)
	_, err = os.Open("PortableApps/LinuxPACom/common.sh")
	if err == nil {
		comEnbld = true
	}
	fi, err := os.Open("PortableApps/LinuxPACom/Info.ini")
	if err != nil {
		fi, err = os.Create("PortableApps/LinuxPACom/Info.ini")
		if err == nil {
			wrt := bufio.NewWriter(fi)
			wrt.WriteString(defIni)
			wrt.Flush()
		}
	}
	if err == nil {
		rdr := bufio.NewReader(fi)
		for err != nil {
			ln, _, error := rdr.ReadLine()
			err = error
			str := string(ln)
			if strings.HasPrefix(str, "theme=") {
				str = strings.TrimPrefix(str, "theme=")
				if str == "lt" {
					darkTheme = false
				}
			}
		}
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
}

func processApp(fold string) (out app) {
	wd, _ := os.Getwd()
	out.dir = wd + "/" + fold
	out.ini = findInfo(fold)
	if out.ini != nil {
		out.name = getName(out.ini)
		out.ini = findInfo(fold)
		out.cat = getCat(out.ini)
		out.ini = findInfo(fold)
	}
	if out.name == "" {
		out.name = strings.TrimPrefix(fold, "PortableApps/")
	}
	if out.cat == "" {
		out.cat = "Other"
	}
	out.icon = getIcon(fold)
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
	}
	return
}

func getCat(ini *os.File) string {
	rdr := bufio.NewReader(ini)
	var ret string
	for line, _, err := rdr.ReadLine(); err == nil; line, _, err = rdr.ReadLine() {
		if strings.HasPrefix(string(line), "Category=") {
			ret = strings.TrimPrefix(string(line), "Category=")
			break
		}
	}
	rdr.Reset(ini)
	return ret
}

func getName(ini *os.File) string {
	rdr := bufio.NewReader(ini)
	var ret string
	for line, _, err := rdr.ReadLine(); err == nil; line, _, err = rdr.ReadLine() {
		if strings.HasPrefix(string(line), "Name=") {
			ret = strings.TrimPrefix(string(line), "Name=")
			break
		}
	}
	rdr.Reset(ini)
	return ret
}

func getIcon(fold string) gxui.Texture {
	var pic *os.File
	if folder, err := os.Open(fold + "/App/AppInfo"); err == nil {
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
			pic, _ = os.Open(fold + "/App/AppInfo/" + pics[ind])
		}
	} else if fi, err := os.Open(fold + "/appicon.png"); err == nil {
		pic = fi
	} else {
		return nil
	}
	img, _, err := image.Decode(pic)
	if err != nil {
		return nil
	}
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, img.Bounds(), img, image.ZP, draw.Src)
	ret := dr.CreateTexture(rgba, 1)
	return ret
}

func findInfo(fold string) *os.File {
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
	if fi, err := os.Open(fold + "/appinfo.ini"); err == nil {
		return fi
	}
	return nil
}
