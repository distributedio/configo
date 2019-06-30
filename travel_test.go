package configo

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTravel_Travel(t *testing.T) {
	tests := []struct {
		name   string
		obj    interface{}
		handle TravelHandle
	}{
		{
			name: "test string",
			obj: struct {
				name string `cfg:"name;;;user name"`
			}{
				name: "name",
			},
			want: map[string]reflect.Vlaue{
				"name": a,
			},
			handle: func(path string, v reflect.Value) {
				fmt.Printf("paht=%s, v=%#v\n", path, v)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := NewTravel(tt.handle)
			tr.Travel(tt.obj)
		})
	}
}
