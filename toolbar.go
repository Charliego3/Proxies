package main

import (
	"github.com/charliego3/proxies/utility"
	"github.com/progrium/macdriver/helper/action"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/objc"
)

const (
	ToolbarAddProxyButtonIdentifier appkit.ToolbarItemIdentifier = "AddProxy"
	ToolbarToggleSidebarIdentifier  appkit.ToolbarItemIdentifier = "ToolbarToggleSidebar"
)

type Toolbar struct {
	appkit.Toolbar
	splitViewController appkit.SplitViewController
}

func createToolbar(controller appkit.SplitViewController) *Toolbar {
	toolbar := new(Toolbar)
	toolbar.splitViewController = controller
	toolbar.Toolbar = appkit.NewToolbar()
	toolbar.SetDisplayMode(appkit.ToolbarDisplayModeIconOnly)
	toolbar.SetShowsBaselineSeparator(true)
	toolbar.SetDelegate(toolbar.getToolbarDelegate())
	toolbar.SetAllowsExtensionItems(true)
	return toolbar
}

func (t Toolbar) identifiers(appkit.Toolbar) []appkit.ToolbarItemIdentifier {
	return []appkit.ToolbarItemIdentifier{
		ToolbarToggleSidebarIdentifier,
		appkit.ToolbarFlexibleSpaceItemIdentifier,
		ToolbarAddProxyButtonIdentifier,
		appkit.ToolbarSidebarTrackingSeparatorItemIdentifier,
		appkit.ToolbarFlexibleSpaceItemIdentifier,
		appkit.ToolbarFlexibleSpaceItemIdentifier,
	}
}

func (t Toolbar) createItem(identifier appkit.ToolbarItemIdentifier, symbol string, handler action.Handler) appkit.ToolbarItem {
	cfg := appkit.ImageSymbolConfiguration_ConfigurationWithScale(appkit.ImageSymbolScaleLarge)
	item := appkit.NewToolbarItemWithItemIdentifier(identifier)
	button := appkit.NewButton()
	button.SetImage(utility.SymbolImage(symbol, cfg))
	button.SetButtonType(appkit.ButtonTypeMomentaryPushIn)
	button.SetBezelStyle(appkit.BezelStyleTexturedRounded)
	button.SetFocusRingType(appkit.FocusRingTypeNone)
	action.Set(button, handler)
	item.SetView(button)
	return item
}

func (t Toolbar) getToolbarDelegate() *appkit.ToolbarDelegate {
	delegate := new(appkit.ToolbarDelegate)
	delegate.SetToolbarAllowedItemIdentifiers(t.identifiers)
	delegate.SetToolbarDefaultItemIdentifiers(t.identifiers)
	delegate.SetToolbarItemForItemIdentifierWillBeInsertedIntoToolbar(func(
		_ appkit.Toolbar,
		identifier appkit.ToolbarItemIdentifier,
		_ bool,
	) appkit.ToolbarItem {
		switch identifier {
		case ToolbarToggleSidebarIdentifier:
			return t.createItem(identifier, "sidebar.leading", func(sender objc.Object) {
				t.splitViewController.ToggleSidebar(nil)
			})
		case ToolbarAddProxyButtonIdentifier:
			return t.createItem(identifier, "plus", func(_ objc.Object) {
				OpenNewProxySheet()
			})
		}
		return appkit.ToolbarItem{}
	})
	return delegate
}

func OpenNewProxySheet() {
	OpenProxySheet("Choose options for your new Proxy:", NewProxy())
}
