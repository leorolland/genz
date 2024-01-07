package main

import (
	"log"

	"github.com/leorolland/genz/cmd/genz"
)

func main() {
	if err := genz.Execute(); err != nil {
		log.Fatal(err)
	}
}
