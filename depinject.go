package gosdk

import (
	"errors"
	"fmt"
	"reflect"
)

// 依赖注入

type DIContainer struct {
	depMap map[reflect.Type]reflect.Value // 依赖映射表
	err    error                          // 错误
}

func (dic *DIContainer) Err() error {
	return dic.err
}

// 尝试更新依赖表
func (dic *DIContainer) updateDepMap(retval reflect.Value) {
	rt_retval := retval.Type()
	// 如果当前返回值类型为error并且值不为nil
	if rt_retval.AssignableTo(reflect.TypeFor[error]()) {
		if retval.Interface() != nil {
			dic.err = retval.Interface().(error)
			return
		}
	} else {
		// 判断是否已经存在对应类型
		_, ok := dic.depMap[rt_retval]
		if ok {
			fmt.Println(rt_retval)
			dic.err = errors.New("type allready exists")
			return
		}
		dic.depMap[rt_retval] = retval
	}
}

// 提供待注入项（方法，返回值允许有error）以及依赖项（指针类型）
// 待注入项经过依赖项注入后可以生成新依赖项
func Inject(v ...any) *DIContainer {
	dic := &DIContainer{depMap: make(map[reflect.Type]reflect.Value)} // 依赖注入容器
	funcs := []reflect.Value{}                                        // 等待注入依赖的方法

	// 先遍历一遍处理纯依赖项和0参数函数，保留其余函数到funcs
	for _, value := range v {
		if value == nil {
			continue
		}
		rv := reflect.ValueOf(value)
		if rv.Kind() == reflect.Ptr && rv.Elem().Kind() == reflect.Struct {
			// 如果是结构体指针
			// 判断是否已经存在对应类型
			dic.updateDepMap(rv)
			if dic.err != nil {
				return dic
			}
		} else if rv.Kind() == reflect.Func {
			// 如果是函数
			// 参数数量为0就先执行
			if rv.Type().NumIn() == 0 {
				retvals := rv.Call(nil)
				for _, retval := range retvals {
					dic.updateDepMap(retval)
					if dic.err != nil {
						return dic
					}
				}
			} else {
				funcs = append(funcs, rv)
			}
		}
		// 如果都不是直接跳过
	}

	// 不断遍历funcs
	for {
		newfuncs := []reflect.Value{} // 剩余的待注入func
		for _, fn := range funcs {
			deps := []reflect.Value{} // 依赖项
			inject := false           // 当前存在的依赖项是否能够注入这个方法了
			// 遍历参数类型切片
			for intype := range fn.Type().Ins() {
				dep, ok := dic.depMap[intype]
				// 如果找不到这个类型的依赖项就跳出
				if !ok {
					inject = false
					break
				}
				deps = append(deps, dep)
				inject = true
			}
			// 如果能够注入
			if inject {
				// 执行并尝试更新依赖表
				retvals := fn.Call(deps)
				for _, retval := range retvals {
					dic.updateDepMap(retval)
					if dic.err != nil {
						return dic
					}
				}
			} else {
				newfuncs = append(newfuncs, fn)
			}
		}
		// 一次遍历后剩余项数量不变，说明已经可以退出了，剩下的都是找不到依赖项的
		if len(newfuncs) == len(funcs) || len(newfuncs) == 0 {
			break
		}
		funcs = newfuncs // 原先使用copy，忘了短copy到长，长度还是长
	}
	return dic
}

// 对全部依赖项进行操作（并不是再次注入，只是进行与指定依赖项相关的操作而已）
func (dic *DIContainer) Invoke(funcs ...any) *DIContainer {
	if dic.err != nil {
		return dic
	}
	for _, fn := range funcs {
		if IsFunction(fn) {
			rv_fn := reflect.ValueOf(fn)
			deps := []reflect.Value{} // 依赖
			inject := false           // 依赖表是否含有全部所需依赖
			// 检查依赖参数是否全部拥有
			for in := range rv_fn.Type().Ins() {
				rv_dep, ok := dic.depMap[in]
				if !ok {
					inject = false
					break
				}
				deps = append(deps, rv_dep)
				inject = true
			}
			// 如果能注入
			if inject {
				rv_fn.Call(deps) // TODO 需要检测执行结果是否含有error并不为nil
			}
		}
	}
	return dic
}
