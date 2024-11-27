package pomtool

import (
	"fmt"
	"os"
	"strconv"
)

type pom4aVersion struct {
	context *Context
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
	return inst.context.AppVersion
}

func (inst *pom4aVersion) rev() string {
	n := inst.context.AppRevision
	return strconv.Itoa(n)
}

func (inst *pom4aVersion) run() error {

	list := []string{}
	name := inst.context.AppName
	title := inst.context.AppTitle

	list = append(list, name+" ("+title+")")
	list = append(list, "     version : "+inst.version())
	list = append(list, "    revision : "+inst.rev())
	list = append(list, "  executable : "+inst.exefile())

	for _, row := range list {
		fmt.Println(row)
	}
	return nil
}
