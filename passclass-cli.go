package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/eczarny/passclass-cli/command"
	"github.com/wsxiaoys/terminal"
)

func authToken() string {
	return fmt.Sprintf("Token %s", os.Getenv("CLASSPASS_TOKEN"))
}

func run(cmd string, args []string) {
	stderr := terminal.Stderr
	switch cmd {
	case "reserve":
		command.Reserve(cmd, args, authToken)
	case "schedule":
		command.Schedule(cmd, args, authToken)
	case "venues":
		command.Venues(cmd, args, authToken)
	default:
		stderr.Color("r").Print(cmd).Print(" is not a passclass-cli command. See 'passclass-cli -h' for help.").Reset().Nl()
	}
}

func main() {
	usage := `passclass-cli - Reserve your favorite classes with ClassPass.

Usage: passclass-cli [--help] [--version] <cmd> [<args>...]

Options:
  --help     Show help
  --version  Show version

Commands:
  reserve    Reserve a class with ClassPass
  schedule   List a class schedule
  venues     List all ClassPass venues
`

	args, _ := docopt.Parse(usage, nil, true, "passclass-cli 0.0.1", false)

	run(args["<cmd>"].(string), args["<args>"].([]string))
}
