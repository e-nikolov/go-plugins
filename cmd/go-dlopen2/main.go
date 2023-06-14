//go:build (linux && cgo) || (darwin && cgo) || (freebsd && cgo)

package main

/*
#cgo linux LDFLAGS: -ldl
#include <dlfcn.h>
#include <limits.h>
#include <stdlib.h>
#include <stdint.h>

#include <stdio.h>

static uintptr_t pluginOpen(const char* path, char** err) {
	void* h = dlopen(path, RTLD_NOW|RTLD_GLOBAL);
	if (h == NULL) {
		*err = (char*)dlerror();
	}
	return (uintptr_t)h;
}

static void* pluginLookup(uintptr_t h, const char* name, char** err) {
	void* r = dlsym((void*)h, name);
	if (r == NULL) {
		*err = (char*)dlerror();
	}
	return r;
}
*/
import "C"
import (
	"errors"
	"fmt"
	"sync"
	"unsafe"

	"github.com/e-nikolov/scratch/pkg/runtime"
)

func main() {
	open("./build/c-shared-plug.so")
}

// Plugin is a loaded Go plugin.
type Plugin struct {
	pluginpath string
	err        string        // set if plugin failed to load
	loaded     chan struct{} // closed when loaded
	syms       map[string]any
}

var (
	pluginsMu sync.Mutex
	plugins   map[string]*Plugin
)

func open(name string) (*Plugin, error) {
	cPath := make([]byte, C.PATH_MAX+1)
	cRelName := make([]byte, len(name)+1)
	copy(cRelName, name)
	if C.realpath(
		(*C.char)(unsafe.Pointer(&cRelName[0])),
		(*C.char)(unsafe.Pointer(&cPath[0]))) == nil {
		return nil, errors.New(`plugin.Open("` + name + `"): realpath failed`)
	}

	filepath := C.GoString((*C.char)(unsafe.Pointer(&cPath[0])))

	pluginsMu.Lock()
	if p := plugins[filepath]; p != nil {
		pluginsMu.Unlock()
		if p.err != "" {
			return nil, errors.New(`plugin.Open("` + name + `"): ` + p.err + ` (previous failure)`)
		}
		<-p.loaded
		return p, nil
	}
	var cErr *C.char
	h := C.pluginOpen((*C.char)(unsafe.Pointer(&cPath[0])), &cErr)
	if h == 0 {
		pluginsMu.Unlock()
		return nil, errors.New(`plugin.Open("` + name + `"): ` + C.GoString(cErr))
	}
	if len(name) > 3 && name[len(name)-3:] == ".so" {
		name = name[:len(name)-3]
	}
	if plugins == nil {
		plugins = make(map[string]*Plugin)
	}
	fmt.Println(h)

	//------------------------------------------------------------------------------------//
	//------------------------------------------------------------------------------------//

	pluginpath, syms, initTasks, errstr := runtime.PluginLastModuleInit()
	_, _, _, _ = pluginpath, syms, initTasks, errstr

	// if errstr != "" {
	// 	plugins[filepath] = &Plugin{
	// 		pluginpath: pluginpath,
	// 		err:        errstr,
	// 	}
	// 	pluginsMu.Unlock()
	// 	return nil, errors.New(`plugin.Open("` + name + `"): ` + errstr)
	// }
	// // This function can be called from the init function of a plugin.
	// // Drop a placeholder in the map so subsequent opens can wait on it.
	// p := &Plugin{
	// 	pluginpath: pluginpath,
	// 	loaded:     make(chan struct{}),
	// }
	// plugins[filepath] = p
	// pluginsMu.Unlock()

	// runtime.DoInit(initTasks)

	// // Fill out the value of each plugin symbol.
	// updatedSyms := map[string]any{}
	// for symName, sym := range syms {
	// 	isFunc := symName[0] == '.'
	// 	if isFunc {
	// 		delete(syms, symName)
	// 		symName = symName[1:]
	// 	}

	// 	fullName := pluginpath + "." + symName
	// 	cname := make([]byte, len(fullName)+1)
	// 	copy(cname, fullName)

	// 	p := C.pluginLookup(h, (*C.char)(unsafe.Pointer(&cname[0])), &cErr)
	// 	if p == nil {
	// 		return nil, errors.New(`plugin.Open("` + name + `"): could not find symbol ` + symName + `: ` + C.GoString(cErr))
	// 	}
	// 	valp := (*[2]unsafe.Pointer)(unsafe.Pointer(&sym))
	// 	if isFunc {
	// 		(*valp)[1] = unsafe.Pointer(&p)
	// 	} else {
	// 		(*valp)[1] = p
	// 	}
	// 	// we can't add to syms during iteration as we'll end up processing
	// 	// some symbols twice with the inability to tell if the symbol is a function
	// 	updatedSyms[symName] = sym
	// }
	// p.syms = updatedSyms

	// close(p.loaded)
	// return p, nil
	// spew.Dump(h)

	return nil, nil
}
