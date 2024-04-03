package session

import (
	"auth-reverse-proxy/config"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/sessions"
	"github.com/hertz-contrib/sessions/cookie"
)

func Init(h *server.Hertz) {
	h.Use(sessions.New(config.Config.Proxy.Session.Name, cookie.NewStore([]byte(config.Config.Proxy.Session.Secret))))
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
