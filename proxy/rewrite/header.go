package rewrite

import (
	"auth-reverse-proxy/config"
	"strings"

	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/hertz-contrib/reverseproxy"
)

func RewriteHeader(proxy *reverseproxy.ReverseProxy, proxyConfig config.ProxyConfig, reverseIndex int) {
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
