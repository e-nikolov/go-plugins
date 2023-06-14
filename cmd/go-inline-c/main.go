package main

/*
#include <stdio.h>

void test() {
	printf("Printed from inline C\n");
}
*/
import "C"

func main() {
	C.test()
}
