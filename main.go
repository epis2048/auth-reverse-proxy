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
	for _, proxyConfig := range config.Config {
		go func(proxyConfig config.ProxyConfig) {
			StartServer(proxyConfig)
		}(proxyConfig)
	}
	ch := make(chan string)
	<-ch
}
