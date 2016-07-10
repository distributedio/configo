package configo

import (
	"reflect"
)

type Validator interface {
	Validate(value string) error
}

//numeric range, format: (min,max),(min,max],[min,max),[min,max], no spaces in format
type nrange struct {
	min int64
	max int64

	left  bool //min value include
	right bool //max value include

	kind reflect.Kind
}

func (n *nrange) Validate(value string) error {
	return nil
}

//regex match, format: /expression/
type regex struct {
	exp string
}

//named validator, format: name
type named struct {
}
