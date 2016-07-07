package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/shafreeck/configo"
	"github.com/shafreeck/configo/example/conf"
)

func main() {
	c := &conf.Config{}

	data, err := ioutil.ReadFile("conf/example.toml")
	if err != nil {
		log.Fatalln(err)
	}

	err = configo.Unmarshal(data, c)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(c)
}
