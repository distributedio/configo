package conf

import "time"

type Config struct {
	Listen  string   `cfg:"listen; :8804; netaddr; The address the server to listen"`
	MaxConn int      `cfg:"max-connection; 10000; numeric; Max number of concurrent connections"`
	Redis   [2]Redis `cfg:"redis;;;redis cluster"`
}

type Redis struct {
	Cluster []string      `cfg:"cluster; ['127.0.0.1:6379'];; the addresses of redis cluster"`
	Timeout time.Duration `cfg:"timeout; 10s"`
}
