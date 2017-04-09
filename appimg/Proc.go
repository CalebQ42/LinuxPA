package appimg

import "reflect"

func convert(in string) (out []tag) {
	for i := 0; i < len(in); i++ {
		v := in[i]
		if v == '<' {
			for j := i; j < len(in); j++ {
				val := in[j]
				if val == '>' {
					var tmp tag
					tmp.process(in[i+1 : j+1])
					if !tmp.end && tmp.typ == "a" {
						tmp.index[0] = i
						tmp.index[1] = j
						nd := fndend(tmp, in[j+1:])
						if !reflect.DeepEqual(nd, tag{}) {
							tmp.Meat = in[j+1 : nd.index[0]]
							out = append(out, tmp)
							str := in[tmp.index[1]:nd.index[0]]
							in = in[:i] + str + in[nd.index[1]+1:]
						}
					}
					break
				}
			}
		}
	}
	return
}

func fndend(fnt tag, area string) tag {
	var count int
	for i, v := range area {
		if v == '<' {
			for j, val := range area[i:] {
				if val == '>' {
					var tmp tag
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
	return tag{}
}
