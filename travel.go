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
		fmt.Printf("travel ptr")
		vValue := v.Elem()
		if !vValue.IsValid() {
			return
		}
		t.travel(path, vValue)
	case reflect.Interface:
		fmt.Printf("travel interface")
		/*
			vValue := v.Elem()
			copyValue := reflect.New(vValue.Type()).Elem()
			travel(copyValue, vValue)
			copy.Set(copyValue)
		*/
	case reflect.Struct:
		fmt.Printf("travel strcut")
		for i := 0; i < v.NumField(); i += 1 {
			if !v.Field(i).IsValid() {
				continue
			}
			tag := extractTag(v.Type().Field(i).Tag.Get(fieldTagName))
			t.travel(path+"."+tag.Name, v.Field(i))
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
		fmt.Printf("default set value %#v\n", v)
		//copy.Set(v)
	}
}
