package main

import (
	"fmt"
	"log"

	"github.com/shafreeck/configo"
	"github.com/shafreeck/configo/example/conf"
)

func main() {
	c := conf.Config{}

	data, err := configo.Marshal(c)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(data))
}
