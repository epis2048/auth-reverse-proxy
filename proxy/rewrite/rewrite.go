package rewrite

import (
	"auth-reverse-proxy/config"
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/hertz-contrib/reverseproxy"
)

func Jump(code int, to string) app.HandlerFunc {
	return func(context context.Context, c *app.RequestContext) {
		c.Redirect(code, []byte(to))
		c.Abort()
	}
}

func Rewrite(proxy *reverseproxy.ReverseProxy, proxyConfig config.ProxyConfig, reverseIndex int) {
	rewriteHeader(proxy, proxyConfig, reverseIndex)
}

func rewriteHeader(proxy *reverseproxy.ReverseProxy, proxyConfig config.ProxyConfig, reverseIndex int) {
	proxy.SetModifyResponse(func(resp *protocol.Response) error {
		for _, v := range proxyConfig.Reverse[reverseIndex].Rewrite.Header {
			header := resp.Header.Get(v.Name)
			if header != "" {
				header = strings.ReplaceAll(header, v.From, v.To)
				resp.Header.Set(v.Name, header)
			}
		}
		return nil
	})
}
