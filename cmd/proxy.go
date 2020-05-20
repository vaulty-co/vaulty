package main

import (
	"fmt"
	"net/http"

	"github.com/urfave/cli/v2"
	"github.com/vaulty/proxy/core"
	"github.com/vaulty/proxy/encrypt"
	"github.com/vaulty/proxy/proxy"
	"github.com/vaulty/proxy/secrets"
	"github.com/vaulty/proxy/storage"
	"github.com/vaulty/proxy/storage/inmem"
	"github.com/vaulty/proxy/transform/action"
)

var proxyCommand = &cli.Command{
	Name:  "proxy",
	Usage: "run proxy server",
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
			Value:   "8080",
		},
		&cli.StringFlag{
			Name: "routes-file",
		},
	},
	Action: func(c *cli.Context) error {
		port := c.String("port")
		configFile := c.String("config")
		config := core.LoadConfig(configFile)

		if c.IsSet("routes-file") {
			config.RoutesFile = c.String("routes-file")
		}

		st := inmem.NewStorage()

		encrypter, err := encrypt.NewEncrypter(config.EncryptionKey)
		if err != nil {
			return err
		}

		secretStorage := secrets.NewEphemeralStorage(encrypter)

		loaderOptions := &storage.LoaderOptions{
			ActionOptions: &action.Options{
				Encrypter:     encrypter,
				SecretStorage: secretStorage,
			},
			Storage: st,
		}

		err = storage.LoadFromFile(config.RoutesFile, loaderOptions)
		if err != nil {
			return err
		}

		proxy, err := proxy.NewProxy(st, config)
		if err != nil {
			return err
		}

		fmt.Printf("==> Vaulty proxy server started on port %v!\n", port)
		err = http.ListenAndServe(":"+port, proxy)
		return err
	},
}
