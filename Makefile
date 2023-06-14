go-stdlib-globals:
	go build -buildmode plugin -o build/stdlib-globals-plug.so ./plugins/stdlib-globals-plug/main.go
	go run ./cmd/go-stdlib-globals

go-inline-c:
	go run ./cmd/go-inline-c

go-extern-c:
	go run ./cmd/go-extern-c

c-dlink:
	go build -o build/c-shared-plug.so -buildmode=c-shared ./plugins/c-shared-plug/main.go
	gcc -o build/c-dlink ./cmd/c-dlink/main.c build/c-shared-plug.so
	./build/c-dlink


c-dlopen:
	go build -o build/c-shared-plug.so -buildmode=c-shared ./plugins/c-shared-plug/main.go
	gcc -o build/c-dlopen ./cmd/c-dlopen/main.c -ldl
	./build/c-dlopen


py-dlopen:
	go build -o build/c-shared-plug.so -buildmode=c-shared ./plugins/c-shared-plug/main.go
	python ./cmd/py-dlopen/main.py
	


py-cffi:
	go build -o build/c-shared-plug.so -buildmode=c-shared ./plugins/c-shared-plug/main.go
	python ./cmd/py-cffi/main.py
	


go-dlopen:
	go build -o build/c-shared-plug.so -buildmode=c-shared ./plugins/c-shared-plug/main.go
	go run ./cmd/go-dlopen/
	

