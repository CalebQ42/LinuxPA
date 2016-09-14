package main

import (
	"bufio"
	"os"
	"path"
	"reflect"
	"sort"
	"strings"

	"github.com/nelsam/gxui/drivers/gl"
)

var (
	appMaster map[string][]prtap
	cats      []string
	wineOnly  []string
	linOnly   []string
	conf      *os.File
	common    string
	commEnbl  bool
)

type prtap struct {
	name string
	cat  string
	ex   string
	desc string
	wine bool
}

func main() {
	commEnbl = true
	appMaster = make(map[string][]prtap)
	os.Mkdir("PortableApps", 0777)
	os.Mkdir("PortableApps/LinuxPACom", 0777)
	common = "PortableApps/LinuxPACom/common.sh"
	_, err := os.Open(common)
	if os.IsNotExist(err) {
		commEnbl = false
	}
	pa, err := os.Open("PortableApps")
	if err != nil {
		panic(err)
	}
	appstmp, _ := pa.Readdir(-1)
	var folds []string
	for _, v := range appstmp {
		if v.IsDir() && v.Name() != "LinuxPACom" && v.Name() != "PortableApps.com" {
			folds = append(folds, v.Name())
		}
	}
	sort.Strings(folds)
	for _, v := range folds {
		fi, _ := os.Open("PortableApps/" + v)
		pat := processApp(fi)
		if (pat != prtap{}) {
			if _, ok := appMaster[pat.cat]; !ok {
				if pat.wine {
					wineOnly = append(wineOnly, pat.cat)
					cats = append(cats, pat.cat)
				} else {
					linOnly = append(linOnly, pat.cat)
					cats = append(cats, pat.cat)
				}
			} else {
				if !pat.wine {
					for i, v := range wineOnly {
						if pat.cat == v {
							wineOnly = append(wineOnly[:i], wineOnly[i+1:]...)
							linOnly = append(linOnly, pat.cat)
							break
						}
					}
				}
			}
			appMaster[pat.cat] = append(appMaster[pat.cat], pat)
		}
	}
	sort.Strings(linOnly)
	sort.Strings(wineOnly)
	sort.Strings(cats)
	gl.StartDriver(uiMain)
}

func processApp(fi *os.File) (out prtap) {
	fis, _ := fi.Readdir(-1)
	if fil, err := os.Open(fi.Name() + "/App/AppInfo/appinfo.ini"); err == nil {
		out.name = getName(fil)
		fil, _ = os.Open(fi.Name() + "/App/AppInfo/appinfo.ini")
		out.cat = getCat(fil)
	} else if fil, err := os.Open(fi.Name() + "/appinfo.ini"); err == nil {
		out.name = getName(fil)
		fil, _ = os.Open(fi.Name() + "/appinfo.ini")
		out.cat = getCat(fil)
	} else {
		out.cat = "Other"
	}
	if out.name == "" {
		out.name = path.Base(fi.Name())
	}
	if out.cat == "" {
		out.cat = "Other"
	}
	//executable detection
	wd, _ := os.Getwd()
	var rdr *bufio.Reader
	for _, v := range fis {
		fil, err := os.Open(wd + "/" + fi.Name() + "/" + v.Name())
		if err == nil {
			stat, _ := fil.Stat()
			if !stat.IsDir() {
				rdr = bufio.NewReader(fil)
				shebang := []byte{'#', '!'}
				two := make([]byte, 2)
				rdr.Read(two)
				if reflect.DeepEqual(shebang, two) {
					out.ex = wd + "/" + fi.Name() + "/" + v.Name()
					rdr.Reset(fil)
					return
				}
			}
		}
	}
	for _, v := range fis {
		fil, err := os.Open(wd + "/" + fi.Name() + "/" + v.Name())
		if err == nil {
			stat, _ := fil.Stat()
			if !stat.IsDir() {
				rdr = bufio.NewReader(fil)
				thr := make([]byte, 4)
				rdr.Read(thr)
				if strings.Contains(string(thr), "ELF") {
					out.ex = wd + "/" + fi.Name() + "/" + v.Name()
					rdr.Reset(fil)
					return
				}
			}
		}
	}
	for _, v := range fis {
		fil, err := os.Open(wd + "/" + fi.Name() + "/" + v.Name())
		if err == nil {
			stat, _ := fil.Stat()
			if !stat.IsDir() && strings.HasSuffix(stat.Name(), "exe") {
				out.wine = true
				out.ex = wd + "/" + fi.Name() + "/" + v.Name()
				out.name += " (Wine)"
				return
			}
		}
	}
	return prtap{}
}

func getCat(fi *os.File) (out string) {
	rdr := bufio.NewReader(fi)
	var err error
	var ln []byte
	for err == nil {
		ln, _, err = rdr.ReadLine()
		str := string(ln)
		if strings.HasPrefix(str, "Category=") {
			out = strings.TrimPrefix(str, "Category=")
			return
		}
	}
	return
}

func getName(fi *os.File) (out string) {
	rdr := bufio.NewReader(fi)
	var err error
	var ln []byte
	for err == nil {
		ln, _, err = rdr.ReadLine()
		str := string(ln)
		if strings.HasPrefix(str, "Name=") {
			out = strings.TrimPrefix(str, "Name=")
			return
		}
	}
	return
}
