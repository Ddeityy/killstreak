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
	"fmt"
	"log"
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

func main() {
	log.Println(RustParseDemo("/home/deity/Dev/killstreak/test_data/unix/2023-11-11_22-30-00.dem"))
	RustCutDemo("/home/deity/Dev/killstreak/test_data/unix/2023-11-11_22-30-00.dem", fmt.Sprint(30756))
}
