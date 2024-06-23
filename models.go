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
	datas []Proxy
	mux   *sync.RWMutex
}

func (d *ProxyDatasource) ByIndex(index int) Proxy {
	if index < 0 {
		return Proxy{}
	}

	d.mux.RLock()
	defer d.mux.RUnlock()

	if len(d.datas)-1 < index {
		return Proxy{}
	}
	return d.datas[index]
}

func (d *ProxyDatasource) Length() int {
	d.mux.RLock()
	defer d.mux.RUnlock()
	return len(d.datas)
}

func (d *ProxyDatasource) Fetch() (proxies []Proxy) {
	d.mux.RLock()
	defer d.mux.RUnlock()

	proxies = make([]Proxy, len(d.datas))
	for i, p := range d.datas {
		proxies[i] = p
	}
	return
}

func (d *ProxyDatasource) AnyUsing() bool {
	d.mux.RLock()
	defer d.mux.RUnlock()

	for _, p := range d.datas {
		if p.InUse {
			return true
		}
	}
	return false
}

func (d *ProxyDatasource) Update(proxy Proxy) {
	d.mux.Lock()
	defer d.mux.Unlock()

	if proxy.ID == "" {
		id, _ := uuid.NewV7()
		proxy.ID = id.String()
	}

	var updated bool
	for i, p := range d.datas {
		if p.ID == proxy.ID {
			d.datas = append(d.datas[:i], append([]Proxy{proxy}, d.datas[i+1:]...)...)
			updated = true
		}
	}

	if !updated {
		d.datas = append(d.datas, proxy)
	}
	d.write()
}

func (d *ProxyDatasource) Delete(index int) {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.datas = append(d.datas[:index], d.datas[index+1:]...)
	d.write()
}

func (d *ProxyDatasource) write() {
	buf, _ := json.Marshal(d.datas)
	defaults.SetObjectForKey(foundation.String_StringWithString(string(buf)), proxyDefaultsKey)
}

type Rule struct {
	ID string
	P  string // pattern
	T  bool   // state
	R  string // remark
}

func (r Rule) Pattern() string {
	return fmt.Sprintf("^%s$", r.P)
}

type RuleDatasource struct {
	datas map[string][]Rule
	mux   *sync.RWMutex
}

func (d *RuleDatasource) FetchAll() map[string][]Rule {
	d.mux.RLock()
	defer d.mux.RUnlock()

	result := make(map[string][]Rule, len(rules.datas))
	for k, v := range d.datas {
		result[k] = append([]Rule{}, v...)
	}
	return result
}

func (d *RuleDatasource) Fetch() (rules []Rule) {
	return d.FetchWithProxy(Window.Sidebar.SelectedId())
}

func (d *RuleDatasource) FetchWithProxy(proxyId string) (rules []Rule) {
	d.mux.RLock()
	defer d.mux.RUnlock()

	if rs, ok := d.datas[proxyId]; ok {
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

	if rs, ok := d.datas[Window.Sidebar.SelectedId()]; ok && len(rs) > index {
		return rs[index]
	}
	return Rule{}
}

func (d *RuleDatasource) Update(rule Rule) {
	d.mux.Lock()
	defer d.mux.Unlock()

	proxyId := Window.Sidebar.SelectedId()
	if rules, ok := d.datas[proxyId]; ok {
		for i := range rules {
			if rules[i].ID == rule.ID {
				rules[i] = rule
			}
		}
		d.datas[proxyId] = rules
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
	if rules, ok := d.datas[proxyId]; ok {
		d.datas[proxyId] = append(rules, rule)
	} else {
		d.datas[proxyId] = []Rule{rule}
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
	if rules, ok := d.datas[proxyId]; ok {
		d.datas[proxyId] = append(rules[:index], rules[index+1:]...)
	}
	d.write()
}

func (d *RuleDatasource) DeleteById(id string) {
	d.mux.Lock()
	defer d.mux.Unlock()

	proxyId := Window.Sidebar.SelectedId()
	if rules, ok := d.datas[proxyId]; ok {
		for i, r := range rules {
			if r.ID != id {
				continue
			}
			d.datas[proxyId] = append(rules[:i], rules[i+1:]...)
		}
	}
	d.write()
}

func (d *RuleDatasource) Remove() {
	d.mux.Lock()
	defer d.mux.Unlock()

	delete(d.datas, Window.Sidebar.SelectedId())
	d.write()
}

func (d *RuleDatasource) write() {
	buf, _ := json.Marshal(d.datas)
	defaults.SetObjectForKey(foundation.String_StringWithString(string(buf)), ruleDefaultsKey)
}

func (d *RuleDatasource) LastIndex() int {
	proxyId := Window.Sidebar.SelectedId()
	return len(d.datas[proxyId]) - 1
}
