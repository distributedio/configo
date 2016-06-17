package configo

import (
	"testing"
)

type Config struct {
	Listen  string   `cfg:"listen, :8805, netaddr, listen address of server"`
	Cluster []string `cfg:"cluster, ['127.0.0.1:8800'], dialstring"`
}

func TestUnmarshal(t *testing.T) {
	//data := []byte("listen = ':8804'")
	data := []byte("")
	c := &Config{}
	if err := Unmarshal(data, c); err != nil {
		t.Error(err)
	}
	t.Logf("Unmarshal result %s\v\n", c)
}
