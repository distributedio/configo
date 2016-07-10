package conf

type Config struct {
	Listen  string `cfg:"listen; :8804; netaddr; The address the server to listen"`
	MaxConn int    `cfg:"max-connection; 10000; numeric; Max number of concurrent connections"`
	Redis   struct {
		Cluster []string `cfg:"cluster; required; dialstring; The addresses of redis cluster"`
	}
}
