package configo

import (
	"reflect"
	"testing"
	"time"

	"github.com/shafreeck/toml"
	"github.com/stretchr/testify/require"
)

func TestTravel_Travel(t *testing.T) {
	tests := []struct {
		name string
		obj  interface{}
		want map[string]interface{}
	}{
		{
			name: "test string",
			obj: struct {
				name string `cfg:"name;;;user name"`
			}{
				name: "name-value",
			},
			want: map[string]interface{}{
				"name": "name-value",
			},
		},
		{
			name: "test int,float,bool",
			obj: struct {
				name  string  `cfg:"name;;;user name"`
				age   int     `cfg:"age;;;user age"`
				score float64 `cfg:"score;;;user score"`
				sex   bool    `cfg:"sex;;;user sex"`
			}{
				name:  "user-name",
				age:   18,
				score: 60.1,
				sex:   true,
			},
			want: map[string]interface{}{
				"name":  "user-name",
				"age":   18,
				"score": 60.1,
				"sex":   true,
			},
		},
		{
			name: "test slice",
			obj: struct {
				Cluster []string `cfg:"cluster;;;the address of redis cluster"`
			}{
				Cluster: []string{"127.0.0.1:6379", "127.0.0.1:7379"},
			},
			want: map[string]interface{}{
				"cluster":   []string{"127.0.0.1:6379", "127.0.0.1:7379"},
				"cluster.0": "127.0.0.1:6379",
				"cluster.1": "127.0.0.1:7379",
			},
		},
		{
			name: "test array",
			obj: struct {
				Cluster [2]string `cfg:"cluster;;;the address of redis cluster"`
			}{
				Cluster: [2]string{"127.0.0.1:6379", "127.0.0.1:7379"},
			},
			want: map[string]interface{}{
				"cluster":   [2]string{"127.0.0.1:6379", "127.0.0.1:7379"},
				"cluster.0": "127.0.0.1:6379",
				"cluster.1": "127.0.0.1:7379",
			},
		},
		{
			name: "test time duration",
			obj: struct {
				Timeout time.Duration `cfg:"timeout;10s;;time out"`
			}{
				Timeout: 3 * time.Second,
			},
			want: map[string]interface{}{
				"timeout": 3 * time.Second,
			},
		},
		{
			name: "test struct in struct",
			obj: struct {
				nest struct {
					Timeout time.Duration `cfg:"timeout;10s;;time out"`
				} `cfg:"nest;;;nest struct"`
			}{
				nest: struct {
					Timeout time.Duration `cfg:"timeout;10s;;time out"`
				}{
					Timeout: 3 * time.Second,
				},
			},
			want: map[string]interface{}{
				"nest.timeout": 3 * time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handle := func(path string, tag *toml.CfgTag, v reflect.Value) {
				if w, ok := tt.want[path]; ok {
					wt := reflect.ValueOf(w)
					testEqual(t, wt, v)
					delete(tt.want, path)
				} else {
					require.FailNowf(t, "unknow path", "path=%s, v=%v", path, v)
				}
			}
			tr := NewTravel(handle)
			tr.Travel(tt.obj)

			if len(tt.want) > 0 {
				require.FailNowf(t, "not get some path", "%#v", tt.want)
			}
		})
	}
}

func testEqual(t *testing.T, want, get reflect.Value) {
	require.Equal(t, want.Kind(), get.Kind(), "data kind not same")
	switch get.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		require.Equal(t, want.Int(), get.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		require.Equal(t, want.Uint(), get.Uint())
	case reflect.Uintptr:
	case reflect.Float32, reflect.Float64:
		require.Equal(t, want.Float(), get.Float())
	case reflect.Bool:
		require.Equal(t, want.Bool(), get.Bool())
	case reflect.String:
		require.Equal(t, want.String(), get.String())

	case reflect.Slice, reflect.Array:
		require.Equal(t, want.Len(), get.Len())
		for i := 0; i < want.Len(); i++ {
			testEqual(t, want.Index(i), get.Index(i))
		}
	default:
		t.Fatalf("uset type %v", get.Kind())
	}
}
