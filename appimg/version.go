package appimg

import (
	"strconv"
	"strings"
)

func compareVersions(imgs []appimg) int {
	for i := range imgs {
		imgs[i].version = removeLetters(imgs[i].version)
	}
	highest := 0
	higharr := strings.Split(imgs[0].version, ".")
	for i := 0; i < len(imgs); i++ {
		if i != highest {
			varr := strings.Split(imgs[i].version, ".")
			if len(higharr) < len(varr) {
				for j := 0; j < len(higharr); j++ {
					h, _ := strconv.Atoi(higharr[j])
					c, _ := strconv.Atoi(varr[j])
					if h > c {
						break
					} else if c > h {
						highest = i
						higharr = varr
						break
					}
				}
			}
		}
	}
	return highest
}
