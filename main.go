package main

import (
	"flag"
	"fmt"

	"github.com/vaulty/proxy/core"
	"github.com/vaulty/proxy/proxy"
)

func main() {
	env := flag.String("e", "development", "proxy environment")
	flag.Parse()
	core.LoadConfig(fmt.Sprintf("config/%s.yml", *env))

	proxy := proxy.NewProxy()
	proxy.Run()
}
