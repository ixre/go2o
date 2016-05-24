package main

import (
	"flag"
)

func main() {
	var areaId int
	var urlAddr string

	areaId = flag.IntVar(&areaId, "a", -1, "id of area")
	urlAddr = flag.StringVar(&urlAddr, "u", "", "the url contain 58 city page")

	if areaId > 0 && len(urlAddr) > 0 {
		collect(areaId, urlAddr)
	}
}

func collect(areaId int, urlAddr string) {

}
