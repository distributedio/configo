package main

import (
	"log"

	"github.com/distributedio/configo"
	"github.com/distributedio/configo/example/conf"
)

func main() {
	c := &conf.Config{}

	if err := configo.Load("conf/example.toml", c); err != nil {
		log.Fatalln(err)
	}
	log.Println(c)
}
