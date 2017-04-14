package appimg

type appimg struct {
	full    string
	name    string
	version string
}

func newApp(name string) appimg {
	var out appimg
	out.full = name
	return out
}
