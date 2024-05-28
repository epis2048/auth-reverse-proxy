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

func StartServer(proxyConfig config.ProxyConfig) {
	h := server.New(server.WithHostPorts(proxyConfig.ListenAddress))
	var authMiddleWare auth.Auth
	switch proxyConfig.Auth {
	case "none":
		authMiddleWare = auth.NewNone(proxyConfig)
	case "token":
		authMiddleWare = auth.NewToken(proxyConfig)
	case "jwt":
		authMiddleWare = auth.NewJwt(proxyConfig)
	case "cas":
		authMiddleWare = auth.NewCas(proxyConfig)
	default:
		log.Fatalln("unknown auth type")
	}

	session.Init(proxyConfig, h)

	g := h.Group("/__auth")
	g.GET("/login", authMiddleWare.HandlerLogin())
	g.GET("/login/status", authMiddleWare.HandlerLoginStatus())
	g.GET("/logout", authMiddleWare.HandlerLogout())
	authMiddleWare.RegisterRouter(g)

	for i, r := range proxyConfig.Reverse {
		switch r.Type {
		case "http_jump":
			for _, u := range r.URL {
				h.Any(u, rewrite.Jump(r.Code, r.To))
			}
		case "http":
			proxy, err := reverseproxy.NewSingleHostReverseProxy(r.Backend)
			if err != nil {
				log.Fatalln(err)
			}
			rewrite.Rewrite(proxy, proxyConfig, i)
			for _, u := range r.URL {
				h.Any(u, authMiddleWare.MiddlewareAuth(), proxy.ServeHTTP)
			}
		case "websocket":
			proxy := reverseproxy.NewWSReverseProxy(r.Backend)
			for _, u := range r.URL {
				h.GET(u, authMiddleWare.MiddlewareAuth(), proxy.ServeHTTP)
			}
		}
	}
	h.Run()
}
