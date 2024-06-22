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
	fmt.Println("listen and serve on 127.0.0.1:49557")
	fmt.Println(http.ListenAndServe(":49557", proxy))
}
