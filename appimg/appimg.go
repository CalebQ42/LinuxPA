package appimg

type appimg struct {
	name string
}

func newApp(name string) appimg {
	var out appimg
	out.name = name
	return out
}
