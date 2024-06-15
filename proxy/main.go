package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	proxy.ConnectDial = func(network, addr string) (net.Conn, error) {
		fmt.Println(network, addr)
		if strings.Contains(addr, "127.0.0.1") ||
			strings.Contains(addr, "localhost") {
			return net.Dial(network, addr)
		}
		https_proxy := "http://172.16.100.150:23128"
		// if strings.Contains(addr, "openai.com") ||
		// 	strings.Contains(addr, "chatgpt.com") ||
		// 	strings.Contains(addr, "anthropic.com") ||
		// 	strings.Contains(addr, "claude.ai") {
		// 	https_proxy = "http://172.16.0.150:23128"
		// }
		return proxy.NewConnectDialToProxy(https_proxy)(network, addr)
	}
	proxy.OnResponse(goproxy.RespConditionFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) bool {
		return true
	})).Do(goproxy.FuncRespHandler(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		fmt.Println("OnResponse:", resp.Status, ctx.Error)
		return resp
	}))
	proxy.OnRequest(goproxy.ReqConditionFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		return true
	})).Do(goproxy.FuncReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		fmt.Println("OnRequest:", req.URL.String(), ctx.Resp)
		return req, ctx.Resp
	}))
	fmt.Println(http.ListenAndServe(":49557", proxy))
}
