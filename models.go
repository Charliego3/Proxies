package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/charliego3/proxies/utility"
	"github.com/google/uuid"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
	"github.com/progrium/macdriver/macos/uti"
	"github.com/progrium/macdriver/objc"
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

type portObj struct {
	Proxy `json:",inline"`
	Rules []Rule `json:"rules"`
}

func Import(objc.Object) {
	app.ActivateIgnoringOtherApps(true)
	panel := appkit.NewOpenPanelWithContentRectStyleMaskBackingDefer(
		utility.RectOf(utility.SizeOf(500, 1000)),
		appkit.WindowStyleMaskTitled|appkit.WindowStyleMaskResizable,
		appkit.BackingStoreBuffered, false)
	panel.SetTitle("Import Proxies")
	panel.SetAllowedContentTypes([]uti.IType{uti.Type_TypeWithFilenameExtension("json")})
	handler := func(result appkit.ModalResponse) {
		if result != appkit.ModalResponseOK {
			return
		}

		contents := foundation.String_StringWithContentsOfURLEncodingError(
			panel.URL(), foundation.String_DefaultCStringEncoding(), nil)
		var list []portObj
		if err := json.Unmarshal([]byte(contents.String()), &list); err != nil {
			return
		}

		for _, obj := range list {
			proxies.mux.Lock()
			var exists bool
			for i, p := range proxies.datas {
				if p.ID == obj.ID {
					exists = true
					proxies.datas[i] = obj.Proxy
				}
			}

			if !exists {
				proxies.datas = append(proxies.datas, obj.Proxy)
			}
			proxies.write()
			proxies.mux.Unlock()

			rules.mux.Lock()
			for j, r := range rules.datas[obj.ID] {
				for j2, r2 := range obj.Rules {
					if r2.ID == r.ID {
						rules.datas[obj.ID][j] = r2
						obj.Rules = append(obj.Rules[:j2], obj.Rules[j2+1:]...)
					}
				}
			}
			rules.datas[obj.ID] = append(rules.datas[obj.ID], obj.Rules...)
			rules.write()
			rules.mux.Unlock()
		}

		app.launchWindow(objc.Object{})
		Window.Sidebar.outline.ReloadData()
		Window.Rules.ReloadData()
		Window.Sidebar.SelectRow(0)
	}

	if Window == nil {
		panel.BeginWithCompletionHandler(handler)
	} else {
		panel.BeginSheetModalForWindowCompletionHandler(Window, handler)
	}
}

func Export(objc.Object) {
	var list []portObj
	proxies := proxies.Fetch()
	for _, p := range proxies {
		rules := rules.FetchWithProxy(p.ID)
		list = append(list, portObj{Proxy: p, Rules: rules})
	}

	app.ActivateIgnoringOtherApps(true)
	if len(list) == 0 {
		utility.ShowAlert(
			utility.WithAlertTitle("Oops!"),
			utility.WithAlertMessage("There is no proxy for export"),
		)
		return
	}

	path, _ := os.UserHomeDir()
	panel := appkit.NewSavePanelWithContentRectStyleMaskBackingDefer(
		utility.RectOf(utility.SizeOf(1000, 1000)),
		appkit.WindowStyleMaskTitled|appkit.WindowStyleMaskResizable,
		appkit.BackingStoreBuffered, false)
	panel.SetTitle("Export Proxies")
	panel.SetCanCreateDirectories(true)
	panel.SetDirectoryURL(foundation.URL_FileURLWithPathComponents([]string{path, "Downloads"}))
	panel.SetNameFieldStringValue("proxies.json")

	handler := func(result appkit.ModalResponse) {
		if result != appkit.ModalResponseOK {
			return
		}

		buf, _ := json.Marshal(list)
		foundation.String_StringWithString(string(buf)).WriteToURLAtomicallyEncodingError(
			panel.URL(),
			true,
			foundation.String_DefaultCStringEncoding(),
			nil,
		)
	}
	if Window == nil {
		panel.BeginWithCompletionHandler(handler)
	} else {
		panel.BeginSheetModalForWindowCompletionHandler(Window, handler)
	}
}
