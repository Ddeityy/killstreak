package main

/*
#cgo LDFLAGS: ${SRCDIR}/lib/rust.so
#include <stdlib.h>
#include "./lib/rust.h"
*/
import "C"
import (
	"unsafe"
)

// Parses demo and returns a JSON string to unmarshal
func RustParseDemo(demoPath string) string {
	cDemoPath := C.CString(demoPath)
	defer C.free(unsafe.Pointer(cDemoPath))
	o := C.parse_demo(cDemoPath)
	return C.GoString(o)
}

func main() {
	//WatchDemosDir()
	FormatDemos()
}
