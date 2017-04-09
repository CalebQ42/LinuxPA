//Package appimg is to download and update AppImages for LinuxPA.
//Converted from github.com/CalebQ42/bbConvert
package appimg

import "reflect"

//Convert converts the input string. Only returns <a> tags
func Convert(in string) (out []Tag) {
	for i := 0; i < len(in); i++ {
		v := in[i]
		if v == '<' {
			for j := i; j < len(in); j++ {
				val := in[j]
				if val == '>' {
					var tmp Tag
					tmp.process(in[i+1 : j+1])
					if !tmp.end {
						tmp.index[0] = i
						tmp.index[1] = j
						nd := fndend(tmp, in[j+1:])
						if !reflect.DeepEqual(nd, Tag{}) && tmp.typ == "a" {
							out = append(out, tmp)
						}
					}
					break
				}
			}
		}
	}
	return
}

func fndend(fnt Tag, area string) Tag {
	var count int
	for i, v := range area {
		if v == '<' {
			for j, val := range area[i:] {
				if val == '>' {
					var tmp Tag
					tmp.process(area[i+1 : i+j+1])
					if tmp.typ == fnt.typ {
						if tmp.end {
							if count == 0 {
								tmp.index[0] = fnt.index[1] + 1 + i
								tmp.index[1] = fnt.index[1] + j + i + 1
								return tmp
							}
							count--
							break
						} else {
							count++
							break
						}
					}
				}
			}
		}
	}
	return Tag{}
}
