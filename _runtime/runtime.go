package _runtime

import (
	"bytes"
	"runtime"
	"strconv"
	_ "unsafe"
)

//go:linkname runtime_procPin runtime.procPin
func runtime_procPin() int

//go:linkname runtime_procUnpin runtime.procUnpin
func runtime_procUnpin() int

func GetProcID() int {
	pid := runtime_procPin()
	runtime_procUnpin()
	return pid
}

func GetGoID() int64 {
	var buf [64]byte
	// 得到的字符串类似于："goroutine 6 [running]:..."
	s := buf[:runtime.Stack(buf[:], false)]
	// 掐头去尾
	s = s[len("goroutine "):]
	s = s[:bytes.IndexByte(s, ' ')]

	gid, _ := strconv.ParseInt(string(s), 10, 64)
	return gid
}
