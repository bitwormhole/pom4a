package pomtool

import (
	"fmt"
	"os"
)

type pom4aVersion struct {
}

func (inst *pom4aVersion) exefile() string {
	args := os.Args
	for index, ar := range args {
		if index == 0 {
			return ar
		}
	}
	return "unknown"
}

func (inst *pom4aVersion) version() string {
	return "v0.0.0"
}

func (inst *pom4aVersion) run() error {

	list := []string{}

	list = append(list, "pom4a (POM Tool for Android)")
	list = append(list, "     version : "+inst.version())
	list = append(list, "  executable : "+inst.exefile())

	for _, row := range list {
		fmt.Println(row)
	}
	return nil
}
