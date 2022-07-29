package stress

import (
	"reflect"
	"runtime"
	"strings"
)

var (
	tests = map[string]func(GoroutineIndex, _RunIndex){}
)

func register(f func(GoroutineIndex, _RunIndex)) {
	funName := getFunctionName(f, '.')
	if _, ok := tests[funName]; ok {
		panic(funName + " already registered")
	}
	tests[funName] = f
}

func getFunctionName(i interface{}, seps ...rune) string {
	// 获取函数名称
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()

	// 用 seps 进行分割
	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})

	// fmt.Println(fields)

	if size := len(fields); size > 0 {
		return fields[size-1]
	}
	return ""
}
