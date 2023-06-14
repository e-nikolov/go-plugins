package main

import (
	"fmt"
	"plugin"

	"github.com/e-nikolov/scratch/pkg/shared"
)

func main() {
	fmt.Printf("STDLIB Plugin - Global Variable Demo\n\n")

	fmt.Printf("Initial value: \n")
	fmt.Printf("GlobalVariable = %q\n\n", shared.GlobalVariable)

	_, err := plugin.Open("build/stdlib-globals-plug.so")

	if err != nil {
		panic(err)
	}

	fmt.Printf("Modified by a plugin: \n")
	fmt.Printf("GlobalVariable: %q\n", shared.GlobalVariable)
}
