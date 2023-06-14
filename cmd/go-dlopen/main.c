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
