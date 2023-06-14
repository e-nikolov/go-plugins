package main

import "github.com/e-nikolov/scratch/pkg/shared"

func init() {
	shared.GlobalVariable = "Value Modified By a Plugin"
}

func main() {}
