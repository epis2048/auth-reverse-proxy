package auth

import (
	"auth-reverse-proxy/config"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"
)

type None struct {
	ProxyConfig config.ProxyConfig
}

func NewNone(proxyConfig config.ProxyConfig) Auth {
	var token None
	token.ProxyConfig = proxyConfig
	return token
}

func (t None) GetConfig() config.ProxyConfig {
	return t.ProxyConfig
}

func (t None) RegisterRouter(g *route.RouterGroup) {}

func (t None) HandlerLogin() app.HandlerFunc {
	return func(context context.Context, c *app.RequestContext) {
		if c.Query("callback") != "" {
			c.Redirect(302, []byte(c.Query("callback")))
			c.Abort()
		} else {
			c.Redirect(302, []byte("/"))
			c.Abort()
		}
	}
}

func (t None) HandlerLogout() app.HandlerFunc {
	return DefaultHandlerLogout(t)
}

func (t None) HandlerLoginStatus() app.HandlerFunc {
	return DefaultHandlerLoginStatus(t)
}

func (t None) UnAuthed() app.HandlerFunc {
	return DefaultUnAuthed(t, "unauthed")
}

func (t None) MiddlewareAuth() app.HandlerFunc {
	return func(context context.Context, c *app.RequestContext) {}
}
