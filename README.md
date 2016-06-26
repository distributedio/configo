# Configo
Configo is a go library to parse toml configuration using struct tags

## Toml
[shafreeck/toml](https://github.com/shafreeck/toml) is a modification version of [naoina/toml](https://github.com/naoina/toml),
adding the abililty to parse complex struct tags.

## Validate
configo use [govalidator](https://github.com/asaskevich/govalidator) to validate conf values

## Features
* Configure parser behaviour using struct tags
* Set default value or required
* Validate value using govalidator(Validators with parameters will be supported soon)
* Generate toml template basing on go struct and tags
* Build generating tools using configo-build

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

### First, install configo-build
```
go get github.com/shafreeck/configo/bin/configo-build
```

### Second, build an executable program basing on your package and struct
```
configo-build <package>.<struct>
```
### Finally, use the built program to generate a toml
```
<built program> > conf.toml
```
and you can patch you toml file if it is already existed
```
<built program> -patch conf.toml
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
