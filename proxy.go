package main

import (
	"fmt"

	"github.com/charliego3/proxies/utility"
	"github.com/progrium/macdriver/helper/action"
	"github.com/progrium/macdriver/helper/layout"
	"github.com/progrium/macdriver/helper/widgets"
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
	label := appkit.NewLabel(title)
	label.SetTranslatesAutoresizingMaskIntoConstraints(false)
	c.AddSubview(label)
	layout.PinAnchorTo(label.LeadingAnchor(), c.LeadingAnchor(), 20)
	layout.PinAnchorTo(label.TopAnchor(), c.TopAnchor(), 15)

	cancel := appkit.NewPushButton("Cancel")
	cancel.SetTranslatesAutoresizingMaskIntoConstraints(false)
	action.Set(cancel, func(sender objc.Object) { c.w.EndSheet(c.p) })
	c.AddSubview(cancel)
	layout.SetMinWidth(cancel, 100)
	layout.PinAnchorTo(cancel.LeadingAnchor(), c.LeadingAnchor(), 20)
	layout.PinAnchorTo(cancel.BottomAnchor(), c.BottomAnchor(), -15)

	ok := appkit.NewPushButton("Save")
	ok.SetTranslatesAutoresizingMaskIntoConstraints(false)
	ok.SetBezelStyle(appkit.BezelStyleRounded)
	ok.SetBezelColor(appkit.Color_BlueColor())
	action.Set(ok, func(sender objc.Object) { c.w.EndSheetReturnCode(c.p, appkit.ModalResponseOK) })
	c.AddSubview(ok)
	layout.SetMinWidth(ok, 100)
	layout.PinAnchorTo(ok.TrailingAnchor(), c.TrailingAnchor(), -20)
	layout.PinAnchorTo(ok.BottomAnchor(), c.BottomAnchor(), -15)

	box := appkit.NewBox()
	box.SetTranslatesAutoresizingMaskIntoConstraints(false)
	box.SetCornerRadius(5)
	box.SetBorderWidth(1)
	box.SetBorderColor(appkit.Color_SeparatorColor())
	box.SetFillColor(utility.ColorHex("#303030"))
	box.SetBoxType(appkit.BoxCustom)
	c.AddSubview(box)
	layout.PinAnchorTo(box.LeadingAnchor(), c.LeadingAnchor(), 20)
	layout.PinAnchorTo(box.TopAnchor(), label.BottomAnchor(), 15)
	layout.PinAnchorTo(box.TrailingAnchor(), c.TrailingAnchor(), -20)
	layout.PinAnchorTo(box.BottomAnchor(), ok.TopAnchor(), -15)

	// form := widgets.NewFormView()
	gv := appkit.GridView_GridViewWithNumberOfColumnsRows(2, 0)
	gv.SetTranslatesAutoresizingMaskIntoConstraints(false)
	gv.SetColumnSpacing(10)
	gv.SetRowSpacing(10)
	gv.SetContentHuggingPriorityForOrientation(654.0, appkit.LayoutConstraintOrientationVertical)
	gv.SetContentHuggingPriorityForOrientation(654.0, appkit.LayoutConstraintOrientationHorizontal)
	form := &widgets.FormView{
		GridView: gv,
	}
	form.AddRow("name string", appkit.NewTextField())
	box.AddSubview(form)
	layout.AliginCenterX(form, box)
	layout.AliginCenterY(form, box)
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
