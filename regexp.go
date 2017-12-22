package onigmo

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -lonigmo
#include <stdlib.h>
#include <string.h>
#include <onigmo.h>

// cgo does not support vargs
extern int onigmo_helper_error_code_with_info_to_str(UChar* err_buf, int err_code, OnigErrorInfo *errInfo);
extern int onigmo_helper_error_code_to_str(UChar* err_buf, OnigPosition err_code);
extern char* onigmo_helper_get(char* str, OnigPosition* beg, OnigPosition* end, int n);

int onigmo_helper_error_code_with_info_to_str(UChar* err_buf, int err_code, OnigErrorInfo *errInfo) {
    return onig_error_code_to_str(err_buf, err_code, errInfo);
}
int onigmo_helper_error_code_to_str(UChar* err_buf, OnigPosition err_code) {
    return onig_error_code_to_str(err_buf, err_code);
}
char* onigmo_helper_get(char* str, OnigPosition* beg, OnigPosition* end, int n) {
    char *res = strndup((char*)(str + beg[n]), end[n] - beg[n]);
    return res;
}
*/
import "C"

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"unsafe"
)

var (
	ONIG_ENCODING_UTF8         = &C.OnigEncodingUTF_8
	ONIG_ENCODING_ASCII        = &C.OnigEncodingASCII
	ONIG_SYNTAX_PERL           = &C.OnigSyntaxPerl
	ONIG_SYNTAX_POSIX_BASIC    = &C.OnigSyntaxPosixBasic
	ONIG_SYNTAX_POSIX_EXTENDED = &C.OnigSyntaxPosixExtended
)

type Regexp struct {
	encoding               C.OnigEncoding
	regex                  C.OnigRegex
	cachedCaptureGroupNums map[string][]C.int
	mu                     sync.Mutex
	matched                bool
	region                 *C.OnigRegion
	input                  string
}

func OnigmoVersion() string {
	return C.GoString(C.onig_version())
}

func NewRegexp(str string) (*Regexp, error) {
	ret := C.onig_init()
	if ret != 0 {
		return nil, errors.New("failed to initialize encoding for the Onigumo regular expression library.")
	}
	result := &Regexp{
		encoding:               ONIG_ENCODING_UTF8,
		cachedCaptureGroupNums: make(map[string][]C.int),
	}
	result.mu.Lock()
	defer result.mu.Unlock()

	patternStart, patternEnd := stringPointers(str)
	defer free(patternStart, patternEnd)

	var errorInfo C.OnigErrorInfo
	r := C.onig_new(&result.regex, patternStart, patternEnd, C.ONIG_OPTION_DEFAULT, result.encoding, C.ONIG_SYNTAX_DEFAULT, &errorInfo)
	if r != C.ONIG_NORMAL {
		return nil, errors.New(errMsgWithInfo(r, &errorInfo))
	}

	return result, nil
}

func Compile(str string) (*Regexp, error) {
	return NewRegexp(str)
}

func MustCompile(str string) *Regexp {
	regexp, error := Compile(str)
	if error != nil {
		panic(`regexp: Compile(` + quote(str) + `): ` + error.Error())
	}

	return regexp
}

func Match(pattern string, s string) (bool, error) {
	re, err := Compile(pattern)
	if err != nil {
		return false, err
	}

	return re.Match(s)
}

func (re *Regexp) HasCaptureGroup(name string) bool {
	_, err := re.getCaptureGroupNums(name)

	return err == nil
}

func (re *Regexp) getCaptureGroupNums(name string) ([]C.int, error) {
	cached, ok := re.cachedCaptureGroupNums[name]
	if ok {
		return cached, nil
	}

	nameStart, nameEnd := stringPointers(name)
	defer free(nameStart, nameEnd)

	var groupNums *C.int
	n := C.onig_name_to_group_numbers(re.regex, nameStart, nameEnd, &groupNums)
	if n <= 0 {
		return nil, fmt.Errorf("%v: no such capture group in pattern", name)
	}

	result := make([]C.int, 0, int(n))
	for i := 0; i < int(n); i++ {
		result = append(result, getPos(groupNums, C.int(i)))
	}

	re.cachedCaptureGroupNums[name] = result

	return result, nil
}

func (re *Regexp) Match(input string) (bool, error) {
	region := C.onig_region_new()
	inputStart, inputEnd := stringPointers(input)
	defer free(inputStart, inputEnd)

	r := C.onig_match(re.regex, inputStart, inputEnd, inputStart, region, C.ONIG_OPTION_NONE)
	if r == C.ONIG_MISMATCH {
		C.onig_region_free(region, 1)
		re.matched = false
		return false, nil

	} else if r < 0 {
		C.onig_region_free(region, 1)
		return false, errors.New(errMsg(r))

	} else {
		re.matched = true
		re.region = region
		re.input = input
		return true, nil
	}
}

func (re *Regexp) Get(name string) (string, error) {
	if !re.matched {
		return "", nil
	}

	groupNums, err := re.getCaptureGroupNums(name)
	if err != nil {
		return "", err
	}

	for _, groupNum := range groupNums {
		w := C.onigmo_helper_get(C.CString(re.input), re.region.beg, re.region.end, groupNum)
		return C.GoString(w), nil
	}

	return "", nil
}

func (re *Regexp) Free() {
	if re.matched {
		C.onig_region_free(re.region, 1)
	}
	C.onig_free(re.regex)
}

func quote(s string) string {
	if strconv.CanBackquote(s) {
		return "`" + s + "`"
	}
	return strconv.Quote(s)
}

func stringPointers(s string) (start, end *C.OnigUChar) {
	start = (*C.OnigUChar)(unsafe.Pointer(C.CString(s)))
	end = (*C.OnigUChar)(unsafe.Pointer(uintptr(unsafe.Pointer(start)) + uintptr(len(s))))
	return
}

func getPos(p *C.int, i C.int) C.int {
	return *(*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + uintptr(i)*unsafe.Sizeof(C.int(0))))
}

func free(start *C.OnigUChar, end *C.OnigUChar) {
	C.memset(unsafe.Pointer(start), C.int(0), C.size_t(uintptr(unsafe.Pointer(end))-uintptr(unsafe.Pointer(start))))
	C.free(unsafe.Pointer(start))
}

func errMsgWithInfo(returnCode C.int, errorInfo *C.OnigErrorInfo) string {
	msg := make([]byte, C.ONIG_MAX_ERROR_MESSAGE_LEN)
	l := C.onigmo_helper_error_code_with_info_to_str((*C.UChar)(&msg[0]), returnCode, errorInfo)
	if l <= 0 {
		return "unknown error"
	} else {
		return string(msg[:l])
	}
}

func errMsg(returnCode C.OnigPosition) string {
	msg := make([]byte, C.ONIG_MAX_ERROR_MESSAGE_LEN)
	l := C.onigmo_helper_error_code_to_str((*C.UChar)(&msg[0]), returnCode)
	if l <= 0 {
		return "unknown error"
	} else {
		return string(msg[:l])
	}
}
