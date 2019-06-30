package configo

import (
	"fmt"
	"reflect"
	"testing"

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handle := func(path string, v reflect.Value) {
				fmt.Printf("paht=%s, v=%#v\n", path, v)
				if w, ok := tt.want[path]; ok {
					require.Equal(t, reflect.ValueOf(w), v)
					delete(tt.want, path)
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
