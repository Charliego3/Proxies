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
		https_proxy := "http://172.16.100.150:23128"
		if strings.Contains(addr, "openai.com") ||
			strings.Contains(addr, "chatgpt.com") {
			https_proxy = "http://172.16.0.150:23128"
		}
		return proxy.NewConnectDialToProxy(https_proxy)(network, addr)
	}
	fmt.Println(http.ListenAndServe(":49557", proxy))
}
