package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var (
	Version string
	Build   string
)

var versionCommand = &cli.Command{
	Name:  "version",
	Usage: "get version of Vaulty",
	Action: func(c *cli.Context) error {
		fmt.Println("Version: ", Version)
		fmt.Println("Build Time: ", Build)
		return nil
	},
}
