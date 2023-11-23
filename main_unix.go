//go:build linux
// +build linux

package main

/*
#cgo LDFLAGS: ${SRCDIR}/lib/rust.so
#include <stdlib.h>
#include "./lib/rust.h"
*/
import "C"
import (
	"flag"
	"unsafe"
)

// Parses demo and returns a JSON string to unmarshal
func RustParseDemo(demoPath string) string {
	cDemoPath := C.CString(demoPath)
	defer C.free(unsafe.Pointer(cDemoPath))
	o := C.parse_demo(cDemoPath)
	return C.GoString(o)
}

// Cuts demo and outputs a cut_demoName.dem
func RustCutDemo(demoPath string, startTick string) {
	cDemoPath := C.CString(demoPath)
	defer C.free(unsafe.Pointer(cDemoPath))

	cStartTick := C.CString(startTick)
	defer C.free(unsafe.Pointer(cStartTick))

	C.cut_demo(cDemoPath, cStartTick)
}

var cut bool

func main() {
	autoCut := flag.Bool("cut", true, "Automatically cut the demo")
	flag.Parse()
	cut = *autoCut
	WatchDemosDir()
}
