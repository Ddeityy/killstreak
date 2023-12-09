package main

/*
#cgo LDFLAGS: ${SRCDIR}/lib/rust.so
#include <stdlib.h>
#include "./lib/rust.h"
*/
import "C"
import (
	"log"
	"os"
	"path"
	"strings"
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
	WatchDemosDir()
}

func formatDemos() {
	demosDir, err := GetDemosDir()
	if err != nil {
		panic(err)
	}
	demos, _ := os.ReadDir(demosDir)
	for _, demo := range demos {
		if strings.Contains(demo.Name(), ".dem") {
			log.Println("------------------------------------------------")
			log.Println("Processing:", demo.Name())
			ProcessDemo(path.Join(demosDir, demo.Name()), demosDir)
		}
	}
}
