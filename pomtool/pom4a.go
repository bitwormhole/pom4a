package pomtool

import "fmt"

type pom4a struct {
	context *Context
}

func (inst *pom4a) install() error {
	i := &pom4aInstall{}
	return i.run()
}

func (inst *pom4a) uninstall() error {
	i := &pom4aUninstall{}
	return i.run()
}

func (inst *pom4a) help() error {
	i := &pom4aHelp{}
	return i.run()
}

func (inst *pom4a) version() error {
	i := &pom4aVersion{
		context: inst.context,
	}
	return i.run()
}

func (inst *pom4a) getCommand(args []string) string {
	for index, ar := range args {
		if index == 1 {
			return ar
		}
	}
	return ""
}

func (inst *pom4a) run(args []string) error {

	cmd := inst.getCommand(args)
	table := make(map[string]func() error)

	table["help"] = inst.help
	table["install"] = inst.install
	table["uninstall"] = inst.uninstall
	table["version"] = inst.version

	fn := table[cmd]
	if fn == nil {
		const msg = "no command named '%s', use 'pom4a help' for more info"
		return fmt.Errorf(msg, cmd)
	}
	return fn()
}

////////////////////////////////////////////////////////////////////////////////

func Run(ctx *Context, args []string) error {
	tool := &pom4a{
		context: ctx,
	}
	return tool.run(args)
}
