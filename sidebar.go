package main

import (
	"fmt"

	"github.com/charliego3/proxies/lib"
	"github.com/charliego3/proxies/utility"
	"github.com/progrium/macdriver/helper/layout"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
	"github.com/progrium/macdriver/objc"
)

const sidebarIdentifier = "sidebarDatasourceIdentifier"

var items []objc.Object

func init() {
	for i := 65; i < 65+26; i++ {
		items = append(items, foundation.String_StringWithString(fmt.Sprintf("%s -- %d", string(rune(i)), i)).Object)
	}
}

type Sidebar struct {
	appkit.IViewController
	w       appkit.IWindow
	view    appkit.IView
	outline appkit.OutlineView
	max     appkit.LayoutConstraint
}

func NewSidebarController(w appkit.IWindow) *Sidebar {
	sidebar := new(Sidebar)
	sidebar.w = w
	sidebar.IViewController = appkit.NewViewController()
	sidebar.Init()
	return sidebar
}

func (s *Sidebar) Init() {
	s.outline = appkit.NewOutlineView()
	s.outline.SetColumnAutoresizingStyle(appkit.TableViewSequentialColumnAutoresizingStyle)
	s.outline.SetUsesAlternatingRowBackgroundColors(false)
	s.outline.SetStyle(appkit.TableViewStyleSourceList)
	s.outline.SetSelectionHighlightStyle(appkit.TableViewSelectionHighlightStyleSourceList)
	s.outline.SetUsesSingleLineMode(true)
	s.outline.SetAllowsColumnSelection(false)
	s.outline.SetAutoresizingMask(appkit.ViewWidthSizable)
	s.outline.SetHeaderView(nil)
	s.outline.AddTableColumn(utility.TableColumn(sidebarIdentifier, ""))

	menu := appkit.NewMenu()
	menu.AddItem(appkit.NewMenuItemWithAction("Edit", "", func(sender objc.Object) {
		fmt.Println("clicked edit")
	}))
	menu.AddItem(appkit.NewMenuItemWithAction("Open", "", func(sender objc.Object) {
		fmt.Println("clicked open")
	}))

	s.outline.SetMenu(menu)
	s.setDelegate()
	s.setDatasource()
	scrollView := appkit.NewScrollView()
	clipView := appkit.ClipViewFrom(scrollView.ContentView().Ptr())
	clipView.SetDocumentView(s.outline)
	clipView.SetAutomaticallyAdjustsContentInsets(false)
	clipView.SetContentInsets(foundation.EdgeInsets{Top: 10})

	s.outline.SelectRowIndexesByExtendingSelection(foundation.NewIndexSetWithIndex(0), true)
	scrollView.SetBorderType(appkit.NoBorder)
	scrollView.SetDrawsBackground(false)
	scrollView.SetAutohidesScrollers(true)
	scrollView.SetHasVerticalScroller(true)
	scrollView.ContentView().ScrollToPoint(foundation.Point{Y: -10})

	s.view = scrollView
	s.IViewController.SetView(s.view)
	s.view.SetFrameSize(utility.SizeOf(260, 0))
	layout.SetMinWidth(s.view, 200)
	s.SetSidebarMaxWidth()
}

func (s *Sidebar) setDelegate() {
	delegate := new(appkit.OutlineViewDelegate)
	delegate.SetOutlineViewViewForTableColumnItem(s.createColumnItem)
	delegate.SetControlTextDidBeginEditing(func(obj foundation.Notification) {})
	delegate.SetOutlineViewHeightOfRowByItem(func(outlineView appkit.OutlineView, item objc.Object) float64 {
		return 30
	})
	po0 := objc.WrapAsProtocol("NSOutlineViewDelegate", appkit.POutlineViewDelegate(delegate))
	objc.SetAssociatedObject(s.outline, objc.AssociationKey("setDelegate"), po0, objc.ASSOCIATION_RETAIN)
	objc.Call[objc.Void](s.outline, objc.Sel("setDelegate:"), po0)
}

func (s *Sidebar) setDatasource() {
	datasource := new(lib.OutlineViewDatasource)
	datasource.SetOutlineViewChildOfItem(func(outlineView appkit.OutlineView, index int, item objc.Object) objc.Object {
		if item.IsNil() {
			return items[index]
		}

		return foundation.StringFrom(item.Ptr()).Object
	})
	datasource.SetOutlineViewIsItemExpandable(func(outlineView appkit.OutlineView, item objc.Object) bool {
		return false
	})
	datasource.SetOutlineViewNumberOfChildrenOfItem(func(_ appkit.OutlineView, item objc.Object) int {
		return len(items)
	})
	appkit.NewTableViewRowAction()
	po1 := objc.WrapAsProtocol("NSOutlineViewDataSource", appkit.POutlineViewDataSource(datasource))
	objc.SetAssociatedObject(s.outline, objc.AssociationKey("setDataSource"), po1, objc.ASSOCIATION_RETAIN)
	objc.Call[objc.Void](s.outline, objc.Sel("setDataSource:"), po1)
}

func (s *Sidebar) createColumnItem(_ appkit.OutlineView, tableColumn appkit.TableColumn, item objc.Object) appkit.View {
	image := appkit.NewImageView()
	image.SetTranslatesAutoresizingMaskIntoConstraints(false)
	image.SetImage(utility.SymbolImage("network", utility.ImageLarge))

	text := appkit.NewTextField()
	text.SetBordered(false)
	text.SetBezelStyle(appkit.TextFieldSquareBezel)
	text.SetEditable(false)
	text.SetDrawsBackground(false)
	text.SetTranslatesAutoresizingMaskIntoConstraints(false)
	text.SetStringValue(foundation.StringFrom(item.Ptr()).CapitalizedString())

	// onoff := appkit.NewButton()
	// onoff.SetButtonType(appkit.ButtonTypeToggle)
	// onoff.SetTranslatesAutoresizingMaskIntoConstraints(false)
	// onoff.SetState(appkit.ControlStateValueOn)

	rowView := appkit.NewTableRowView()
	rowView.AddSubview(image)
	rowView.AddSubview(text)
	// rowView.AddSubview(onoff)

	layout.AliginLeading(image, rowView)
	layout.AliginCenterY(image, rowView)
	layout.PinAnchorTo(text.LeadingAnchor(), image.TrailingAnchor(), 5)
	layout.AliginCenterY(text, rowView)
	// image.LeadingAnchor().ConstraintEqualToAnchor(rowView.LeadingAnchor()).SetActive(true)
	// image.CenterYAnchor().ConstraintEqualToAnchor(rowView.CenterYAnchor()).SetActive(true)
	// text.LeadingAnchor().ConstraintEqualToAnchorConstant(image.TrailingAnchor(), 3).SetActive(true)
	// text.TrailingAnchor().ConstraintEqualToAnchor(rowView.TrailingAnchor()).SetActive(true)
	// text.CenterYAnchor().ConstraintEqualToAnchor(rowView.CenterYAnchor()).SetActive(true)
	// layout.AliginCenterY(onoff, rowView)
	// layout.PinAnchorTo(onoff.LeadingAnchor(), text.TrailingAnchor(), 5)
	// layout.PinAnchorTo(onoff.TrailingAnchor(), rowView.TrailingAnchor(), 1)
	return rowView.View
}

func (s *Sidebar) SetSidebarMaxWidth() {
	if !s.max.IsNil() {
		s.max.SetActive(false)
	}
	s.max = s.view.WidthAnchor().ConstraintLessThanOrEqualToConstant(s.w.Frame().Size.Width / 2)
	s.max.SetActive(true)
}
