package main

/*
#include <stdio.h>
#include "foo.c"
*/
import "C"

func main() {
	C.ACFunction()
}
