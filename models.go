package main

import (
	"fmt"
	"strings"
)

type Proxy struct {
	ID       string
	Name     string
	Type     string
	Host     string
	Port     int
	Auth     bool
	Username string
	Password string
	InUse    bool
}

func (p Proxy) URL() string {
	return strings.ToLower(fmt.Sprintf("%s://%s:%d", p.Type, p.Host, p.Port))
}

func NewProxy() Proxy {
	return Proxy{Type: "HTTP"}
}

type Rule struct {
	N string
	T bool
	R string
}

var RuleDataSource = make(map[string]*Rule)
