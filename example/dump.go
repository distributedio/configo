package main

import (
	"log"

	"github.com/distributedio/configo"
	"github.com/distributedio/configo/example/conf"
)

func main() {
	c := conf.Config{}

	if err := configo.Dump("dump.toml", c); err != nil {
		log.Fatalln(err)
	}
	log.Println("Dumpped to dump.toml")
}
