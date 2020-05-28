package main

import (
	"fmt"
	"net/http"

	"github.com/urfave/cli/v2"
	"github.com/vaulty/vaulty/api"
	"github.com/vaulty/vaulty/storage/inmem"
)

var apiCommand = &cli.Command{
	Name:  "api",
	Usage: "run REST api server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Value:   "vaulty.yml",
			Usage:   "Vaulty configuration file",
		},
		&cli.StringFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Value:   "3000",
		},
	},
	Action: func(c *cli.Context) error {
		port := c.String("port")

		storage := inmem.NewStorage()
		server := api.NewServer(storage)

		fmt.Printf("==> Vaulty API server started on port %v!", port)
		err := http.ListenAndServe(":"+port, server)
		return err
	},
}
