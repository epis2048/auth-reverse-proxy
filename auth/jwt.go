package auth

import (
	"auth-reverse-proxy/config"
	"auth-reverse-proxy/session"
	"auth-reverse-proxy/utils"
	"context"
	"encoding/base64"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/golang-jwt/jwt/v5"
)

type Jwt struct {
	ProxyConfig config.ProxyConfig
}

func NewJwt(proxyConfig config.ProxyConfig) Auth {
	var jwt Jwt
	jwt.ProxyConfig = proxyConfig
	return jwt
}

func (t Jwt) GetConfig() config.ProxyConfig {
	return t.ProxyConfig
}

func (t Jwt) RegisterRouter(g *route.RouterGroup) {}

func (t Jwt) HandlerLogin() app.HandlerFunc {
	return func(context context.Context, c *app.RequestContext) {
		var uid string
		var err error
		tokenString := c.Query("token")
		_, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(t.ProxyConfig.Jwt.Valid.Secret), nil
		})
		if err != nil {
			c.AbortWithMsg("JWT is invalid: "+err.Error(), consts.StatusOK)
			return
		}
		parts := strings.Split(tokenString, ".")
		decodedPayload, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			c.AbortWithMsg("Failed to decode payload segment: "+err.Error(), consts.StatusOK)
			return
		}
		uid, err = utils.ParseUIDFromJson(decodedPayload, t.ProxyConfig.Jwt.Valid.JSONPath)
		if err != nil {
			c.AbortWithMsg("Parse UID From JWT err: "+err.Error(), consts.StatusOK)
			return
		}
		session.Set(c, "uid", uid)
		if c.Query("callback") != "" {
			c.Redirect(302, []byte(c.Query("callback")))
			c.Abort()
		} else {
			c.Redirect(302, []byte("/"))
			c.Abort()
		}
	}
}

func (t Jwt) HandlerLogout() app.HandlerFunc {
	return DefaultHandlerLogout(t)
}

func (t Jwt) HandlerLoginStatus() app.HandlerFunc {
	return DefaultHandlerLoginStatus(t)
}

func (t Jwt) UnAuthed() app.HandlerFunc {
	return DefaultUnAuthed(t, "unauthed")
}

func (t Jwt) MiddlewareAuth() app.HandlerFunc {
	return DefaultSessionMiddleWareAuth(t)
}
