package appimg

import "strings"

func removeLetters(vers string) string {
	vers = strings.ToLower(vers)
	letters := []string{"abcdefghijklmnopqrstuvwxyz"}
	for _, v := range letters {
		vers = strings.Replace(vers, v, "", -1)
	}
	return vers
}
