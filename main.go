package main

import (
	"bufio"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/nelsam/gxui/drivers/gl"
)

var (
	appMaster map[string][]prtap
	cats      []string
)

type prtap struct {
	name string
	cat  string
	ex   string
	desc string
}

func main() {
	appMaster = make(map[string][]prtap)
	os.Mkdir("PortableApps", 0777)
	pa, err := os.Open("PortableApps")
	if err != nil {
		panic(err)
	}
	appstmp, _ := pa.Readdir(-1)
	var folds []string
	for _, v := range appstmp {
		if v.IsDir() {
			folds = append(folds, v.Name())
		}
	}
	sort.Strings(folds)
	for _, v := range folds {
		fi, _ := os.Open("PortableApps/" + v)
		pat := processApp(fi)
		if (pat != prtap{}) {
			if _, ok := appMaster[pat.cat]; !ok {
				cats = append(cats, pat.cat)
			}
			appMaster[pat.cat] = append(appMaster[pat.cat], pat)
		}
	}
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
		out.cat = "other"
	}
	if out.name == "" {
		out.name = path.Base(fi.Name())
	}
	if out.cat == "" {
		out.cat = "other"
	}
	for _, v := range fis {
		if !v.IsDir() && strings.HasSuffix(v.Name(), ".sh") {
			//do os check here for possible cross platform support
			out.ex = fi.Name() + "/" + v.Name()
			return
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
