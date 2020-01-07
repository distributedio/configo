package configo

import (
	"fmt"
	"reflect"

	"github.com/shafreeck/toml"
)

type TravelHandle func(path string, tag *toml.CfgTag, v reflect.Value)

type Travel struct {
	handle TravelHandle
}

func NewTravel(h TravelHandle) *Travel {
	return &Travel{handle: h}
}

func (t *Travel) Travel(obj interface{}) {
	t.travel("", nil, reflect.ValueOf(obj))
}

func (t *Travel) travel(path string, tag *toml.CfgTag, v reflect.Value) {
	switch v.Kind() {
	case reflect.Ptr:
		vValue := v.Elem()
		if !vValue.IsValid() {
			return
		}
		t.travel(path, tag, vValue)
	case reflect.Interface:
		vValue := v.Elem()
		if !vValue.IsValid() {
			return
		}
		t.travel(path, tag, vValue)
	case reflect.Struct:
		for i := 0; i < v.NumField(); i += 1 {
			if !v.Field(i).IsValid() {
				continue
			}
			tag := extractTag(v.Type().Field(i).Tag.Get(fieldTagName))
			p := tag.Name
			if len(path) > 0 {
				p = path + "." + tag.Name
			}
			t.travel(p, tag, v.Field(i))
		}
	case reflect.Slice, reflect.Array:
		// handle slice & array as a whole
		t.handle(path, tag, v)
		for i := 0; i < v.Len(); i++ {
			p := fmt.Sprintf("%d", i)
			if len(path) > 0 {
				p = path + "." + p
			}
			// handle every element
			t.travel(p, tag, v.Index(i))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fallthrough
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fallthrough
	case reflect.Float32, reflect.Float64:
		fallthrough
	case reflect.Bool:
		fallthrough
	case reflect.String:
		t.handle(path, tag, v)
	default:
		panic(fmt.Sprintf("config file use unsupport type. %v", v.Type()))
	}
}
