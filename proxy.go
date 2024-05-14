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
	p appkit.Panel
}

func NewCreator(w appkit.IWindow, p appkit.Panel) Creator {
	return Creator{IView: appkit.NewView(), w: w, p: p}
}

func (c Creator) Init(title string) {
	c.p.SetContentView(c)
	c.p.MakeFirstResponder(c.w.FirstResponder())
	fmt.Println(c.p.AcceptsFirstResponder(), c.p.BecomeFirstResponder())
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
	ok.SetTranslatesAutoresizingMaskIntoConstraints(false)
	action.Set(ok, func(sender objc.Object) {
		c.w.EndSheetReturnCode(c.p, appkit.ModalResponseOK)
	})
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

	i := appkit.NewTextField()
	i.SetEditable(true)
	i.SetEnabled(true)
	i.SetTranslatesAutoresizingMaskIntoConstraints(false)
	c.AddSubview(i)
	layout.SetMinWidth(i, 100)
	layout.SetMinHeight(i, 30)
	layout.PinAnchorTo(i.TopAnchor(), c.TopAnchor(), 5)
	layout.PinAnchorTo(i.TrailingAnchor(), c.TrailingAnchor(), -5)

	content := appkit.NewViewWithFrame(utility.RectOf(utility.SizeOf(400, 300)))
	content.SetTranslatesAutoresizingMaskIntoConstraints(false)
	content.SetWantsLayer(true)
	content.Layer().SetBorderWidth(1)
	content.Layer().SetBorderColor(appkit.Color_BlackColor().CGColor())
	content.Layer().SetCornerRadius(5)
	box.AddSubview(content)
	layout.SetMinWidth(content, 400)
	layout.SetMinHeight(content, 300)
	layout.AliginCenterX(content, box)
	layout.AliginCenterY(content, box)

	name := appkit.NewLabel("Proxy name")
	name.SetTranslatesAutoresizingMaskIntoConstraints(false)
	content.AddSubview(name)
	layout.AliginLeading(name, content)
	layout.AliginTop(name, content)

	ni := appkit.NewTextField()
	ni.SetTranslatesAutoresizingMaskIntoConstraints(false)
	ni.SetEditable(true)
	content.AddSubview(ni)
	layout.PinAnchorTo(ni.LeadingAnchor(), name.TrailingAnchor(), 20)
	layout.PinAnchorTo(ni.TrailingAnchor(), content.TrailingAnchor(), -5)
	layout.AliginTop(ni, content)

	t := appkit.NewLabel("Proxy Type")
	t.SetTranslatesAutoresizingMaskIntoConstraints(false)
	content.AddSubview(t)
	layout.AliginLeading(t, content)
	layout.PinAnchorTo(t.TopAnchor(), name.BottomAnchor(), 5)

	ti := appkit.NewTextField()
	ti.SetEnabled(true)
	ti.SetEditable(true)
	ti.SetTranslatesAutoresizingMaskIntoConstraints(false)
	content.AddSubview(ti)
	layout.PinAnchorTo(ti.LeadingAnchor(), t.TrailingAnchor(), 20)
	layout.PinAnchorTo(ti.TopAnchor(), ni.BottomAnchor(), 5)
	layout.PinAnchorTo(ti.TrailingAnchor(), content.TrailingAnchor(), -5)
	layout.AliginTrailing(ti, content)
}

func (Creator) Handle(code appkit.ModalResponse) {
	fmt.Println(code, ".......")
}

func OpenNewPanelSheet(w appkit.IWindow) {
	panel := appkit.NewPanelWithContentRectStyleMaskBackingDefer(
		utility.RectOf(utility.SizeOf(600, 500)),
		appkit.WindowStyleMaskFullSizeContentView,
		appkit.BackingStoreBuffered,
		false,
	)

	creator := NewCreator(w, panel)
	creator.Init("Choose options for your new Proxy")
	w.BeginSheetCompletionHandler(panel, creator.Handle)
}
