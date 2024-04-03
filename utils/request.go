package utils

import (
	"context"
	"crypto/tls"

	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/network/standard"
	"github.com/cloudwego/hertz/pkg/protocol"
)

func RequestGet(url string) *protocol.Response {
	client, _ := client.NewClient(client.WithTLSConfig(&tls.Config{
		InsecureSkipVerify: true,
	}),
		client.WithDialer(standard.NewDialer()))
	req, resp := protocol.AcquireRequest(), protocol.AcquireResponse()
	req.SetRequestURI(url)
	req.SetMethod("GET")
	client.Do(context.Background(), req, resp)
	return resp
}
