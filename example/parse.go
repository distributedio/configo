package main

import (
	"fmt"

	"github.com/shafreeck/configo"
)

type Config struct {
	Listen   string `cfg:"listen, :8804, netaddr, server listen address"`
	MaxConns int    `cfg:"maxconns, 1000, numeric, max number of connections"`
}

func main() {
	data := []byte("")
	v := &Config{}
	if err := configo.Unmarshal(data, v); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(v)
}
