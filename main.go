package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/guessi/cloudtrail-cli/cmd"
	"github.com/guessi/cloudtrail-cli/pkg/constants"
	"github.com/urfave/cli/v2"
)

func showVersion() {
	r, _ := regexp.Compile(`v[0-9]\.[0-9]+\.[0-9]+`)
	versionInfo := r.FindString(constants.GitVersion)
	fmt.Println("cloudtrail-cli", versionInfo)
	fmt.Println(" Git Commit:", constants.GitVersion)
	fmt.Println(" Build with:", constants.GoVersion)
	fmt.Println(" Build time:", constants.BuildTime)
}

func main() {
	app := &cli.App{
		Name:    constants.NAME,
		Usage:   constants.USAGE,
		Version: constants.GitVersion,
		Flags:   cmd.Flags,
		Action: func(c *cli.Context) error {
			cmd.Wrapper(c)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print version number",
				Action: func(cCtx *cli.Context) error {
					showVersion()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
