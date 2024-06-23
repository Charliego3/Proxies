package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charliego3/proxies/utility"
	"github.com/progrium/macdriver/helper/layout"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
	"github.com/progrium/macdriver/objc"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	sidebarIdentifier = "sidebarDatasourceIdentifier"

	proxyDefaultsKey = "proxiesKey"
)

type Sidebar struct {
	appkit.IViewController
	view    appkit.IView
	outline appkit.OutlineView
	max     appkit.LayoutConstraint
}

func NewSidebarController() *Sidebar {
	sidebar := new(Sidebar)
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
	menu.AddItem(utility.MenuItem("Toggle State", "togglepower", func(sender objc.Object) {
		proxy := proxies.ByIndex(s.outline.ClickedRow())
		proxy.InUse = !proxy.InUse
		proxies.Update(proxy)
		selected := s.outline.SelectedRow()
		s.Update()
		s.SelectRow(selected)
	}))
	menu.AddItem(utility.MenuItem("Edit Options", "pencil.line", func(sender objc.Object) {
		OpenProxySheet("Update options for Proxy:", proxies.ByIndex(s.outline.ClickedRow()))
	}))
	menu.AddItem(utility.MenuItem("Delete Proxy", "trash", func(objc.Object) {
		index := s.outline.ClickedRow()
		proxy := proxies.ByIndex(index)
		utility.ShowAlert(
			utility.WithAlertTitle("Are you sure delete this proxy?"),
			utility.WithAlertMessage(fmt.Sprintf("%s\n%s", cases.Title(language.English).String(proxy.Name), proxy.URL())),
			utility.WithAlertWindow(Window),
			utility.WithAlertStyle(appkit.AlertStyleWarning),
			utility.WithAlertButtons("OK", "Cancel"),
			utility.WithAlertHandler(func(code appkit.ModalResponse) {
				if code == appkit.AlertFirstButtonReturn {
					proxies.Delete(index)
					s.Update()
					s.SelectRow(s.outline.NumberOfRows() - 1)
				}
			}),
		)
	}))

	s.outline.SetMenu(menu)
	s.setDelegate()
	s.setDatasource()
	s.SelectRow(0)
	scrollView := appkit.NewScrollView()
	clipView := appkit.ClipViewFrom(scrollView.ContentView().Ptr())
	clipView.SetDocumentView(s.outline)
	clipView.SetAutomaticallyAdjustsContentInsets(false)
	clipView.SetContentInsets(foundation.EdgeInsets{Top: 10})

	scrollView.SetBorderType(appkit.NoBorder)
	scrollView.SetDrawsBackground(false)
	scrollView.SetAutohidesScrollers(true)
	scrollView.SetHasVerticalScroller(true)
	scrollView.ContentView().ScrollToPoint(foundation.Point{Y: -10})

	s.view = scrollView
	s.IViewController.SetView(s.view)
	s.view.SetFrameSize(utility.SizeOf(200, 0))
	layout.SetMinWidth(s.view, 200)
	s.SetSidebarMaxWidth()
}

func (s *Sidebar) setDelegate() {
	delegate := new(appkit.OutlineViewDelegate)
	delegate.SetOutlineViewViewForTableColumnItem(s.createColumnItem)
	delegate.SetControlTextDidBeginEditing(func(obj foundation.Notification) {})
	delegate.SetOutlineViewHeightOfRowByItem(func(outlineView appkit.OutlineView, item objc.Object) float64 {
		return 40
	})
	delegate.SetOutlineViewSelectionDidChange(func(notification foundation.Notification) {
		proxy := proxies.ByIndex(s.outline.SelectedRow())
		Window.SetTitle(cases.Title(language.English).String(proxy.Name))
	})
	po0 := objc.WrapAsProtocol("NSOutlineViewDelegate", appkit.POutlineViewDelegate(delegate))
	objc.SetAssociatedObject(s.outline, objc.AssociationKey("setDelegate"), po0, objc.ASSOCIATION_RETAIN)
	objc.Call[objc.Void](s.outline, objc.Sel("setDelegate:"), po0)
}

func (s *Sidebar) setDatasource() {
	datasource := new(OutlineViewDatasource)
	datasource.SetOutlineViewChildOfItem(func(outlineView appkit.OutlineView, index int, item objc.Object) objc.Object {
		proxy := proxies.ByIndex(index)
		return foundation.String_StringWithString(fmt.Sprintf("%s#%s#%t",
			proxy.Name, proxy.URL(), proxy.InUse)).Object
	})
	datasource.SetOutlineViewIsItemExpandable(func(outlineView appkit.OutlineView, item objc.Object) bool {
		return false
	})
	datasource.SetOutlineViewNumberOfChildrenOfItem(func(_ appkit.OutlineView, item objc.Object) int {
		return proxies.Length()
	})
	appkit.NewTableViewRowAction()
	po1 := objc.WrapAsProtocol("NSOutlineViewDataSource", appkit.POutlineViewDataSource(datasource))
	objc.SetAssociatedObject(s.outline, objc.AssociationKey("setDataSource"), po1, objc.ASSOCIATION_RETAIN)
	objc.Call[objc.Void](s.outline, objc.Sel("setDataSource:"), po1)
}

func (s *Sidebar) createColumnItem(_ appkit.OutlineView, tableColumn appkit.TableColumn, item objc.Object) appkit.View {
	content := foundation.StringFrom(item.Ptr())
	arrs := strings.Split(content.String(), "#")
	inUse, _ := strconv.ParseBool(arrs[2])

	symbol := "globe.europe.africa"
	if inUse {
		symbol += ".fill"
	}
	image := appkit.NewImageView()
	image.SetTranslatesAutoresizingMaskIntoConstraints(false)
	image.SetImage(utility.SymbolImage(symbol, utility.ImageLarge))

	text := appkit.NewLabel(cases.Title(language.English).String(arrs[0]))
	text.SetTranslatesAutoresizingMaskIntoConstraints(false)
	url := appkit.NewTextField()
	url.SetBordered(false)
	url.SetEditable(false)
	url.SetDrawsBackground(false)
	url.SetTranslatesAutoresizingMaskIntoConstraints(false)
	url.SetAttributedStringValue(foundation.NewAttributedStringWithStringAttributes(
		strings.ToLower(arrs[1]), map[foundation.AttributedStringKey]objc.IObject{
			"NSColor": appkit.Color_SecondaryLabelColor(),
			"NSFont":  appkit.Font_LabelFontOfSize(10),
		}))

	info := appkit.NewView()
	info.SetTranslatesAutoresizingMaskIntoConstraints(false)
	info.AddSubview(text)
	info.AddSubview(url)
	layout.AliginLeading(text, info)
	layout.AliginTop(text, info)
	layout.AliginTrailing(text, info)
	layout.PinAnchorTo(url.TopAnchor(), text.BottomAnchor(), 0)
	layout.AliginLeading(url, info)
	layout.AliginTrailing(url, info)
	layout.AliginBottom(url, info)

	rowView := appkit.NewTableRowView()
	rowView.AddSubview(image)
	rowView.AddSubview(info)

	layout.AliginLeading(image, rowView)
	layout.AliginCenterY(image, rowView)
	layout.SetWidth(image, 20)
	layout.PinAnchorTo(info.LeadingAnchor(), image.TrailingAnchor(), 5)
	layout.AliginCenterY(info, rowView)
	layout.AliginTrailing(info, rowView)
	return rowView.View
}

func (s *Sidebar) SetSidebarMaxWidth() {
	if !s.max.IsNil() {
		s.max.SetActive(false)
	}
	s.max = s.view.WidthAnchor().ConstraintLessThanOrEqualToConstant(Window.Frame().Size.Width / 2)
	s.max.SetActive(true)
}

func (s *Sidebar) Update() {
	row := s.outline.SelectedRow()
	if row < 0 {
		row = 0
	}
	length := proxies.Length()
	if row >= length {
		row = length - 1
	}
	s.outline.ReloadData()
}

func (s *Sidebar) ScrollToBottom() {
	rows := s.outline.NumberOfRows()
	s.outline.ScrollRowToVisible(s.outline.NumberOfRows() - 1)
	s.SelectRow(rows - 1)
}

func (s *Sidebar) SelectRow(row int) {
	s.outline.SelectRowIndexesByExtendingSelection(foundation.NewIndexSetWithIndex(uint(row)), true)
}

func (s *Sidebar) SelectedIndex() int {
	return s.outline.SelectedRow()
}

func (s *Sidebar) SelectedId() string {
	return proxies.ByIndex(s.SelectedIndex()).ID
}
