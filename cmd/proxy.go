package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	"github.com/urfave/cli/v2"
	vaulty "github.com/vaulty/vaulty"
)

var config = vaulty.NewConfig()

var proxyCommand = &cli.Command{
	Name:  "proxy",
	Usage: "run proxy server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "address, a",
			Value:       ":8080",
			Usage:       "address that vaulty should listen on",
			Destination: &config.Address,
		},
		&cli.StringFlag{
			Name:        "routes-file, r",
			Value:       "./routes.json",
			Usage:       "routes file",
			Destination: &config.RoutesFile,
		},
		&cli.StringFlag{
			Name:        "ca-path, ca",
			Value:       "./",
			Usage:       "path to CA key and cert",
			Destination: &config.CAPath,
		},
		&cli.StringFlag{
			Name:        "proxy-pass, p",
			Usage:       "forward proxy password",
			EnvVars:     []string{"PROXY_PASS"},
			Destination: &config.ProxyPassword,
		},
		&cli.StringFlag{
			Name:        "key, k",
			Usage:       "forward proxy password",
			EnvVars:     []string{"ENCRYPTION_KEY"},
			Destination: &config.EncryptionKey,
		},
	},
	Action: func(c *cli.Context) error {
		log.SetFormatter(&prefixed.TextFormatter{
			FullTimestamp: true,
		})

		if err := config.GenerateMissedValues(); err != nil {
			return fmt.Errorf("Error with generating missed values: %s", err)
		}

		return vaulty.Run(config)
	},
}
