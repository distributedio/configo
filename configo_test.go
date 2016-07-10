package configo

import (
	//"github.com/naoina/toml"
	"fmt"
	"testing"
)

type Config struct {
	Listen  string   `cfg:"listen; :8805; netaddr; listen address of server"`
	Cluster []string `cfg:"cluster; ['127.0.0.1:8800']; dialstring"`
	MaxConn int      `cfg:"max-conn; required; numeric"`
}

func TestUnmarshal(t *testing.T) {
	//data := []byte("listen = ':8804'")
	data := []byte("max-conn=100")
	c := &Config{}
	if err := Unmarshal(data, c); err != nil {
		t.Error(err)
	}
	t.Logf("Unmarshal result %v\n", c)
}
func TestMarshal(t *testing.T) {
	c := &Config{}
	if data, err := Marshal(c); err != nil {
		t.Error(err)
	} else {
		fmt.Printf("%s", data)
	}
}
