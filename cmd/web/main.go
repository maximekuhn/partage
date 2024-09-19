package main

import (
	"github.com/maximekuhn/partage/internal/app/web"
)

func main() {
	// FIXME: passing a nil will break, doing this for tmp testing...
	s := web.NewServer(nil)
	s.Run()
}
