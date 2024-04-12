package auth

import (
	"auth-reverse-proxy/config"
	"auth-reverse-proxy/session"
	"auth-reverse-proxy/utils"
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
)

type Token struct {
	ProxyConfig config.ProxyConfig
}

func NewToken(proxyConfig config.ProxyConfig) Auth {
	var token Token
	token.ProxyConfig = proxyConfig
	return token
}

func (t Token) GetConfig() config.ProxyConfig {
	return t.ProxyConfig
}

func (t Token) RegisterRouter(g *route.RouterGroup) {}

func (t Token) HandlerLogin() app.HandlerFunc {
	return func(context context.Context, c *app.RequestContext) {
		var uid string
		var err error
		resp := utils.RequestGet(strings.ReplaceAll(t.ProxyConfig.Token.Valid.URL, "{token}", c.Query("token")))

		switch t.ProxyConfig.Token.Valid.Format {
		case "json":
			uid, err = utils.ParseUIDFromJson(resp.Body(), t.ProxyConfig.Token.Valid.JSONPath)
			if err != nil {
				c.AbortWithMsg(err.Error(), consts.StatusOK)
			}
		case "xml":
			uid, err = utils.ParseUIDFromXml(resp.Body(), t.ProxyConfig.Token.Valid.XMLPath)
			if err != nil {
				c.AbortWithMsg(err.Error(), consts.StatusOK)
			}
		default:
			c.AbortWithMsg("unknown token valid format", consts.StatusOK)
		}

		session.Set(c, "uid", uid)
		if !c.IsAborted() {
			if c.Query("callback") != "" {
				c.Redirect(302, []byte(c.Query("callback")))
				c.Abort()
			} else {
				c.Redirect(302, []byte("/"))
				c.Abort()
			}
		}
	}
}

func (t Token) HandlerLogout() app.HandlerFunc {
	return DefaultHandlerLogout(t)
}

func (t Token) HandlerLoginStatus() app.HandlerFunc {
	return DefaultHandlerLoginStatus(t)
}

func (t Token) UnAuthed() app.HandlerFunc {
	return DefaultUnAuthed(t, "unauthed")
}

func (t Token) MiddlewareAuth() app.HandlerFunc {
	return DefaultSessionMiddleWareAuth(t)
}
