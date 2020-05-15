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
			Name:    "port",
			Aliases: []string{"p"},
			Value:   "8080",
		},
	},
	Action: func(c *cli.Context) error {
		port := c.String("port")
		environment := c.String("environment")
		config := core.LoadConfig(fmt.Sprintf("config/%s.yml", environment))
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

		fmt.Printf("==> Vaulty proxy server started on port %v! in %v environment\n", port, environment)
		err = http.ListenAndServe(":"+port, proxy)
		return err
	},
}
