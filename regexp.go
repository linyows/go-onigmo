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
	encoding    C.OnigEncoding
	regex       C.OnigRegex
	mu          sync.Mutex
	expr        string
	matchResult *MatchResult
}

type NamedGroupNums map[string][]C.int

type MatchResult struct {
	regex          C.OnigRegex
	matched        bool
	input          string
	region         *C.OnigRegion
	namedGroupNums NamedGroupNums
	errorMessage   string
}

func OnigmoVersion() string {
	return C.GoString(C.onig_version())
}

func NewRegexp(expr string) (*Regexp, error) {
	ret := C.onig_init()
	if ret != 0 {
		return nil, errors.New("failed to initialize encoding for the Onigumo regular expression library.")
	}

	re := &Regexp{
		encoding: ONIG_ENCODING_UTF8,
		expr:     expr,
	}

	re.mu.Lock()
	defer re.mu.Unlock()

	beginning, end := stringPointers(expr)
	defer free(beginning, end)

	var errorInfo C.OnigErrorInfo
	r := C.onig_new(&re.regex, beginning, end, C.ONIG_OPTION_DEFAULT, re.encoding, C.ONIG_SYNTAX_DEFAULT, &errorInfo)
	if r != C.ONIG_NORMAL {
		return nil, errors.New(errMsgWithInfo(r, &errorInfo))
	}

	return re, nil
}

func Compile(expr string) (*Regexp, error) {
	return NewRegexp(expr)
}

func MustCompile(expr string) *Regexp {
	regexp, error := Compile(expr)
	if error != nil {
		panic(`regexp: Compile(` + quote(expr) + `): ` + error.Error())
	}
	return regexp
}

func Match(pattern string, b []byte) bool {
	re, err := Compile(pattern)
	if err != nil {
		return false
	}
	return re.match(b)
}

func MatchString(pattern string, s string) bool {
	re, err := Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(s)
}

func (m *MatchResult) HasCaptureGroup(name string) bool {
	_, err := m.getNamedGroupNums(name)

	return err == nil
}

func (re *Regexp) match(b []byte) bool {
	region := C.onig_region_new()
	beginning, end := bytePointers(b)
	defer free(beginning, end)
	input := string(b)

	r := C.onig_match(re.regex, beginning, end, beginning, region, C.ONIG_OPTION_NONE)
	if r == C.ONIG_MISMATCH {
		C.onig_region_free(region, 1)
		re.matchResult = &MatchResult{
			matched:        false,
			input:          input,
			namedGroupNums: make(map[string][]C.int),
		}
		return false

	} else if r < 0 {
		C.onig_region_free(region, 1)
		re.matchResult = &MatchResult{
			matched:        false,
			input:          input,
			namedGroupNums: make(map[string][]C.int),
			errorMessage:   errMsg(r),
		}
		return false

	} else {
		re.matchResult = &MatchResult{
			matched:        true,
			region:         region,
			input:          input,
			regex:          re.regex,
			namedGroupNums: make(map[string][]C.int),
		}
		return true
	}
}

func (re *Regexp) MatchString(s string) bool {
	b := []byte(s)
	return re.match(b)
}

func (re *Regexp) search(b []byte) bool {
	region := C.onig_region_new()
	beginning, end := bytePointers(b)
	searchBeginning := beginning
	searchEnd := end
	defer free(beginning, end)
	input := string(b)

	r := C.onig_search(re.regex, beginning, end, searchBeginning, searchEnd, region, C.ONIG_OPTION_NONE)
	if r == C.ONIG_MISMATCH {
		C.onig_region_free(region, 1)
		re.matchResult = &MatchResult{
			matched:        false,
			input:          input,
			namedGroupNums: make(map[string][]C.int),
		}
		return false

	} else if r < 0 {
		C.onig_region_free(region, 1)
		re.matchResult = &MatchResult{
			matched:        false,
			input:          input,
			namedGroupNums: make(map[string][]C.int),
			errorMessage:   errMsg(r),
		}
		return false

	} else {
		re.matchResult = &MatchResult{
			matched:        true,
			region:         region,
			input:          input,
			regex:          re.regex,
			namedGroupNums: make(map[string][]C.int),
		}
		return true
	}
}

func (re *Regexp) SearchString(s string) bool {
	b := []byte(s)
	return re.search(b)
}

func (m *MatchResult) getNamedGroupNums(s string) ([]C.int, error) {
	cached, ok := m.namedGroupNums[s]
	if ok {
		return cached, nil
	}

	beginning, end := stringPointers(s)
	defer free(beginning, end)

	var groupNums *C.int
	n := C.onig_name_to_group_numbers(m.regex, beginning, end, &groupNums)
	if n <= 0 {
		return nil, fmt.Errorf("%v: no such capture group in pattern", s)
	}

	result := make([]C.int, 0, int(n))
	for i := 0; i < int(n); i++ {
		result = append(result, getPos(groupNums, C.int(i)))
	}

	m.namedGroupNums[s] = result

	return result, nil
}

func (m *MatchResult) Get(s string) (string, error) {
	if !m.matched {
		return "", nil
	}

	groupNums, err := m.getNamedGroupNums(s)
	if err != nil {
		return "", err
	}

	for _, groupNum := range groupNums {
		w := C.onigmo_helper_get(C.CString(m.input), m.region.beg, m.region.end, groupNum)
		word := C.GoString(w)
		if word != "" {
			return word, nil
		}
	}

	return "", nil
}

func (re *Regexp) Free() {
	C.onig_free(re.regex)
}

func (m *MatchResult) Free() {
	if m.matched {
		C.onig_region_free(m.region, 1)
	}
}

func quote(s string) string {
	if strconv.CanBackquote(s) {
		return "`" + s + "`"
	}
	return strconv.Quote(s)
}

func stringPointers(s string) (beginning, end *C.OnigUChar) {
	beginning = (*C.OnigUChar)(unsafe.Pointer(C.CString(s)))
	end = (*C.OnigUChar)(unsafe.Pointer(uintptr(unsafe.Pointer(beginning)) + uintptr(len(s))))
	return
}

func bytePointers(b []byte) (beginning, end *C.OnigUChar) {
	beginning = (*C.OnigUChar)(C.CBytes(b))
	end = (*C.OnigUChar)(unsafe.Pointer(uintptr(unsafe.Pointer(beginning)) + uintptr(len(b))))
	return
}

func getPos(p *C.int, i C.int) C.int {
	return *(*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + uintptr(i)*unsafe.Sizeof(C.int(0))))
}

func free(beginning *C.OnigUChar, end *C.OnigUChar) {
	C.memset(unsafe.Pointer(beginning), C.int(0), C.size_t(uintptr(unsafe.Pointer(end))-uintptr(unsafe.Pointer(beginning))))
	C.free(unsafe.Pointer(beginning))
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
