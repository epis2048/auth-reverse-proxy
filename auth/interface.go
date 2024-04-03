package auth

import (
	"auth-reverse-proxy/config"
	"auth-reverse-proxy/session"
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
)

type Auth interface {
	Init()

	// handlers
	RegisterRouter(g *route.RouterGroup)
	HandlerLogin() app.HandlerFunc
	HandlerLogout() app.HandlerFunc
	HandlerLoginStatus() app.HandlerFunc

	// public functions
	UnAuthed() app.HandlerFunc
	MiddlewareAuth() app.HandlerFunc
}

func DefaultHandlerLoginStatus(a Auth) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.AbortWithMsg(session.Get(c, "uid"), consts.StatusOK)
	}
}

func DefaultHandlerLogout(a Auth) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		session.Set(c, "uid", "")
	}
}

func DefaultSessionMiddleWareAuth(a Auth) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		uid := session.Get(c, "uid")
		if uid == "" {
			a.UnAuthed()(ctx, c)
		}
	}
}

func DefaultUnAuthed(a Auth, msg string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		switch config.Config.Proxy.UnAuthedResponse {
		case "text":
			c.AbortWithMsg(msg, consts.StatusOK)
		case "jump":
			c.Redirect(302, []byte(strings.ReplaceAll(msg, "{callback}", string(c.URI().Scheme())+"://"+string(c.URI().Host())+"/__auth/login")))
			c.Abort()
		}
	}
}
