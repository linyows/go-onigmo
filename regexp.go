package onigmo

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -lonigmo
#include <stdlib.h>
#include <string.h>
#include <onigmo.h>

// cgo does not support vargs
extern int onig_helper_error_code_with_info_to_str(UChar* err_buf, int err_code, OnigErrorInfo *errInfo);
extern int onig_helper_error_code_to_str(UChar* err_buf, OnigPosition err_code);
int onig_helper_error_code_with_info_to_str(UChar* err_buf, int err_code, OnigErrorInfo *errInfo) {
    return onig_error_code_to_str(err_buf, err_code, errInfo);
}
int onig_helper_error_code_to_str(UChar* err_buf, OnigPosition err_code) {
    return onig_error_code_to_str(err_buf, err_code);
}
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

type Regexp struct {
	regex                  C.OnigRegex
	cachedCaptureGroupNums map[string][]C.int
}

type MatchResult struct {
	match  bool
	regex  *Regexp
	region *C.OnigRegion
	input  string
}

func OnigmoVersion() string {
	return C.GoString(C.onig_version())
}

func Compile(pattern string) (*Regexp, error) {
	ret := C.onig_init()
	if ret != 0 {
		return nil, errors.New("failed to initialize encoding for the Oniguruma regular expression library.")
	}
	result := &Regexp{
		cachedCaptureGroupNums: make(map[string][]C.int),
	}
	patternStart, patternEnd := pointers(pattern)
	defer free(patternStart, patternEnd)
	var errorInfo C.OnigErrorInfo
	r := C.onig_new(&result.regex, patternStart, patternEnd, C.ONIG_OPTION_DEFAULT, &C.OnigEncodingASCII, C.ONIG_SYNTAX_DEFAULT, &errorInfo)
	if r != C.ONIG_NORMAL {
		return nil, errors.New(errMsgWithInfo(r, &errorInfo))
	}
	return result, nil
}

func (regex *Regexp) Free() {
	C.onig_free(regex.regex)
}

func (regex *Regexp) HasCaptureGroup(name string) bool {
	_, err := regex.getCaptureGroupNums(name)
	return err == nil
}

func (r *Regexp) getCaptureGroupNums(name string) ([]C.int, error) {
	cached, ok := r.cachedCaptureGroupNums[name]
	if ok {
		return cached, nil
	}
	nameStart, nameEnd := pointers(name)
	defer free(nameStart, nameEnd)
	var groupNums *C.int
	n := C.onig_name_to_group_numbers(r.regex, nameStart, nameEnd, &groupNums)
	if n <= 0 {
		return nil, fmt.Errorf("%v: no such capture group in pattern", name)
	}
	result := make([]C.int, 0, int(n))
	for i := 0; i < int(n); i++ {
		result = append(result, getPosI(groupNums, C.int(i)))
	}
	r.cachedCaptureGroupNums[name] = result
	return result, nil
}

func (regex *Regexp) Match(input string) (*MatchResult, error) {
	region := C.onig_region_new()
	inputStart, inputEnd := pointers(input)
	defer free(inputStart, inputEnd)
	r := C.onig_match(regex.regex, inputStart, inputEnd, inputStart, region, C.ONIG_OPTION_NONE)
	if r == C.ONIG_MISMATCH {
		C.onig_region_free(region, 1)
		return &MatchResult{
			match: false,
		}, nil
	} else if r < 0 {
		C.onig_region_free(region, 1)
		return nil, errors.New(errMsg(r))
	} else {
		return &MatchResult{
			match:  true,
			regex:  regex,
			region: region,
			input:  input,
		}, nil
	}
}

func (m *MatchResult) Get(name string) (string, error) {
	if !m.match {
		return "", nil
	}
	groupNums, err := m.regex.getCaptureGroupNums(name)
	if err != nil {
		return "", err
	}
	for _, groupNum := range groupNums {
		beg := getPos(m.region.beg, groupNum)
		end := getPos(m.region.end, groupNum)
		if beg > end || beg < 0 || int(end) > len(m.input) {
			return "", fmt.Errorf("%v: unexpected result when calling onig_name_to_group_numbers()", name)
		} else if beg == end {
			continue
		} else {
			return m.input[beg:end], nil
		}
	}
	return "", nil
}

func (m *MatchResult) IsMatch() bool {
	return m.match
}

func (m *MatchResult) Free() {
	if m.match {
		C.onig_region_free(m.region, 1)
	}
}

func pointers(s string) (start, end *C.OnigUChar) {
	start = (*C.OnigUChar)(unsafe.Pointer(C.CString(s)))
	end = (*C.OnigUChar)(unsafe.Pointer(uintptr(unsafe.Pointer(start)) + uintptr(len(s))))
	return
}

func getPos(p *C.OnigPosition, i C.int) C.int {
	return *(*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + uintptr(i)*unsafe.Sizeof(C.int(0))))
}

func getPosI(p *C.int, i C.int) C.int {
	return *(*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + uintptr(i)*unsafe.Sizeof(C.int(0))))
}

func free(start *C.OnigUChar, end *C.OnigUChar) {
	C.memset(unsafe.Pointer(start), C.int(0), C.size_t(uintptr(unsafe.Pointer(end))-uintptr(unsafe.Pointer(start))))
	C.free(unsafe.Pointer(start))
}

func errMsgWithInfo(returnCode C.int, errorInfo *C.OnigErrorInfo) string {
	msg := make([]byte, C.ONIG_MAX_ERROR_MESSAGE_LEN)
	l := C.onig_helper_error_code_with_info_to_str((*C.UChar)(&msg[0]), returnCode, errorInfo)
	if l <= 0 {
		return "unknown error"
	} else {
		return string(msg[:l])
	}
}

func errMsg(returnCode C.OnigPosition) string {
	msg := make([]byte, C.ONIG_MAX_ERROR_MESSAGE_LEN)
	l := C.onig_helper_error_code_to_str((*C.UChar)(&msg[0]), returnCode)
	if l <= 0 {
		return "unknown error"
	} else {
		return string(msg[:l])
	}
}
