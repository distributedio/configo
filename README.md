conf by golang struct tag

```go
//`cfg: "name, required or default value, validate, descripion"`
type Config struct {
    Listen string `cfg: "listen, required, netaddr, server listen address"`
    MaxConns int `cfg: "max-conns, 1000, , max number of connections`
}
```
