package configo

import (
	"fmt"
	"reflect"
)

type TravelHandle func(path string, v reflect.Value)

type Travel struct {
	handle TravelHandle
}

func NewTravel(h TravelHandle) *Travel {
	return &Travel{handle: h}
}

func (t *Travel) Travel(obj interface{}) {
	t.travel("", reflect.ValueOf(obj))
}

func (t *Travel) travel(path string, v reflect.Value) {
	switch v.Kind() {
	case reflect.Ptr:
		vValue := v.Elem()
		if !vValue.IsValid() {
			return
		}
		t.travel(path, vValue)
	case reflect.Interface:
		vValue := v.Elem()
		if !vValue.IsValid() {
			return
		}
		t.travel(path, vValue)
	case reflect.Struct:
		for i := 0; i < v.NumField(); i += 1 {
			if !v.Field(i).IsValid() {
				continue
			}
			p := getPath(path, v.Type().Field(i))
			t.travel(p, v.Field(i))
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			p := fmt.Sprintf("%d", i)
			if len(path) > 0 {
				p = path + "." + p
			}
			t.travel(p, v.Index(i))
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
		t.handle(path, v)
	default:
		panic(fmt.Sprintf("config file use unsupport type. %v", v.Type()))
	}
}

func getPath(prefix string, f reflect.StructField) string {
	tag := extractTag(f.Tag.Get(fieldTagName))
	if len(prefix) == 0 {
		return tag.Name
	}
	return prefix + "." + tag.Name
}
