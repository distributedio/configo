package main

import (
	"fmt"

	"github.com/shafreeck/configo"
)

type Config struct {
	Listen string `cfg:"listen, :8805, netaddr, listen address of server"`
	Redis  struct {
		Cluster []string `cfg:"cluster, ['127.0.0.1:8800'], dialstring"`
		Net     struct {
			Timeout int
		}
	}
	MaxConn int `cfg:"max-conn, required, numeric"`
}

func main() {
	c := &Config{}
	if data, err := configo.Marshal(c); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Printf("%s", data)
	}
}
