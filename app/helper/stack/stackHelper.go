package stack

import (
	"fmt"
	"github.com/golang/glog"
	"runtime"
)

func PrintRecoverFromPanic(isPanic ...*bool) {
	var panicStack []interface{}
	if err := recover(); err != nil {
		panicStack = append(panicStack, fmt.Sprintf("panicStack Recover:%s↓\n", err))
		i := 0
		funcName, file, line, ok := runtime.Caller(i)
		for ok {
			panicStack = append(panicStack, fmt.Sprintf("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line))
			i++
			funcName, file, line, ok = runtime.Caller(i)
			if len(isPanic) != 0 {
				*isPanic[0] = true
			}
		}
		glog.Errorf("%+v", panicStack)
		glog.Flush()
	}
}
