package main

import (
	"bufio"
	"fmt"
	"os"
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
	var hasEx bool
	out.cat = "other"
	fis, _ := fi.Readdir(-1)
	for _, v := range fis {
		if v.IsDir() && v.Name() == "App" {
			fild, err := os.Open(fi.Name() + "/App/AppInfo/appinfo.ini")
			fmt.Println(fi.Name() + "/App/AppInfo/appinfo.ini")
			if err == nil {
				fmt.Println("working!")
				out.name = getName(*fild)
				fild, _ = os.Open(fi.Name() + "/App/AppInfo/appinfo.ini")
				out.cat = getCat(*fild)
				fmt.Println("Name:", out.name)
			}
		} else if !v.IsDir() {
			//do os check here
			if strings.HasSuffix(v.Name(), ".sh") {
				hasEx = true
				out.ex = fi.Name() + "/" + v.Name()
				if out.name == "" {
					out.name = strings.TrimSuffix(v.Name(), ".sh")
				}
			}
		}
	}
	if hasEx {
		return
	}
	return prtap{}
}

func getCat(fi os.File) (out string) {
	rdr := bufio.NewReader(&fi)
	var err error
	var ln []byte
	for err == nil {
		ln, _, err = rdr.ReadLine()
		str := string(ln)
		if strings.HasPrefix(str, "Category=") {
			fmt.Println(str)
			out = strings.TrimPrefix(str, "Category=")
			return
		}
	}
	return
}

func getName(fi os.File) (out string) {
	rdr := bufio.NewReader(&fi)
	var err error
	var ln []byte
	for err == nil {
		ln, _, err = rdr.ReadLine()

		str := string(ln)
		fmt.Println(str)
		if strings.HasPrefix(str, "Name=") {
			out = strings.TrimPrefix(str, "Name=")
			return
		}
	}
	return
}
