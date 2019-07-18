package configo

import (
	"flag"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/shafreeck/toml"
)

/*
	`flags` 实现将对象中的变量添加到`flag`中，从而实现通过命令行设置变量的功能。

		import (
			"log"

			"github.com/distributedio/configo"
		)

		type Config struct {
			Key       string   `cfg:"key; default;; simple type example"`
			Child     *Child   `cfg:"child; ;; class type "`
			Array     []string `cfg:"array;;; simple array type"`
			CompArray []*Child `cfg:"comp;;; complex array type"`
		}

		type Child struct {
			Name string `cfg:"name; noname;; child class item`
		}

		func main() {
			conf := &Config{}
			configo.AddFlags(conf)
			flag.Parse()

			if err := configo.Load("conf/example.toml", conf); err != nil {
				log.Fatalln(err)
			}
		}

	首先，需要在`flag.Parse()`之前调用`AddFlags()`将对象中的变量添加到`falg`中。
	`configo.Load()`会在内部调用`ApplyFlags()`方法，将`flag`中设置的变量应用到
	对象中。

	对象中的变量按照如下规则对应`flag`中的`key`：

	* 简单数据类型，直接使用`cfg`中的`name`作为`flag`中的`key`。
	  如`Conf.Key`，对应`flag`中的`key`。
	* 对象数据类型，需要添加上一层对象的名称。
	  如 `Conf.Child.Name` 对应`flag`中的`child.name`
	* 数组或slice类型，要增加下标作为一个层级。
	  如 `Conf.CompArray[0].Name`，对应`flag`中的`comp.0.name`
	* 对于简单数据类型的数组或slice也可以使用名称作为`flag`中的`key`，
	  使用字符串表示一个数组。
	  例如：`Conf.Array`，对应`flag`中的`array`。同时在执行中，使用如下的
	  方式设置`array`:
		./cmd -array="[\"a1\", \"a2\"]"
*/

const (
	ConfigoFlagSuffix = "[configo]"
)

// AddFlags 将对象中的变量加入到flag中，从而可以通过命令行设置对应的变量。
//
// * `obj`  为待加入到`flag`中的对象的实例
// * `keys` 限定加入`flag`中变量的范围，**不设置**的时候表示将所有变量都加入到`flag`中。
func AddFlags(fs *flag.FlagSet, obj interface{}, keys ...string) {
	flagMap := make(map[string]struct{}, len(keys))
	for i := range keys {
		flagMap[keys[i]] = struct{}{}
	}
	t := NewTravel(func(path string, tag *toml.CfgTag, fv reflect.Value) {
		if _, ok := flagMap[path]; len(flagMap) > 0 && !ok {
			return
		}
		var err error
		switch fv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64:
			var v int64
			if v, err = strconv.ParseInt(tag.Value, 10, 64); err != nil {
				if fv.Kind() == reflect.Int64 {
					//try to parse a time.Duration
					if d, err := time.ParseDuration(tag.Value); err == nil {
						fs.Duration(path, time.Duration(d), tag.Description+ConfigoFlagSuffix)
						return
					}
				}
				log.Fatalln(err)
				return
			}
			fs.Int64(path, v, tag.Description+ConfigoFlagSuffix)
		case reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64:
			var v uint64
			if v, err = strconv.ParseUint(tag.Value, 10, 64); err != nil {
				log.Fatalln(err)
				return
			}
			fs.Uint64(path, v, tag.Description+ConfigoFlagSuffix)
		case reflect.Float32, reflect.Float64:
			var v float64
			if v, err = strconv.ParseFloat(tag.Value, 64); err != nil {
				log.Fatalln(err)
				return
			}
			fs.Float64(path, v, tag.Description+ConfigoFlagSuffix)
		case reflect.Bool:
			var v bool
			if v, err = strconv.ParseBool(tag.Value); err != nil {
				log.Fatalln(err)
				return
			}
			fs.Bool(path, v, tag.Description+ConfigoFlagSuffix)
		case reflect.String:
			fs.String(path, tag.Value, tag.Description+ConfigoFlagSuffix)
		case reflect.Slice, reflect.Array:
			// TODO 使用flag.Var设置变量
			fs.String(path, tag.Value, tag.Description+ConfigoFlagSuffix)
		default:
			log.Printf("unknow type %s for set flag", fv.Type())
		}
	})
	t.Travel(obj)
}

// ApplyFlags 将命令行中设置的变量值应用到`obj`中。
//
// **注意：** configo中的函数默认会调用这个函数设置配置文件，所以不需要显示调用。
func ApplyFlags(fs *flag.FlagSet, obj interface{}) {
	actualFlags := make(map[string]*flag.Flag)
	fs.Visit(func(f *flag.Flag) {
		if strings.Contains(f.Usage, ConfigoFlagSuffix) {
			actualFlags[f.Name] = f
		}
	})
	t := NewTravel(func(path string, tag *toml.CfgTag, fv reflect.Value) {
		f, ok := actualFlags[path]
		if !ok {
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
			// TODO NOT support
			//	if err := unmarshalArray("name", f.Value.String(), &s); err != nil {
			//		log.Fatalln(err)
			//		return
			//	}
			//	fv.Set(reflect.ValueOf(s.Name))
			//	log.Printf("get list =%#v\n", s)
		default:
			log.Printf("unknow type %s for set flag", fv.Type())
		}
	})
	t.Travel(obj)
}
