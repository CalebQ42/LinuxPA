package main

type portableApp struct {
	exes []executable
	name string
	icon []byte
	cats []string
}

func (p portableApp) exesNoWine() (out []executable) {
	for _, e := range p.exes {
		if !e.wine {
			out = append(out, e)
		}
	}
	return
}
