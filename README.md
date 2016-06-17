# Configo is a go library to parse toml configuration using struct tags

## Toml
[shafreeck/toml](https://github.com/shafreeck/toml) is a modification version of [shafreeck](naoina/toml), Add the abililty to parse complex struct tags

## Validate
configo use [govalidate](https://github.com/asaskevich/govalidator) to validate conf values

## Features
* Configure parser behaviour using struct tags
* Set default value or required
* Validate value using govalidate(Validators with parameters will be supported soon)
* Generate toml template basing on go struct and tags(TODO)

## Example
```go
package main

import (
        "fmt"

        "github.com/shafreeck/configo"
)

//`cfg: "name, required or default value, validate, descripion"`
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
```
