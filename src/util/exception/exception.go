package exception

import (
	"runtime"
	"fmt"
	"util/logger"
)

func CatchExecption()  {
	var trace = make([]byte, 1024, 1024)
	if err := recover(); err != nil {
		count := runtime.Stack(trace, true)
		errMsg := fmt.Sprintf("Recover from panic: %s\n", err)
		traceMsg := fmt.Sprintf("Stack of %d bytes: %s\n", count, trace)
		logger.Error.Println(errMsg)
		logger.Error.Println(traceMsg)
		return
	}
}