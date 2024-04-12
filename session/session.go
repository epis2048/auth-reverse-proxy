package session

import (
	"auth-reverse-proxy/config"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/sessions"
	"github.com/hertz-contrib/sessions/cookie"
)

func Init(proxyConfig config.ProxyConfig, h *server.Hertz) {
	h.Use(sessions.New(proxyConfig.Session.Name, cookie.NewStore([]byte(proxyConfig.Session.Secret))))
}

func Set(c *app.RequestContext, key string, value string) {
	session := sessions.Default(c)
	session.Set(key, value)
	session.Save()
}

func Get(c *app.RequestContext, key string) string {
	session := sessions.Default(c)
	v := session.Get(key)
	if v != nil {
		return v.(string)
	}
	return ""
}
