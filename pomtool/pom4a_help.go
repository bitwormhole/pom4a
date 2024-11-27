package pomtool

import "fmt"

type pom4aHelp struct{}

func (inst *pom4aHelp) run() error {

	list := []string{}

	list = append(list, "Usage:")
	list = append(list, "    pom4a [cmd]")
	list = append(list, "cmd: [help|install|uninstall|version]")

	for _, row := range list {
		fmt.Println(row)
	}
	return nil
}
