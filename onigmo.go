package main

/*
#cgo pkg-config: onigmo
#include <stdlib.h>
#include <string.h>
#include <onigmo.h>
#include "gonigmo.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"sync"
	"unsafe"
)

func main() {
	fmt.Println(Version())

	expr := "a(.*)b|[e-f]+"
	regex := MustCompile(expr)
	fmt.Printf("%#v\n", regex)

	boo := regex.MatchString("zzzzaffffffffb")
	if boo {
		fmt.Println("matched!!!!!!")
	} else {
		fmt.Println("no matched......")
	}
}

func Version() string {
	return C.GoString(C.onig_version())
}

type Matches struct {
	count   int
	indexes [][]int32
}

type Regexp struct {
	expr      string
	regex     C.OnigRegex
	region    *C.OnigRegion
	errorInfo *C.OnigErrorInfo
	errorBuf  *C.char
	matches   *Matches
	mu        sync.Mutex
}

func compile(expr string, option int) (*Regexp, error) {
	re := &Regexp{expr: expr}
	exprPtr := C.CString(expr)
	defer C.free(unsafe.Pointer(exprPtr))

	re.mu.Lock()
	defer re.mu.Unlock()

	r := C.GonigNewDefault(exprPtr, C.int(len(expr)), C.int(option), &re.regex, &re.region, &re.errorInfo, &re.errorBuf)
	if r != C.ONIG_NORMAL {
		return nil, errors.New(C.GoString(re.errorBuf))
	}

	return re, nil
}

func Compile(expr string) (*Regexp, error) {
	return compile(expr, C.ONIG_OPTION_DEFAULT)
}

func MustCompile(str string) *Regexp {
	regexp, error := Compile(str)
	if error != nil {
		panic(`regexp: Compile(` + quote(str) + `): ` + error.Error())
	}
	return regexp
}

func (re *Regexp) Clear() {
	matches := re.matches
	matches.count = 0
}

func (re *Regexp) doMatch(b []byte, n int, offset int) bool {
	re.Clear()
	if n == 0 {
		b = []byte{0}
	}

	ptr := unsafe.Pointer(&b[0])
	r := C.GonigSearch((ptr), C.int(n), C.int(offset), C.int(C.ONIG_OPTION_DEFAULT), re.regex,
		re.region, re.errorInfo, (*C.char)(nil), (*C.int)(nil), (*C.int)(nil))

	pos := (int)(r)
	return pos >= 0
}

func (re *Regexp) Match(b []byte) bool {
	return re.doMatch(b, len(b), 0)
}

func (re *Regexp) MatchString(s string) bool {
	return re.Match([]byte(s))
}

/*
type Matches struct {
	match  bool
	regex  *Regexp
	region *C.OnigRegion
	input  string
}

func (re *Regexp) Free() {
	C.onig_free(re.regex)
}

func (re *Regexp) Encoding() C.OnigEncoding {
	//return &C.OnigEncodingASCII
	return &C.OnigEncodingUTF_8
}

func (re *Regexp) Init() (C.OnigEncoding, error) {
	encoding := re.Encoding()
	encodings := []C.OnigEncoding{encoding}
	r := C.onig_initialize(&encodings[0], C.int(len(encodings)))
	if r != C.ONIG_NORMAL {
		return nil, errors.New("failed to onigmo initialization")
	}
	return encoding, nil
}

func (re *Regexp) Match(input string) (*Matches, error) {
	region := C.onig_region_new()

	start, end := pointers(input)
	defer free(start, end)

	r := C.onig_match(re.regex, start, end, start, region, C.ONIG_OPTION_NONE)
	if r == C.ONIG_MISMATCH {
		C.onig_region_free(region, 1)
		return &Matches{
			match: false,
			regex: re,
			input: input,
		}, nil

	} else if r < 0 {
		C.onig_region_free(region, 1)
		return nil, errors.New("error")

	} else {
		return &Matches{
			match:  true,
			regex:  re,
			region: region,
			input:  input,
		}, nil
	}
}

func (m *Matches) IsMatch() bool {
	return m.match
}

func (m *Matches) Free() {
	if m.match {
		C.onig_region_free(m.region, 1)
	}
}

func pointers(s string) (start, end *C.OnigUChar) {
	start = (*C.OnigUChar)(unsafe.Pointer(C.CString(s)))
	end = (*C.OnigUChar)(unsafe.Pointer(uintptr(unsafe.Pointer(start)) + uintptr(len(s))))
	return
}

func free(start *C.OnigUChar, end *C.OnigUChar) {
	C.memset(unsafe.Pointer(start), C.int(0), C.size_t(uintptr(unsafe.Pointer(end))-uintptr(unsafe.Pointer(start))))
	C.free(unsafe.Pointer(start))
}
*/
