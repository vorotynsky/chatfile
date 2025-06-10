package main

import "github.com/alexflint/go-arg"

func main() {
	var args struct {
		Run *RunCmd `arg:"subcommand:run" help:"Run a chatfile"`
	}
	arg.MustParse(&args)

	if args.Run != nil {
		args.Run.Execute()
	}
}
