package main

import (
	"os"

	"github.com/guessi/cloudtrail-cli/cmd"
	"github.com/guessi/cloudtrail-cli/pkg/constants"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    constants.NAME,
		Usage:   constants.USAGE,
		Version: constants.VERSION,
		Flags:   cmd.Flags,
		Action: func(c *cli.Context) error {
			cmd.Wrapper(c)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}
