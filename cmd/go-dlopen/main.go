//go:build (linux && cgo) || (darwin && cgo) || (freebsd && cgo)

package main

// #cgo CFLAGS: -g -Wall
// #include "main.h"
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"math/rand"
	"unsafe"

	"github.com/e-nikolov/scratch/pkg/shared"
)

func main() {
	fmt.Printf("Main Go - %q\n", shared.GlobalVariable)

	LoadPlugin("./build/c-shared-plug.so")

	a, b := rand.Int()%100, rand.Int()%100
	fmt.Printf("%d + %d = %d\n", a, b, Add(a, b))

	c := rand.Float64()
	fmt.Printf("cos(%.5f) = %.5f\n", c, Cosine(c))

	vals := rand.Perm(10)
	fmt.Println(vals)
	Sort(vals)
	fmt.Println(vals)

	Log("Hello from Go")

	fmt.Printf("Main Go - %q\n", shared.GlobalVariable)

	fmt.Printf("\n%v\n", StartHTTPServer(":8080", "Hello from the main Go"))
}

func LoadPlugin(path string) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	C.LoadPlugin(C.go_str{cpath, C.longlong(len(path))})
}

func Add(a, b int) int {
	return int(C.Add(C.longlong(a), C.longlong(b)))
}

func Cosine(x float64) float64 {
	return float64(C.Cosine(C.double(x)))
}

func Sort(vals []int) {
	C.Sort(C.go_slice{unsafe.Pointer(&vals[0]), C.longlong(len(vals)), C.longlong(cap(vals))})
}

func Log(msg string) int {
	cs := C.CString(msg)
	defer C.free(unsafe.Pointer(cs))
	return int(C.Log(C.go_str{cs, C.longlong(len(msg))}))
}

func StartHTTPServer(addr string, message string) error {
	caddr := C.CString(addr)
	defer C.free(unsafe.Pointer(caddr))

	cmsg := C.CString(message)
	defer C.free(unsafe.Pointer(cmsg))

	e := C.StartHTTPServer(caddr, cmsg)

	er := C.GoString(e)
	// defer C.free(unsafe.Pointer(e.p))

	if er != "" {
		return fmt.Errorf(er)
	}

	return nil
}
