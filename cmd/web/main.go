package main

import (
	"log"

	"github.com/maximekuhn/partage/internal/app/web"
)

func main() {
	// FIXME: passing a nil will break, doing this for tmp testing...
	s := web.NewServer(nil)
	log.Fatal(s.Run())
}
