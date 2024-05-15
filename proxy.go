package main

import (
	"fmt"

	"github.com/charliego3/proxies/utility"
	"github.com/progrium/macdriver/helper/action"
	"github.com/progrium/macdriver/helper/layout"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/objc"
)

type Creator struct {
	appkit.IView
	w appkit.IWindow
	p appkit.IWindow
}

func NewCreator(w appkit.IWindow, p appkit.IWindow) Creator {
	return Creator{IView: appkit.NewView(), w: w, p: p}
}

func (c Creator) Init(title string) {
	c.p.SetContentView(c)
	label := appkit.NewLabel(title)
	label.SetTranslatesAutoresizingMaskIntoConstraints(false)
	c.AddSubview(label)
	layout.PinAnchorTo(label.LeadingAnchor(), c.LeadingAnchor(), 20)
	layout.PinAnchorTo(label.TopAnchor(), c.TopAnchor(), 15)

	cancel := appkit.NewPushButton("Cancel")
	cancel.SetBezelStyle(appkit.BezelStyleRounded)
	cancel.SetTranslatesAutoresizingMaskIntoConstraints(false)
	action.Set(cancel, func(sender objc.Object) { c.w.EndSheet(c.p) })
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
	types.SetTranslatesAutoresizingMaskIntoConstraints(false)
	layout.SetMinWidth(types, 250)
	namei := appkit.NewTextField()
	namei.SetToolTip("toop tips")
	hosti := appkit.NewTextField()
	porti := appkit.NewTextField()
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

	grid := appkit.GridView_GridViewWithViews([][]appkit.IView{
		{appkit.NewLabel("Name:"), namei},
		{appkit.NewLabel("Type:"), types},
		{appkit.NewLabel("Host:"), hosti},
		{appkit.NewLabel("Port:"), porti},
	})
	grid.SetColumnSpacing(10)
	grid.SetRowSpacing(10)
	grid.SetXPlacement(appkit.GridCellPlacementTrailing)
	grid.SetTranslatesAutoresizingMaskIntoConstraints(false)
	box.AddSubview(grid)
	layout.SetMinWidth(grid, 300)
	layout.AliginCenterX(grid, box)
	layout.AliginCenterY(grid, box)

	action.Set(ok, func(sender objc.Object) {
		c.w.EndSheetReturnCode(c.p, appkit.ModalResponseOK)
	})
}

func (Creator) Handler(code appkit.ModalResponse) {
	fmt.Println(code, ".......")
}

func OpenNewPanelSheet(w appkit.IWindow) {
	panel := appkit.NewWindowWithSizeAndStyle(
		500, 400,
		appkit.WindowStyleMaskTitled|
			appkit.WindowStyleMaskFullSizeContentView,
	)

	creator := NewCreator(w, panel)
	creator.Init("Choose options for your new Proxy:")
	w.BeginSheetCompletionHandler(panel, creator.Handler)
}
