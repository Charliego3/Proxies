package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/progrium/macdriver/macos/foundation"
)

const ()

type Proxy struct {
	ID       string
	Name     string
	Schema   string
	Host     string
	Port     int
	Auth     bool
	Username string
	Password string
	InUse    bool
}

func (p Proxy) URL() string {
	return strings.ToLower(fmt.Sprintf("%s://%s:%d", p.Schema, p.Host, p.Port))
}

func NewProxy() Proxy {
	return Proxy{Schema: "HTTP"}
}

type ProxyDatasource struct {
	proxies []Proxy
	mux     *sync.RWMutex
}

func (d *ProxyDatasource) ByIndex(index int) Proxy {
	if index < 0 {
		return Proxy{}
	}

	d.mux.RLock()
	defer d.mux.RUnlock()

	if len(d.proxies)-1 < index {
		return Proxy{}
	}
	return d.proxies[index]
}

func (d *ProxyDatasource) Length() int {
	d.mux.RLock()
	defer d.mux.RUnlock()
	return len(d.proxies)
}

func (d *ProxyDatasource) Fetch() (proxies []Proxy) {
	d.mux.RLock()
	defer d.mux.RUnlock()

	proxies = make([]Proxy, len(d.proxies))
	for i, p := range d.proxies {
		proxies[i] = p
	}
	return
}

func (d *ProxyDatasource) Update(proxy Proxy) {
	d.mux.Lock()
	defer d.mux.Unlock()

	if proxy.ID == "" {
		id, _ := uuid.NewV7()
		proxy.ID = id.String()
	}

	var updated bool
	for i, p := range d.proxies {
		if p.ID == proxy.ID {
			d.proxies = append(d.proxies[:i], append([]Proxy{proxy}, d.proxies[i+1:]...)...)
			updated = true
		}
	}

	if !updated {
		d.proxies = append(d.proxies, proxy)
	}
	d.write()
}

func (d *ProxyDatasource) Delete(index int) {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.proxies = append(d.proxies[:index], d.proxies[index+1:]...)
	d.write()
}

func (d *ProxyDatasource) write() {
	buf, _ := json.Marshal(d.proxies)
	defaults.SetObjectForKey(foundation.String_StringWithString(string(buf)), proxyDefaultsKey)
}

type Rule struct {
	ID string
	N  string // name
	T  bool   // state
	R  string // remark
}

type RuleDatasource struct {
	rules map[string][]Rule
	mux   *sync.RWMutex
}

func (d *RuleDatasource) Fetch() (rules []Rule) {
	d.mux.RLock()
	defer d.mux.RUnlock()

	if rs, ok := d.rules[Window.Sidebar.SelectedId()]; ok {
		rules = make([]Rule, len(rs))
		for i, r := range rs {
			rules[i] = r
		}
	}
	return
}

func (d *RuleDatasource) ByIndex(index int) Rule {
	d.mux.RLock()
	defer d.mux.RUnlock()

	if rs, ok := d.rules[Window.Sidebar.SelectedId()]; ok && len(rs) > index {
		return rs[index]
	}
	return Rule{}
}

func (d *RuleDatasource) Update(rule Rule) {
	d.mux.Lock()
	defer d.mux.Unlock()

	proxyId := Window.Sidebar.SelectedId()
	if rules, ok := d.rules[proxyId]; ok {
		for i := range rules {
			if rules[i].ID == rule.ID {
				rules[i] = rule
			}
		}
		d.rules[proxyId] = rules
	}
	d.write()
}

func (d *RuleDatasource) Add(rule Rule) {
	d.mux.Lock()
	defer d.mux.Unlock()

	if rule.ID == "" {
		id, _ := uuid.NewV7()
		rule.ID = id.String()
	}

	proxyId := Window.Sidebar.SelectedId()
	if rules, ok := d.rules[proxyId]; ok {
		d.rules[proxyId] = append(rules, rule)
	} else {
		d.rules[proxyId] = []Rule{rule}
	}
	d.write()
}

func (d *RuleDatasource) Delete(index int) {
	if index < 0 {
		return
	}

	d.mux.Lock()
	defer d.mux.Unlock()

	proxyId := Window.Sidebar.SelectedId()
	if rules, ok := d.rules[proxyId]; ok {
		d.rules[proxyId] = append(rules[:index], rules[index+1:]...)
	}
	d.write()
}

func (d *RuleDatasource) DeleteById(id string) {
	d.mux.Lock()
	defer d.mux.Unlock()

	proxyId := Window.Sidebar.SelectedId()
	if rules, ok := d.rules[proxyId]; ok {
		for i, r := range rules {
			if r.ID != id {
				continue
			}
			d.rules[proxyId] = append(rules[:i], rules[i+1:]...)
		}
	}
	d.write()
}

func (d *RuleDatasource) Remove() {
	d.mux.Lock()
	defer d.mux.Unlock()

	delete(d.rules, Window.Sidebar.SelectedId())
	d.write()
}

func (d *RuleDatasource) write() {
	buf, _ := json.Marshal(d.rules)
	defaults.SetObjectForKey(foundation.String_StringWithString(string(buf)), ruleDefaultsKey)
}

func (d *RuleDatasource) LastIndex() int {
	proxyId := Window.Sidebar.SelectedId()
	return len(d.rules[proxyId]) - 1
}
