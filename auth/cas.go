package auth

import (
	"auth-reverse-proxy/config"
	"auth-reverse-proxy/session"
	"auth-reverse-proxy/utils"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
)

type Cas struct {
	ProxyConfig config.ProxyConfig
}

func NewCas(proxyConfig config.ProxyConfig) Auth {
	var cas Cas
	cas.ProxyConfig = proxyConfig
	return cas
}

func (t Cas) GetConfig() config.ProxyConfig {
	return t.ProxyConfig
}

func (t Cas) RegisterRouter(g *route.RouterGroup) {}

func (t Cas) HandlerLogin() app.HandlerFunc {
	return func(context context.Context, c *app.RequestContext) {
		var uid string
		var err error
		resp := utils.RequestGet(t.ProxyConfig.Cas.EndPoint + "/serviceValidate?service=" + string(c.URI().Scheme()) + "://" + string(c.URI().Host()) + "/__auth/login" + "&ticket=" + c.Query("ticket"))

		uid, err = utils.ParseUIDFromXml(resp.Body(), t.ProxyConfig.Cas.XMLPath)
		if err != nil {
			c.AbortWithMsg(err.Error(), consts.StatusOK)
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

func (t Cas) HandlerLogout() app.HandlerFunc {
	return DefaultHandlerLogout(t)
}

func (t Cas) HandlerLoginStatus() app.HandlerFunc {
	return DefaultHandlerLoginStatus(t)
}

func (t Cas) UnAuthed() app.HandlerFunc {
	return DefaultUnAuthed(t, t.ProxyConfig.Cas.EndPoint+"/login?service={callback}")
}

func (t Cas) MiddlewareAuth() app.HandlerFunc {
	return DefaultSessionMiddleWareAuth(t)
}
