package main

import (
	"os"

	"github.com/bitwormhole/pom4a/pomtool"
)

func main() {
	args := os.Args
	err := pomtool.Run(args)
	if err != nil {
		panic(err)
	}
}
