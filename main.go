package main

import (
	"flag"
	"fmt"

	"github.com/vaulty/proxy/core"
	"github.com/vaulty/proxy/proxy"
	"github.com/vaulty/proxy/storage"
	"github.com/vaulty/proxy/transformer"
)

func main() {
	env := flag.String("e", "development", "proxy environment")
	flag.Parse()
	config := core.LoadConfig(fmt.Sprintf("config/%s.yml", *env))

	redisClient := core.NewRedisClient(config)
	storage := storage.NewStorage(redisClient)
	transformer := transformer.NewTransformer(redisClient)

	proxy := proxy.NewProxy(storage, transformer, config)
	proxy.Run()
}
