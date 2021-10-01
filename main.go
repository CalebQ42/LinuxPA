package main

import (
	"embed"
)

var (
	version = "3.0.0.0"
	apps    []*app
	//go:embed embed
	embedFS embed.FS
)

func main() {
}
