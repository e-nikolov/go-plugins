# go-plugins

This repo contains some examples of alternatives to the plugin package in Go's standard library, inspired by https://github.com/vladimirvivien/go-cshared-examples.

## Dynamic loading of Go plugins in a Go program via CGO

- Run the demo via:


```bash

go build -o build/c-shared-plug.so -buildmode=c-shared ./plugins/c-shared-plug/main.go
go run ./cmd/go-dlopen/
```


- Main program: https://github.com/e-nikolov/go-plugins/blob/main/cmd/go-dlopen/main.go

```Go
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
```

- Plugin: https://github.com/e-nikolov/go-plugins/blob/main/plugins/c-shared-plug/main.go

```Go
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
```

- Plugin Loader in C


```C
#include "main.h"
#include <stdio.h>
#include <dlfcn.h>
#include <dlfcn.h>
#include <stdlib.h>
#include <string.h>

go_int (*add)(go_int, go_int);

go_int Add(go_int a, go_int b) {
	return add(a, b);
}

go_float64 (*cosine)(go_float64);

go_float64 Cosine(go_float64 x) {
	return cosine(x);
}

void (*sort)(go_slice);

void Sort(go_slice vals) {
	sort(vals);
}

go_int (*logg)(go_str);

go_int Log(go_str msg) {
    char *c_msg = malloc(msg.len + 10);
    memcpy(c_msg, msg.p, msg.len);
    memcpy(c_msg + msg.len, "; from C", 8);
    go_str c_msg_str = {c_msg, msg.len + 8};
    go_int ret = logg(c_msg_str);
    free(c_msg);
    return ret;
}

char* (*startHTTPServer)(char* addr, char* message);

char* StartHTTPServer(char* addr, char* message) {
    return startHTTPServer(addr, message);
}

int LoadPlugin(go_str path) {
    void *handle;
    char *error;

    // use dlopen to load shared object
    handle = dlopen(path.p, RTLD_LAZY);
    if (!handle) {
        fputs (dlerror(), stderr);
        exit(1);
    }

    // resolve Add symbol and assign to fn ptr
    add = dlsym(handle, "Add");
    if ((error = dlerror()) != NULL)  {
        fputs(error, stderr);
        exit(1);
    }

    // resolve Cosine symbol
    cosine = dlsym(handle, "Cosine");
    if ((error = dlerror()) != NULL)  {
        fputs(error, stderr);
        exit(1);
    }
    // resolve Sort symbol
    sort = dlsym(handle, "Sort");
    if ((error = dlerror()) != NULL)  {
        fputs(error, stderr);
        exit(1);
    }

    // resolve Log symbol
    logg = dlsym(handle, "Log");
    if ((error = dlerror()) != NULL)  {
        fputs(error, stderr);
        exit(1);
    }
    // resolve StartHTTPServer symbol
    startHTTPServer = dlsym(handle, "StartHTTPServer");
    if ((error = dlerror()) != NULL)  {
        fputs(error, stderr);
        exit(1);
    }

    // close file handle when done
    dlclose(handle);

    return 0;
}
```
