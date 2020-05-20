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
