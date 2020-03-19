package main

import (
	"flag"

	"github.com/vaulty/proxy/core"
	"github.com/vaulty/proxy/proxy"
)

func main() {
	env := flag.String("e", "development", "proxy environment")
	flag.Parse()
	core.LoadConfig(*env)

	core.InitRedisClient(core.Config())

	proxy := proxy.NewProxy()
	proxy.Run()
}
