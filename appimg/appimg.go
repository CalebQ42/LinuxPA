package appimg

type appimg struct {
	url  string
	name string
}

func newApp(url, name string) appimg {
	var out appimg
	out.url = url
	out.name = name
	return out
}
