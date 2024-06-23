package main

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/charliego3/proxies/utility"
	"github.com/progrium/macdriver/helper/action"
	"github.com/progrium/macdriver/helper/layout"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/objc"
)

type ProxySheet struct {
	appkit.IView
	p appkit.IWindow
}

func NewProxySheet(p appkit.IWindow) ProxySheet {
	return ProxySheet{IView: appkit.NewView(), p: p}
}

func (c ProxySheet) Init(title string, proxy Proxy) {
	c.p.SetContentView(c)
	label := appkit.NewLabel(title)
	label.SetTranslatesAutoresizingMaskIntoConstraints(false)
	c.AddSubview(label)
	layout.PinAnchorTo(label.LeadingAnchor(), c.LeadingAnchor(), 20)
	layout.PinAnchorTo(label.TopAnchor(), c.TopAnchor(), 15)

	cancel := appkit.NewPushButton("Cancel")
	cancel.SetBezelStyle(appkit.BezelStyleRounded)
	cancel.SetTranslatesAutoresizingMaskIntoConstraints(false)
	action.Set(cancel, func(sender objc.Object) { Window.EndSheet(c.p) })
	c.AddSubview(cancel)
	layout.SetMinWidth(cancel, 100)
	layout.PinAnchorTo(cancel.LeadingAnchor(), c.LeadingAnchor(), 20)
	layout.PinAnchorTo(cancel.BottomAnchor(), c.BottomAnchor(), -15)

	ok := appkit.NewPushButton("Save")
	ok.SetBezelStyle(appkit.BezelStyleRounded)
	ok.SetState(appkit.MixedState)
	ok.SetTranslatesAutoresizingMaskIntoConstraints(false)
	c.AddSubview(ok)
	layout.SetMinWidth(ok, 100)
	layout.PinAnchorTo(ok.TrailingAnchor(), c.TrailingAnchor(), -20)
	layout.PinAnchorTo(ok.BottomAnchor(), c.BottomAnchor(), -15)

	box := appkit.NewBox()
	box.SetTranslatesAutoresizingMaskIntoConstraints(false)
	box.SetBoxType(appkit.BoxCustom)
	box.SetCornerRadius(5)
	box.SetBorderWidth(1)
	box.SetBorderColor(appkit.Color_SeparatorColor())
	box.SetContentViewMargins(utility.SizeOf(0, 0))
	utility.AddAppearanceObserver(func() {
		box.SetFillColor(utility.ColorWithAppearance(appkit.Color_WhiteColor(), utility.ColorHex("#303030")))
	})
	c.AddSubview(box)
	layout.PinAnchorTo(box.LeadingAnchor(), c.LeadingAnchor(), 20)
	layout.PinAnchorTo(box.TopAnchor(), label.BottomAnchor(), 15)
	layout.PinAnchorTo(box.TrailingAnchor(), c.TrailingAnchor(), -20)
	layout.PinAnchorTo(box.BottomAnchor(), ok.TopAnchor(), -15)

	types := appkit.NewPopUpButton()
	types.AddItemsWithTitles([]string{
		"HTTP",
		"HTTPS",
	})
	types.SelectItemWithTitle(proxy.Schema)
	types.SetTranslatesAutoresizingMaskIntoConstraints(false)
	layout.SetMinWidth(types, 250)
	namei := appkit.NewTextField()
	namei.SetStringValue(proxy.Name)
	hosti := appkit.NewTextField()
	hosti.SetStringValue(proxy.Host)
	porti := appkit.NewTextField()
	if proxy.Port > 0 {
		porti.SetStringValue(strconv.Itoa(proxy.Port))
	}
	typesHandler := func(sender objc.Object) {
		if "HTTPS" == types.TitleOfSelectedItem() {
			hosti.SetPlaceholderString("example.com")
			porti.SetPlaceholderString("443")
		} else {
			hosti.SetPlaceholderString("192.168.1.2")
			porti.SetPlaceholderString("80")
		}
	}
	typesHandler(objc.NewObject())
	action.Set(types, typesHandler)

	authed := "Username & Password"
	authentication := appkit.NewPopUpButton()
	authentication.AddItemsWithTitles([]string{"None", authed})
	authentication.SetTranslatesAutoresizingMaskIntoConstraints(false)
	var paddingTop float64 = 75
	if proxy.Auth {
		authentication.SelectItemWithTitle(authed)
		paddingTop = 50
	}
	layout.SetMinWidth(authentication, 250)

	unamel := appkit.NewLabel("Username:")
	unamei := appkit.NewTextField()
	unamei.SetStringValue(proxy.Username)
	passl := appkit.NewLabel("Password:")
	passi := appkit.NewTextField()
	passi.SetStringValue(proxy.Password)
	grid := appkit.GridView_GridViewWithViews([][]appkit.IView{
		{appkit.NewLabel("Name:"), namei},
		{appkit.NewLabel("Type:"), types},
		{appkit.NewLabel("Host:"), hosti},
		{appkit.NewLabel("Port:"), porti},
		{appkit.NewLabel("Authenication"), authentication},
		{unamel, unamei},
		{passl, passi},
	})
	grid.SetColumnSpacing(7)
	grid.SetRowSpacing(10)
	grid.SetXPlacement(appkit.GridCellPlacementTrailing)
	grid.SetTranslatesAutoresizingMaskIntoConstraints(false)
	box.AddSubview(grid)
	layout.SetMinWidth(grid, 300)
	layout.AliginCenterX(grid, box)
	constraint := grid.TopAnchor().ConstraintEqualToAnchorConstant(box.TopAnchor(), paddingTop)
	constraint.SetActive(true)

	authenticationHandler := func(sender objc.Object) {
		paddingTop = 50
		if authentication.TitleOfSelectedItem() == authed {
			unamei.SetHidden(false)
			unamel.SetHidden(false)
			passl.SetHidden(false)
			passi.SetHidden(false)
		} else {
			unamei.SetHidden(true)
			unamel.SetHidden(true)
			passl.SetHidden(true)
			passi.SetHidden(true)
			paddingTop = 75
		}
		constraint.SetActive(false)
		constraint = grid.TopAnchor().ConstraintEqualToAnchorConstant(box.TopAnchor(), paddingTop)
		constraint.SetActive(true)
		grid.SetNeedsUpdateConstraints(true)
	}
	authenticationHandler(objc.NewObject())
	action.Set(authentication, authenticationHandler)
	action.Set(ok, func(sender objc.Object) {
		ok.SetEnabled(false)
		showWaring := func(message string) {
			utility.ShowAlert(
				utility.WithAlertTitle("Oops!"),
				utility.WithAlertMessage("Options are invalid: "+message),
				utility.WithAlertWindow(c.p),
				utility.WithAlertStyle(appkit.AlertStyleWarning),
			)
			ok.SetEnabled(true)
			ok.SetState(appkit.MixedState)
		}
		name := namei.StringValue()
		port := porti.StringValue()
		host := hosti.StringValue()
		if name == "" {
			showWaring("Name can not be empty")
			return
		}
		if host == "" {
			showWaring("Host can not be empty")
			return
		}
		schema := types.TitleOfSelectedItem()
		if port == "" {
			if schema == "HTTP" {
				port = "80"
			} else {
				port = "443"
			}
		}
		_, err := url.Parse(fmt.Sprintf("%s://%s:%s", strings.ToLower(schema), host, port))
		if err != nil {
			showWaring(err.Error())
			return
		}
		proxy.Name = name
		proxy.Schema = schema
		proxy.Host = host
		proxy.Port, _ = strconv.Atoi(port)
		proxy.Auth = authentication.TitleOfSelectedItem() == authed
		proxy.Username = unamei.StringValue()
		proxy.Password = passi.StringValue()
		proxies.Update(proxy)
		Window.Sidebar.Update()
		Window.Sidebar.ScrollToBottom()
		Window.EndSheetReturnCode(c.p, appkit.ModalResponseOK)
	})
}

func (ProxySheet) Handler(code appkit.ModalResponse) {}

func OpenProxySheet(title string, proxy Proxy) {
	panel := appkit.NewWindowWithSizeAndStyle(
		500, 400,
		appkit.WindowStyleMaskTitled|
			appkit.WindowStyleMaskFullSizeContentView,
	)

	panel.SetMenu(appkit.NewMenuWithTitle("title string"))
	creator := NewProxySheet(panel)
	creator.Init(title, proxy)
	Window.BeginSheetCompletionHandler(panel, creator.Handler)
}
