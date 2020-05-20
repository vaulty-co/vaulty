package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:                 "vaulty",
		Usage:                "Vaulty command line utility",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "vaulty.yml",
				Usage:   "Vaulty configuration file",
			},
		},
	}

	app.Commands = []*cli.Command{
		apiCommand,
		proxyCommand,
		versionCommand,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
