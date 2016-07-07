package main

import (
	"log"

	"github.com/shafreeck/configo"
	"github.com/shafreeck/configo/example/conf"
)

func main() {
	c := conf.Config{}

	if err := configo.Dump("dump.toml", c); err != nil {
		log.Fatalln(err)
	}
	log.Println("Dumpped to dump.toml")
}
