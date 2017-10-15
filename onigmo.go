package main

/*
#cgo pkg-config: onigmo
#include <stdio.h>
#include <stdlib.h>
#include <onigmo.h>
*/
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

func main() {
	// fmt.Println(Version())
	Simple()
}

type Regexp struct {
	pattern  string
	reg      C.regex_t
	einfo    *C.OnigErrorInfo
	errorBuf *C.char
}

func Simple() {
	var r int
	p := "a(.*)b|[e-f]+"
	re := &Regexp{pattern: p}
	// region *C.OnigRegion
	// s := "zzzzaffffffffb"
	pattern := C.CString(p)
	defer C.free(unsafe.Pointer(pattern))
	// str := C.CString(s)
	plen := C.int(len(p))

	r = C.onig_new(&re.reg, pattern, plen, C.ONIG_OPTION_DEFAULT, 0, C.ONIG_SYNTAX_DEFAULT, &re.einfo)
	if r != C.ONIG_NORMAL {
		fmt.Fprintf(os.Stdout, "ERROR: %d\n", r)
		// var strerr C.OnigUChar
		// C.onig_error_code_to_str(strerr, r, &re.einfo)
		// fmt.Fprintf(os.Stdout, "ERROR: %s\n", C.GoString(strerr))
	}
}

func Version() string {
	return C.GoString(C.onig_version())
}
