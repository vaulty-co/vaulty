package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/vaulty/vaulty"
	"github.com/vaulty/vaulty/config"
)

var conf = config.NewConfig()

var proxyCommand = &cli.Command{
	Name:  "proxy",
	Usage: "run proxy server",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:        "debug",
			Usage:       "enable debug (exposes request and response bodies)",
			Destination: &conf.Debug,
		},
		&cli.StringFlag{
			Name:        "address",
			Aliases:     []string{"a"},
			Value:       ":8080",
			Usage:       "address that vaulty should listen on",
			Destination: &conf.Address,
		},
		&cli.StringFlag{
			Name:        "routes-file",
			Aliases:     []string{"r"},
			Value:       "./routes.json",
			Usage:       "routes file",
			Destination: &conf.RoutesFile,
		},
		&cli.StringFlag{
			Name:        "ca-path",
			Aliases:     []string{"ca"},
			Value:       "./",
			Usage:       "path to CA key and cert",
			Destination: &conf.CAPath,
		},
		&cli.StringFlag{
			Name:        "proxy-pass",
			Aliases:     []string{"p"},
			Usage:       "forward proxy password",
			EnvVars:     []string{"PROXY_PASS"},
			Destination: &conf.ProxyPassword,
		},
		&cli.StringFlag{
			Name:        "key",
			Aliases:     []string{"k"},
			Usage:       "forward proxy password",
			Destination: &conf.Encryption.Key,
		},
		&cli.StringFlag{
			Name:        "hash-salt",
			Usage:       "salt for the hash action",
			EnvVars:     []string{"HASH_SALT"},
			Destination: &conf.Salt,
		},
	},
	Action: func(c *cli.Context) error {
		if err := conf.FromEnvironment(); err != nil {
			return err
		}

		if err := conf.GenerateMissedValues(); err != nil {
			return fmt.Errorf("Error with generating missed values: %s", err)
		}

		return vaulty.Run(conf)
	},
}
