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
//`cfg: "name, required or default value, validate, descripion"`
type Config struct {
    Listen string `cfg: "listen, required, netaddr, server listen address"`
    MaxConns int `cfg: "max-conns, 1000, , max number of connections`
}
```
