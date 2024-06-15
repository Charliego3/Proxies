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

const (
	colRuleIdentifier   = appkit.UserInterfaceItemIdentifier("ColumnRule")
	colStateIdentifier  = appkit.UserInterfaceItemIdentifier("ColumnRemark")
	colRemarkIdentifier = appkit.UserInterfaceItemIdentifier("ColumnState")
)

type Rules struct {
	appkit.IViewController
	Table appkit.TableView
	ditem appkit.IMenuItem
	titem appkit.IMenuItem
}

func NewRulesViewController() *Rules {
	return new(Rules).Init()
}

var (
	colRule   []foundation.String
	colRemark []foundation.String
)

func init() {
	for i := 65; i < 65+2; i++ {
		colRule = append(colRule, foundation.String_StringWithString(string(rune(i))))
		colRemark = append(colRemark, foundation.String_StringWithString(fmt.Sprintf("%d", i)))
	}
}

func (r *Rules) Init() *Rules {
	r.Table = appkit.NewTableView()
	r.Table.SetColumnAutoresizingStyle(appkit.TableViewSequentialColumnAutoresizingStyle)
	r.Table.SetUsesAlternatingRowBackgroundColors(true)
	r.Table.SetStyle(appkit.TableViewStyleAutomatic)
	r.Table.SetSelectionHighlightStyle(appkit.TableViewSelectionHighlightStyleRegular)
	r.Table.SetUsesSingleLineMode(true)
	r.Table.SetAllowsColumnSelection(false)
	r.Table.SetAutoresizingMask(appkit.ViewWidthSizable)
	r.Table.SetTranslatesAutoresizingMaskIntoConstraints(false)
	r.Table.SetRowHeight(25)
	r.Table.AddTableColumn(createTableColumn(colRuleIdentifier, "Rule"))
	r.Table.AddTableColumn(createTableColumn(colRemarkIdentifier, "Status"))
	r.Table.AddTableColumn(createTableColumn(colStateIdentifier, "Remark"))

	menu := appkit.NewMenu()
	menu.AddItem(utility.MenuItem("New Rule", "plus.circle.fill", func(objc.Object) {
		colRule = append(colRule, foundation.String_StringWithString(""))
		colRemark = append(colRemark, foundation.String_StringWithString(""))
		last := len(colRule) - 1
		r.Table.ReloadData()
		r.Table.ScrollRowToVisible(last)
		r.Table.SelectRowIndexesByExtendingSelection(foundation.NewIndexSetWithIndex(uint(last)), true)
		row := r.Table.RowViewAtRowMakeIfNecessary(last, false)
		fmt.Printf("%+v\n", row)
		if row.IsNil() {
			return
		}
		// fmt.Println(row.Subviews())
		// fmt.Printf("%+v\n", row.Subviews()[0])
		// t := appkit.TextFieldFrom(row.Subviews()[0].Ptr())
		// fmt.Println(t.StringValue())
		t := appkit.TextFieldFrom(row.Subviews()[0].Subviews()[0].Ptr())
		fmt.Println(t.StringValue())
		t.SetPlaceholderString("new rule regex")
		// t.CurrentEditor().KeyDown(appkit.NewEvent())
	}))
	r.titem = utility.MenuItem("Disable", "switch.2", func(objc.Object) {

	})
	r.ditem = utility.MenuItem("Delete", "trash.fill", func(objc.Object) {

	})
	menu.AddItem(r.titem)
	menu.AddItem(r.ditem)
	md := new(appkit.MenuDelegate)
	md.SetMenuWillOpen(func(menu appkit.Menu) {
		hidden := r.Table.ClickedRow() == -1
		if !hidden {
			r.titem.SetTitle("Enable")
		}
		r.titem.SetHidden(hidden)
		r.ditem.SetHidden(hidden)
	})
	menu.SetDelegate(md)
	r.Table.SetMenu(menu)

	datasource := &lib.TableViewDataSourceDelegate{}
	datasource.SetNumberOfRowsInTableView(func(table appkit.TableView) int {
		return len(colRule)
	})
	r.Table.SetDataSource(datasource)

	delegate := &appkit.TableViewDelegate{}
	delegate.SetTableViewViewForTableColumnRow(func(_ appkit.TableView, column appkit.TableColumn, row int) appkit.View {
		if column.Identifier() == colRemarkIdentifier {
			// multiply.circle.fill
			return appkit.ImageView_ImageViewWithImage(
				utility.SymbolImage("checkmark.circle.fill",
					appkit.ImageSymbolConfiguration_ConfigurationWithPaletteColors(
						[]appkit.IColor{appkit.Color_WhiteColor(), utility.ColorHex("00CC00")},
					))).View
		}
		return getRowView(func() string {
			if column.Identifier() == colRuleIdentifier {
				return colRule[row].String()
			}
			return colRemark[row].String()
		})
	})
	r.Table.SetDelegate(delegate)
	r.Table.SelectRowIndexesByExtendingSelection(foundation.NewIndexSetWithIndex(0), true)

	scrollView := appkit.NewScrollView()
	scrollView.SetBorderType(appkit.NoBorder)
	scrollView.SetScrollerKnobStyle(appkit.ScrollerKnobStyleDefault)
	scrollView.SetScrollerStyle(appkit.ScrollerStyleOverlay)
	scrollView.SetFindBarPosition(appkit.ScrollViewFindBarPositionAboveContent)
	scrollView.SetAutohidesScrollers(true)
	scrollView.SetDocumentView(r.Table)
	scrollView.SetHasVerticalScroller(true)
	scrollView.SetTranslatesAutoresizingMaskIntoConstraints(false)

	layout.SetMinWidth(scrollView, 300)
	controller := appkit.NewViewController()
	controller.SetView(scrollView)
	r.IViewController = controller
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

func getRowView(supplier func() string) appkit.View {
	text := appkit.NewTextField()
	text.SetBordered(false)
	text.SetBezelStyle(appkit.TextFieldSquareBezel)
	text.SetDrawsBackground(false)
	text.SetTranslatesAutoresizingMaskIntoConstraints(false)
	text.SetStringValue(supplier())
	text.SetLineBreakStrategy(appkit.LineBreakStrategy(appkit.LineBreakByTruncatingTail))
	text.SetUsesSingleLineMode(true)
	delegate := new(appkit.TextFieldDelegate)
	delegate.SetControlTextDidEndEditing(func(obj foundation.Notification) {
		fmt.Println("text field value did edited:", text.StringValue())
	})
	text.SetDelegate(delegate)

	rowView := appkit.NewView()
	rowView.AddSubview(text)

	text.LeadingAnchor().ConstraintEqualToAnchor(rowView.LeadingAnchor()).SetActive(true)
	text.TrailingAnchor().ConstraintEqualToAnchor(rowView.TrailingAnchor()).SetActive(true)
	text.CenterYAnchor().ConstraintEqualToAnchor(rowView.CenterYAnchor()).SetActive(true)
	return rowView
}
