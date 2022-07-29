package Util

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"runtime"

	log "github.com/cihub/seelog"
	"bytes"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
)

// 产生panic时的调用栈打印
// 调用方法:在主循环头 defer PrintPanicStack()
func PrintPanicStack(extras ...interface{}) string {
	var buff bytes.Buffer
	var haveErr = false
	if x := recover(); x != nil {
		haveErr = true
		buff.WriteString(fmt.Sprintf("dump:%v\n", x))
		//log.Error("dump:%v", x)
		i := 0
		funcName, file, line, ok := runtime.Caller(i)
		for ok {
			buff.WriteString(fmt.Sprintf("F%v:[%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line))
			//log.Error("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			i++
			funcName, file, line, ok = runtime.Caller(i)
		}

		for k := range extras {
			buff.WriteString(fmt.Sprintf("EXRAS#%v DATA:%v", k, spew.Sdump(extras[k])))
		}
	}
	if haveErr {
		log.Error(buff.String())
		return buff.String()
	}
	return ""
}

// 获取调用栈
func GetPanicStackString() string {
	var buff bytes.Buffer
	i := 0
	funcName, file, line, ok := runtime.Caller(i)
	for ok {
		buff.WriteString(fmt.Sprintf("F%v:[%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line))
		i++
		funcName, file, line, ok = runtime.Caller(i)
	}
	return buff.String()

}

func Dump(objs ...interface{}) string {
	return spew.Sdump(objs...)
}

func DumpJson(obj interface{}) string {
	s,_ := json.Marshal(obj)
	return string(s)
}