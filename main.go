package main

import (
	"auth-reverse-proxy/config"
	"flag"
)

var configFile string

func parseFlags() {
	flag.StringVar(&configFile, "config", "config.yaml", "config file name")
	flag.Parse()
}

func main() {
	parseFlags()
	config.Init(configFile)
	StartServer()
}
