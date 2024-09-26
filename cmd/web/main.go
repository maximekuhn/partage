package main

import (
	"log"

	"github.com/maximekuhn/partage/internal/app/web"
)

func main() {
	s := web.NewServer()
	log.Fatal(s.Run())
}
