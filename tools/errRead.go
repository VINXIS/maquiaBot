package tools

import (
	"log"
	"runtime"
)

// ErrRead will check to see if there is an error; it will print the error if there is any
func ErrRead(err error) {
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Fatalf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
	return
}
