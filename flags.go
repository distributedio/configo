package configo

// Flags 将类中的变量加入到flags中，从而可以通过命令行进行设置
// 可以通过keys指定范围的成员加入到flags中，如果不指定则将所有的
// 变量都加入到flags中
// 不同类型变量的书写规则
// 1. 一级变量，即只有一层变量
// 2. 多级变量：使用 root.child 的形式作为变量名称
// 3. 数组变量：数组下标作为一个层级，例如要设置root[0].key, 则flags中的key名称为root.0.key
func Flags(obj interface{}, keys ...string) {
	/*
		t := newTravel(func(path string, v reflect.Value) {
			if len(keys) > 0 && !match(path, keys) {
				return
			}
			switch v.Kind() {
			case Int:
				value := flags.IntVal()
				addtomap
			default:
			}
		})
		t.travel(obj)
	*/

	// 1. 遍历所有变量
	// 2. 如果keys没有设置，则默认是将所有的变量都加入到flags中
	// 3. 如果指定了具体的变量名称，则只加入指定的变量到flags中
	// 4. 在读取配置文件之后，添加完默认配置之后，应用命令行中的配置
}

func ApplyFlags(obj interface{}) {
	/*
		actualFlags := make(map[string]bool)
		flag.Visit(func(f *flag.Flag) {
			actualFlags[f.Name] = f
		})
		if len(actualFlags) == 0 {
			return
		}
		t := newTravel(func(path string, v reflect.Value) {
			if _, ok := actualFlags[path]; !ok {
				return
			}

			if len(keys) > 0 && !match(path, keys) {
				return
			}
			switch v.Kind() {
			case Int:
				value := flags.IntVal()
				addtomap
			default:
			}
		})
		t.travel(obj)
	*/
}
