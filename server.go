package main

import (
	"auth-reverse-proxy/auth"
	"auth-reverse-proxy/config"
	"auth-reverse-proxy/proxy/rewrite"
	"auth-reverse-proxy/session"
	"log"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/reverseproxy"
)

func StartServer() {
	h := server.New(server.WithHostPorts(config.Config.ListenAddress))
	proxy, err := reverseproxy.NewSingleHostReverseProxy(config.Config.Proxy.Reverse.Backend)
	if err != nil {
		log.Fatalln(err)
	}

	var authMiddleWare auth.Auth
	switch config.Config.Proxy.Auth {
	case "token":
		authMiddleWare = auth.Token{}
	case "cas":
		authMiddleWare = auth.Cas{}
	default:
		log.Fatalln("unknown auth type")
	}

	session.Init(h)
	authMiddleWare.Init()

	g := h.Group("/__auth")
	g.GET("/login", authMiddleWare.HandlerLogin())
	g.GET("/login/status", authMiddleWare.HandlerLoginStatus())
	g.GET("/logout", authMiddleWare.HandlerLogout())
	authMiddleWare.RegisterRouter(g)

	rewrite.RewriteHeader(proxy)
	for _, u := range config.Config.Proxy.Reverse.URL {
		h.Any(u, authMiddleWare.MiddlewareAuth(), proxy.ServeHTTP)
	}

	h.Spin()
}
