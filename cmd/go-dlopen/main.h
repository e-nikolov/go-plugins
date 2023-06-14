#ifndef _MAIN_H
#define _MAIN_H

// define types needed
typedef long long go_int;
typedef double go_float64;
typedef struct{void *arr; go_int len; go_int cap;} go_slice;
typedef struct{const char *p; go_int len;} go_str;
// typedef struct { void *t; void *v; } go_error;
// typedef struct { string; } go_error;

extern int LoadPlugin(go_str path);

extern go_int Add(go_int a, go_int b);
extern go_float64 Cosine(go_float64 x);
extern go_int Log(go_str msg);
extern void Sort(go_slice vals);
extern char* StartHTTPServer(char* addr, char* message);

#endif