package main

import "C"

import (
	"fmt"
	"math"
	"net/http"
	"sort"
	"sync"

	"github.com/e-nikolov/scratch/pkg/shared"
)

var count int
var mtx sync.Mutex

func init() {
	shared.GlobalVariable = "Does init() get called?"
}

//export Add
func Add(a, b int) int {
	fmt.Printf("Go Plugin - %q\n", shared.GlobalVariable)

	shared.GlobalVariable = "Value Modified By a Plugin"

	fmt.Printf("Go Plugin - %q\n", shared.GlobalVariable)

	return a + b
}

//export Cosine
func Cosine(x float64) float64 {
	return math.Cos(x)
}

//export Sort
func Sort(vals []int) {
	sort.Ints(vals)
}

//export Log
func Log(msg string) int {
	mtx.Lock()
	defer mtx.Unlock()
	fmt.Println(msg + "; from Go")
	count++
	return count
}

//export StartHTTPServer
func StartHTTPServer(cAddr *C.char, cMessage *C.char) (err *C.char) {
	addr := C.GoString(cAddr)
	message := C.GoString(cMessage)

	fmt.Printf("\nStarting HTTP Server at %q\n", addr)
	e := http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "HTTP Server Demo Inside a Go plugin\n\nMesage from the main program: %q\nPath: %q\n", message, r.URL.Path)
	}))

	if e != nil {
		return C.CString(e.Error())
	}

	return nil
}

func main() {}
