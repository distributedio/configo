package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/distributedio/configo"
	"github.com/distributedio/configo/example/conf"
)

func main() {
	c := &conf.Config{}

	configo.Flags(c, "listen", "redis", "redis.0.cluster")
	flag.Parse()

	data, err := ioutil.ReadFile("conf/example.toml")
	if err != nil {
		log.Fatalln(err)
	}

	err = configo.Unmarshal(data, c)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%#v\n", c)
}
