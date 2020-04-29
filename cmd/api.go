package main

import (
	"fmt"
	"net/http"

	"github.com/urfave/cli/v2"
	"github.com/vaulty/proxy/api"
	"github.com/vaulty/proxy/core"
	"github.com/vaulty/proxy/storage"
)

var apiCommand = &cli.Command{
	Name:  "api",
	Usage: "run REST api server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Value:   "3000",
		},
	},
	Action: func(c *cli.Context) error {
		port := c.String("port")
		environment := c.String("environment")
		config := core.LoadConfig(fmt.Sprintf("config/%s.yml", environment))
		redisClient := core.NewRedisClient(config)
		storage := storage.NewRedisStorage(redisClient)

		server := api.NewServer(storage)

		fmt.Printf("==> Vaulty API server started on port %v! in %v environment\n", port, environment)
		err := http.ListenAndServe(":"+port, server)
		return err
	},
}
