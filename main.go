package main

import (
	"os"

	"github.com/bitwormhole/pom4a/pomtool"
)

func main() {
	args := os.Args
	ctx := &pomtool.Context{
		AppName:     theModuleName,
		AppTitle:    theModuleTitle,
		AppVersion:  theModuleVersion,
		AppRevision: theModuleRevision,
	}
	err := pomtool.Run(ctx, args)
	if err != nil {
		panic(err)
	}
}
