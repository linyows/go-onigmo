package main

/*
#cgo LDFLAGS: -lonigmo
#include <onigmo.h>
#include "go_onigmo.c"
*/
import "C"
import "fmt"

func main() {
	r, _ := C.go_onigmo()
	fmt.Printf("%#v\n", r)
	s, _ := C.example()
	fmt.Printf("%#v\n", s)
}
