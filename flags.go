package configo

import (
	"flag"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/shafreeck/toml"
)

var flagMap = make(map[string]struct{})

// Flags 将类中的变量加入到flags中，从而可以通过命令行进行设置
// 可以通过keys指定范围的成员加入到flags中，如果不指定则将所有的
// 变量都加入到flags中
// 不同类型变量的书写规则
// 1. 一级变量，即只有一层变量
// 2. 多级变量：使用 root.child 的形式作为变量名称
// 3. 数组变量：数组下标作为一个层级，例如要设置root[0].key, 则flags中的key名称为root.0.key
func Flags(obj interface{}, keys ...string) {
	for i := range keys {
		flagMap[keys[i]] = struct{}{}
	}
	t := NewTravel(func(path string, tag *toml.CfgTag, fv reflect.Value) {
		if _, ok := flagMap[path]; len(flagMap) > 0 && !ok {
			return
		}
		var err error
		//Set the default value
		switch fv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64:
			var v int64
			if v, err = strconv.ParseInt(tag.Value, 10, 64); err != nil {
				if fv.Kind() == reflect.Int64 {
					//try to parse a time.Duration
					if d, err := time.ParseDuration(tag.Value); err == nil {
						flag.Duration(path, time.Duration(d), tag.Description)
						return
					}
				}
				log.Fatalln(err)
				return
			}
			flag.Int64(path, v, tag.Description)
		case reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64:
			var v uint64
			if v, err = strconv.ParseUint(tag.Value, 10, 64); err != nil {
				log.Fatalln(err)
				return
			}
			flag.Uint64(path, v, tag.Description)
		case reflect.Float32, reflect.Float64:
			var v float64
			if v, err = strconv.ParseFloat(tag.Value, 64); err != nil {
				log.Fatalln(err)
				return
			}
			flag.Float64(path, v, tag.Description)
		case reflect.Bool:
			var v bool
			if v, err = strconv.ParseBool(tag.Value); err != nil {
				log.Fatalln(err)
				return
			}
			flag.Bool(path, v, tag.Description)
		case reflect.String:
			flag.String(path, tag.Value, tag.Description)
		case reflect.Slice, reflect.Array:
			// set as string type
			flag.String(path, tag.Value, tag.Description)
		default:
			log.Fatalf("unsupport type %s for set flag", fv.Type())
		}
	})
	t.Travel(obj)

	// 1. 遍历所有变量
	// 2. 如果keys没有设置，则默认是将所有的变量都加入到flags中
	// 3. 如果指定了具体的变量名称，则只加入指定的变量到flags中
	// 4. 在读取配置文件之后，添加完默认配置之后，应用命令行中的配置
}

func ApplyFlags(obj interface{}) {
	actualFlags := make(map[string]*flag.Flag)
	flag.Visit(func(f *flag.Flag) {
		actualFlags[f.Name] = f
	})
	if len(actualFlags) == 0 {
		return
	}
	t := NewTravel(func(path string, tag *toml.CfgTag, fv reflect.Value) {
		f, ok := actualFlags[path]
		if !ok {
			return
		}
		if _, ok := flagMap[path]; len(flagMap) > 0 && !ok {
			return
		}
		var err error
		switch fv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64:
			var v int64
			if v, err = strconv.ParseInt(f.Value.String(), 10, 64); err != nil {
				if fv.Kind() == reflect.Int64 {
					//try to parse a time.Duration
					if d, err := time.ParseDuration(f.Value.String()); err == nil {
						fv.SetInt(int64(d))
						return
					}
				}
				log.Fatalln(err)
				return
			}
			fv.SetInt(v)
		case reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64:
			var v uint64
			if v, err = strconv.ParseUint(f.Value.String(), 10, 64); err != nil {
				log.Fatalln(err)
				return
			}
			fv.SetUint(v)
		case reflect.Float32, reflect.Float64:
			var v float64
			if v, err = strconv.ParseFloat(f.Value.String(), 64); err != nil {
				log.Fatalln(err)
				return
			}
			fv.SetFloat(v)
		case reflect.Bool:
			var v bool
			if v, err = strconv.ParseBool(f.Value.String()); err != nil {
				log.Fatalln(err)
				return
			}
			fv.SetBool(v)
		case reflect.String:
			fv.SetString(f.Value.String())
		case reflect.Slice, reflect.Array:
			// FIXME TODO
			/*
				v := rv.Addr().Interface()
				if err := unmarshalArray(ft.Name, f.Value.String(), v); err != nil {
					return err
				}
			*/
		default:
			log.Fatalf("unsupport type %s for set flag", fv.Type())
		}
	})
	t.Travel(obj)
}
