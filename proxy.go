package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	rules   = new(RuleDatasource)
	proxies = new(ProxyDatasource)
	server  = NewProxyServer()
)

func init() {
	proxies.mux = new(sync.RWMutex)
	proxyvalue := defaults.StringForKey(proxyDefaultsKey)
	if proxyvalue != "" {
		json.Unmarshal([]byte(proxyvalue), &proxies.datas)
	}

	rules.datas = make(map[string][]Rule)
	rules.mux = new(sync.RWMutex)
	rulevalue := defaults.StringForKey(ruleDefaultsKey)
	if rulevalue != "" {
		json.Unmarshal([]byte(rulevalue), &rules.datas)
	}

	server.Start()
}

type ProxyServer struct {
	mux *sync.Mutex

	http *http.Server

	httpListener net.Listener
	sockListener net.Listener

	regexs *sync.Map
}

func NewProxyServer() *ProxyServer {
	server := new(ProxyServer)
	server.mux = new(sync.Mutex)
	server.http = new(http.Server)
	server.regexs = new(sync.Map)
	return server
}

func (g *ProxyServer) RemoveRegex(pattern string) {
	g.regexs.Delete(pattern)
}

func (g *ProxyServer) Start() (err error) {
	g.mux.Lock()
	defer g.mux.Unlock()

	g.httpListener, err = net.Listen("tcp", "0.0.0.0:48557")
	if err != nil {
		return err
	}

	g.sockListener, err = net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return err
	}

	g.http.Handler = http.HandlerFunc(g.handle)
	go g.http.Serve(g.httpListener)
	go func() {
		for {
			conn, err := g.sockListener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}
				continue
			}

			go g.handleSOCKS(conn)
		}
	}()
	g.setupSystem()
	fmt.Println("http proxy listen on", g.httpListener.Addr().String())
	fmt.Println("sock proxy listen on", g.sockListener.Addr().String())
	return nil
}

func (g *ProxyServer) setupSystem() {

}

func (g *ProxyServer) aaaaa() {

}

func (g *ProxyServer) Shutdown() {
	g.mux.Lock()
	defer g.mux.Unlock()

	g.http.Shutdown(context.Background())
	g.sockListener.Close()
}

func (g *ProxyServer) handle(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if r.Method == http.MethodConnect {
		g.handleHTTPS(w, r)
	} else {
		g.handleHTTP(w, r)
	}
}

var unique sync.Map

func (g *ProxyServer) getProxyForRequest(r *http.Request) string {
	proxies := proxies.Fetch()
	for proxyId, rules := range rules.FetchAll() {
		for _, rule := range rules {
			if !rule.T || rule.P == "" {
				continue
			}

			regex, ok := g.regexs.Load(rule.P)
			if !ok {
				regex = regexp.MustCompile(rule.Pattern())
				g.regexs.Store(rule.P, regex)
			}

			host := r.Host
			if strings.HasSuffix(host, ":443") {
				host = host[:len(host)-4]
			} else if strings.HasSuffix(host, ":80") {
				host = host[:len(host)-3]
			}
			if !regex.(*regexp.Regexp).MatchString(host) {
				continue
			}

			for _, p := range proxies {
				if p.InUse && p.ID == proxyId {
					return p.URL()
				}
			}
		}
	}

	if _, ok := unique.LoadOrStore(r.Host, struct{}{}); !ok {
		fmt.Println("未匹配", r.Host)
	}
	return ""
}

func (g *ProxyServer) handleHTTP(w http.ResponseWriter, r *http.Request) {
	var transport http.RoundTripper
	proxyAddr := g.getProxyForRequest(r)
	if proxyAddr != "" {
		proxyURL, err := url.Parse(proxyAddr)
		if err != nil {
			http.Error(w, "Invalid proxy address", http.StatusInternalServerError)
			return
		}

		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	} else {
		transport = &http.Transport{}
	}

	client := &http.Client{Timeout: time.Minute * 3, Transport: transport}
	resp, err := client.Do(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (g *ProxyServer) handleHTTPS(w http.ResponseWriter, r *http.Request) {
	proxyAddr := g.getProxyForRequest(r)
	var destConn net.Conn
	var err error
	if proxyAddr != "" {
		proxyURL, err := url.Parse(proxyAddr)
		if err != nil {
			http.Error(w, "Invalid proxy address", http.StatusInternalServerError)
			return
		}

		destConn, err = net.DialTimeout("tcp", proxyURL.Host, time.Second*5)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		fmt.Fprintf(destConn, "CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", r.Host, r.Host)
		br := bufio.NewReader(destConn)
		resp, err := http.ReadResponse(br, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		if resp.StatusCode != 200 {
			http.Error(w, fmt.Sprintf("Proxy responded with non 200 status: %s", resp.Status), resp.StatusCode)
			return
		}
	} else {
		destConn, err = net.DialTimeout("tcp", r.Host, time.Second*5)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}

func (g *ProxyServer) handleSOCKS(conn net.Conn) {
	defer conn.Close()

	// Read the SOCKS handshake to determine the target
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Printf("Failed to read SOCKS handshake: %v", err)
		return
	}

	// Create a mock request to use with the proxy selector
	mockRequest := &http.Request{
		Header: make(http.Header),
	}

	// Extract destination from SOCKS request (simplified)
	if n > 3 && buffer[0] == 5 { // SOCKS5
		mockRequest.Host = fmt.Sprintf("%s:%d", net.IP(buffer[4:8]).String(), uint16(buffer[8])<<8|uint16(buffer[9]))
	}

	var proxyConn net.Conn
	proxyAddr := g.getProxyForRequest(mockRequest)
	if proxyAddr != "" {
		proxyConn, err = net.Dial("tcp", proxyAddr)
		if err != nil {
			log.Printf("Failed to connect to proxy: %v", err)
			return
		}
		defer proxyConn.Close()

		// Forward the initial SOCKS handshake
		proxyConn.Write(buffer[:n])
	} else {
		targetAddr := mockRequest.Host
		proxyConn, err = net.Dial("tcp", targetAddr)
		if err != nil {
			log.Printf("Failed to connect to target: %v", err)
			return
		}
		defer proxyConn.Close()
		// Handle SOCKS protocol with the client
		// This is a simplified version and needs to be expanded for full SOCKS support
		conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	}

	// Transfer data between connections
	go io.Copy(proxyConn, conn)
	io.Copy(conn, proxyConn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("transfer panic", err)
		}
	}()
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}
