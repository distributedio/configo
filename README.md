# Configo is a go library to parse toml configuration using struct tags

## Toml
[shafreeck/toml](https://github.com/shafreeck/toml) is a modification version of [naoina/toml](https://github.com/naoina/toml),
adding the abililty to parse complex struct tags.

## Validate
configo use [govalidate](https://github.com/asaskevich/govalidator) to validate conf values

## Features
* Configure parser behaviour using struct tags
* Set default value or required
* Validate value using govalidate(Validators with parameters will be supported soon)
* Generate toml template basing on go struct and tags

## Example

### Unmarshal toml to config
```go
package main

import (
        "fmt"

        "github.com/shafreeck/configo"
)

//`cfg:"name, required or default value, validate, descripion"`
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

## Generate toml from struct
```go
package main

import (
        "fmt"

        "github.com/shafreeck/configo"
)

type Config struct {
        Listen string `cfg:"listen, :8805, netaddr, listen address of server"`
        MaxConn int `cfg:"max-conn, required, numeric"`
        Redis  struct {
                Cluster []string `cfg:"cluster, ['127.0.0.1:8800'], dialstring"`
                Net     struct {
                        Timeout int
                }
        }
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
```
Output

```toml
#type:        string
#rules:       netaddr
#description: listen address of server
#default:     :8805
#listen=""

#type:        int
#rules:       numeric
#required
max-conn=0

[redis]

#type:        []string
#rules:       dialstring
#default:     ['127.0.0.1:8800']
#cluster=[]

[redis.net]

#type:        int
timeout=0
```
