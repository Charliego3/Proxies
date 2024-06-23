package main

import (
	"fmt"
	"unsafe"

	"github.com/charliego3/proxies/utility"
	"github.com/progrium/macdriver/dispatch"
	"github.com/progrium/macdriver/helper/layout"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
	"github.com/progrium/macdriver/objc"
)

const (
	colRuleIdentifier   = appkit.UserInterfaceItemIdentifier("ColumnRule")
	colStateIdentifier  = appkit.UserInterfaceItemIdentifier("ColumnRemark")
	colRemarkIdentifier = appkit.UserInterfaceItemIdentifier("ColumnState")

	ruleDefaultsKey = "proxiesRulesKey"
)

type RuleView struct {
	appkit.IViewController
	appkit.TableView
	ditem appkit.IMenuItem
	titem appkit.IMenuItem
}

func NewRulesViewController() *RuleView {
	return new(RuleView).Init()
}

func (r *RuleView) Init() *RuleView {
	r.TableView = appkit.NewTableView()
	r.SetColumnAutoresizingStyle(appkit.TableViewSequentialColumnAutoresizingStyle)
	r.SetUsesAlternatingRowBackgroundColors(true)
	r.SetStyle(appkit.TableViewStyleAutomatic)
	r.SetSelectionHighlightStyle(appkit.TableViewSelectionHighlightStyleRegular)
	r.SetUsesSingleLineMode(true)
	r.SetAllowsColumnSelection(false)
	r.SetAutoresizingMask(appkit.ViewWidthSizable)
	r.SetTranslatesAutoresizingMaskIntoConstraints(false)
	r.SetRowHeight(25)
	r.AddTableColumn(createTableColumn(colRuleIdentifier, "Rule"))
	r.AddTableColumn(createTableColumn(colRemarkIdentifier, "Status"))
	r.AddTableColumn(createTableColumn(colStateIdentifier, "Remark"))

	menu := appkit.NewMenu()
	menu.AddItem(utility.MenuItem("New Rule", "plus.rectangle", func(objc.Object) {
		rules.Add(Rule{T: true})
		lastIndex := rules.LastIndex()
		r.ReloadData()
		r.ScrollRowToVisible(lastIndex)
		r.SelectRowIndexesByExtendingSelection(foundation.NewIndexSetWithIndex(uint(lastIndex)), true)
		dispatch.MainQueue().DispatchAsync(func() {
			text := appkit.TextFieldFrom(lastptr)
			if !text.IsNil() {
				text.BecomeFirstResponder()
			}
		})
	}, appkit.ImageSymbolConfiguration_ConfigurationPreferringMulticolor()))
	r.titem = utility.MenuItem("Switch State", "bolt", func(objc.Object) {
		rule := rules.ByIndex(r.ClickedRow())
		rule.T = !rule.T
		rules.Update(rule)
		r.ReloadData()
	}, appkit.ImageSymbolConfiguration_ConfigurationWithPaletteColors([]appkit.IColor{
		appkit.Color_SystemGreenColor(),
		appkit.Color_SystemYellowColor(),
	}))
	r.ditem = utility.MenuItem("Delete", "trash", func(objc.Object) {
		clicked := r.ClickedRow()
		server.RemoveRegex(rules.ByIndex(clicked).P)
		rules.Delete(clicked)
		r.ReloadData()
	}, appkit.ImageSymbolConfiguration_ConfigurationPreferringMulticolor())
	menu.AddItem(r.titem)
	menu.AddItem(r.ditem)
	md := new(appkit.MenuDelegate)
	md.SetMenuWillOpen(func(menu appkit.Menu) {
		index := r.ClickedRow()
		hidden := index == -1
		if !hidden {
			rule := rules.ByIndex(index)
			if rule.T {
				r.titem.SetTitle("Disable")
				r.titem.SetImage(utility.SymbolImage("bolt.slash.fill", r.titem.Image().SymbolConfiguration()))
			} else {
				r.titem.SetTitle("Enable")
				r.titem.SetImage(utility.SymbolImage("bolt.fill", r.titem.Image().SymbolConfiguration()))
			}
		}
		r.titem.SetHidden(hidden)
		r.ditem.SetHidden(hidden)
	})
	menu.SetDelegate(md)
	r.TableView.SetMenu(menu)

	datasource := &TableViewDataSourceDelegate{}
	datasource.SetNumberOfRowsInTableView(func(table appkit.TableView) int {
		return rules.LastIndex() + 1
	})
	r.SetDataSource(datasource)

	delegate := &appkit.TableViewDelegate{}
	delegate.SetTableViewViewForTableColumnRow(func(_ appkit.TableView, column appkit.TableColumn, row int) appkit.View {
		rule := rules.ByIndex(row)
		if column.Identifier() == colRemarkIdentifier {
			color := appkit.Color_SystemRedColor()
			if rule.T {
				color = utility.ColorHex("00CC00")
			}
			return appkit.ImageView_ImageViewWithImage(utility.SymbolImage("circle.fill",
				appkit.ImageSymbolConfiguration_ConfigurationWithHierarchicalColor(color))).View
		}

		return r.getRowView(column, row, rule)
	})
	r.SetDelegate(delegate)
	r.SelectRowIndexesByExtendingSelection(foundation.NewIndexSetWithIndex(0), true)

	scrollView := appkit.NewScrollView()
	scrollView.SetBorderType(appkit.NoBorder)
	scrollView.SetScrollerKnobStyle(appkit.ScrollerKnobStyleDefault)
	scrollView.SetScrollerStyle(appkit.ScrollerStyleOverlay)
	scrollView.SetFindBarPosition(appkit.ScrollViewFindBarPositionAboveContent)
	scrollView.SetAutohidesScrollers(true)
	scrollView.SetDocumentView(r.TableView)
	scrollView.SetHasVerticalScroller(true)
	scrollView.SetTranslatesAutoresizingMaskIntoConstraints(false)

	layout.SetMinWidth(scrollView, 300)
	controller := appkit.NewViewController()
	controller.SetView(scrollView)
	r.IViewController = controller

	appkit.Event_AddLocalMonitorForEventsMatchingMaskHandler(appkit.EventMaskKeyDown, func(event appkit.Event) appkit.Event {
		if event.KeyCode() == 53 {
			return appkit.Event{}
		}
		return event
	})
	return r
}

func createTableColumn(identifier appkit.UserInterfaceItemIdentifier, title string) appkit.TableColumn {
	column := appkit.NewTableColumn()
	column.SetIdentifier(identifier)
	column.SetTitle(title)
	if identifier == colRemarkIdentifier {
		column.SetWidth(38)
	}
	return column
}

var lastptr unsafe.Pointer

func (r *RuleView) getRowView(column appkit.TableColumn, row int, rule Rule) appkit.View {
	text := appkit.NewTextField()
	value := rule.R
	if column.Identifier() == colRuleIdentifier {
		value = rule.P
		text.SetPlaceholderString("rule pattern")
	}
	text.SetBordered(false)
	text.SetBezelStyle(appkit.TextFieldSquareBezel)
	text.SetDrawsBackground(false)
	text.SetTranslatesAutoresizingMaskIntoConstraints(false)
	text.SetStringValue(value)
	text.SetLineBreakStrategy(appkit.LineBreakStrategy(appkit.LineBreakByTruncatingTail))
	text.SetUsesSingleLineMode(true)

	delegate := new(appkit.TextFieldDelegate)
	delegate.SetControlTextDidEndEditing(func(obj foundation.Notification) {
		if column.Identifier() == colRuleIdentifier {
			rule.P = text.StringValue()
			if rule.P == "" {
				rules.DeleteById(rule.ID)
				r.ReloadData()
				return
			}
		} else {
			rule.R = text.StringValue()
		}
		rules.Update(rule)
		r.ReloadData()
	})
	delegate.SetControlTextDidBeginEditing(func(obj foundation.Notification) {
		fmt.Println("proxy server will be remove regexp:", text.StringValue())
		server.RemoveRegex(text.StringValue())
	})
	text.SetDelegate(delegate)
	if row == rules.LastIndex() && column.Identifier() == colRuleIdentifier {
		lastptr = text.Ptr()
	}

	rowView := appkit.NewView()
	rowView.AddSubview(text)
	text.LeadingAnchor().ConstraintEqualToAnchor(rowView.LeadingAnchor()).SetActive(true)
	text.TrailingAnchor().ConstraintEqualToAnchor(rowView.TrailingAnchor()).SetActive(true)
	text.CenterYAnchor().ConstraintEqualToAnchor(rowView.CenterYAnchor()).SetActive(true)
	return rowView
}
