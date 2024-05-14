package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/elazarl/goproxy"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	proxy.ConnectDial = func(network, addr string) (net.Conn, error) {
		return proxy.NewConnectDialToProxy("http://172.16.100.150:23128")(network, addr)
	}
	fmt.Println(http.ListenAndServe(":49557", proxy))
}
