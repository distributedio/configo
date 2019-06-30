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
		/*
			vValue := v.Elem()
			copyValue := reflect.New(vValue.Type()).Elem()
			travel(copyValue, vValue)
			copy.Set(copyValue)
		*/
	case reflect.Struct:
		for i := 0; i < v.NumField(); i += 1 {
			t.travel(path, v.Field(i))
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i += 1 {
			p := fmt.Sprintf("%s.%d.", path, i)
			t.travel(p, v.Index(i))
		}
	case reflect.Map:
		/*
			for _, key := range v.MapKeys() {
				vValue := v.MapIndex(key)
				copyValue := reflect.New(vValue.Type()).Elem()
				t.travel(copyValue, vValue)
				copy.SetMapIndex(key, copyValue)
			}
		*/
	case reflect.String:
		// TODO get path
		t.handle(path, v)
	default:
		//copy.Set(v)
	}
}
