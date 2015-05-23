package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func ReportError(f string, a ...interface{}) {
	_, path, line, _ := runtime.Caller(1)
	_, file := filepath.Split(path)
	f = fmt.Sprintf("%v:%v", file, line) + f
	fmt.Fprintf(os.Stderr, f, a...)
}
