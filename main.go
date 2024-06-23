package main

import (
	"runtime"

	"github.com/charliego3/proxies/utility"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
	"github.com/progrium/macdriver/objc"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	windowFrame = utility.SizeOf(600, 500)
	defaults    = appkit.UserDefaultsController_SharedUserDefaultsController().Defaults()
)

type App struct {
	appkit.Application

	status appkit.StatusItem
}

func newApp() *App {
	app := appkit.Application_SharedApplication()
	return &App{Application: app}
}

type ProxyWindow struct {
	appkit.Window

	Sidebar *Sidebar
	Rules   *RuleView
}

var (
	app    *App
	Window *ProxyWindow
)

func (app *App) launching(foundation.Notification) {
	app.setupSystemBar()
	app.SetActivationPolicy(appkit.ApplicationActivationPolicyRegular)

	if proxies.Length() == 0 {
		app.launchWindow(objc.Object{})
		OpenNewProxySheet()
	}
}

func (app *App) launchWindow(_ objc.Object) {
	app.ActivateIgnoringOtherApps(true)
	if Window != nil {
		Window.OrderFrontRegardless()
		return
	}

	Window = new(ProxyWindow)
	Window.Window = appkit.NewWindowWithSizeAndStyle(
		windowFrame.Width, windowFrame.Height,
		appkit.WindowStyleMaskTitled|
			appkit.WindowStyleMaskClosable|
			appkit.WindowStyleMaskMiniaturizable|
			appkit.WindowStyleMaskResizable|
			appkit.WindowStyleMaskFullSizeContentView,
	)
	Window.SetTitle("Proxies")
	objc.Retain(Window)

	Window.Sidebar = NewSidebarController()
	Window.Rules = NewRulesViewController()
	controller := appkit.NewSplitViewController()
	controller.SetSplitViewItems([]appkit.ISplitViewItem{
		appkit.SplitViewItem_SidebarWithViewController(Window.Sidebar),
		appkit.SplitViewItem_SplitViewItemWithViewController(Window.Rules),
	})

	delegate := new(appkit.WindowDelegate)
	delegate.SetWindowDidEndLiveResize(func(notification foundation.Notification) {
		Window.Sidebar.SetSidebarMaxWidth()
	})
	delegate.SetWindowWillClose(func(notification foundation.Notification) {
		Window = nil
	})

	Window.SetDelegate(delegate)
	Window.SetTitlebarSeparatorStyle(appkit.TitlebarSeparatorStyleShadow)
	Window.SetToolbarStyle(appkit.WindowToolbarStyleUnifiedCompact)
	Window.SetTitlebarAppearsTransparent(false)
	Window.SetToolbar(createToolbar(controller))
	Window.SetContentViewController(controller)
	Window.SetContentSize(windowFrame)
	Window.SetContentMinSize(windowFrame)
	Window.Center()
	Window.MakeKeyAndOrderFront(nil)
}

func (app *App) setupSystemBar() {
	status := appkit.StatusBar_SystemStatusBar().StatusItemWithLength(appkit.VariableStatusItemLength)
	status.Button().SetImage(utility.SymbolImage("globe.asia.australia"+utility.Ternary(proxies.AnyUsing(), ".fill", ""), utility.ImageLarge))
	objc.Retain(&status)
	app.status = status

	var showTitled bool
	delegate := new(appkit.MenuDelegate)
	delegate.SetMenuWillOpen(func(menu appkit.Menu) {
		menu.RemoveAllItems()
		label := appkit.NewMenuItem()
		label.SetAttributedTitle(foundation.NewAttributedStringWithStringAttributes("Proxies", utility.FontAttribute(15)))
		menu.AddItem(label)
		for _, p := range proxies.Fetch() {
			item := appkit.NewMenuItem()
			attributed := foundation.NewMutableAttributedString()
			attributed.AppendAttributedString(foundation.NewAttributedStringWithStringAttributes(
				cases.Title(language.English).String(p.Name),
				utility.FontAttribute(12)),
			)
			attributed.AppendAttributedString(foundation.NewAttributedStringWithStringAttributes("\n"+p.URL(), utility.FontAttribute(9)))
			item.SetAttributedTitle(attributed)
			item.SetImage(utility.SymbolImage("circle.fill", appkit.ImageSymbolConfiguration_ConfigurationWithHierarchicalColor(
				utility.Ternary(p.InUse, utility.ColorHex("#1DAD03"), utility.ColorHex("#871313")))))
			submenu := appkit.NewMenu()
			for _, r := range rules.FetchWithProxy(p.ID) {
				submenu.AddItem(utility.MenuItem(r.P, utility.Ternary(r.T, "circle.circle.fill", "circle.circle"), func(objc.Object) {

				}))
			}
			item.SetSubmenu(submenu)
			menu.AddItem(item)
		}

		menu.AddItem(appkit.MenuItem_SeparatorItem())
		menu.AddItem(utility.MenuItem(utility.Ternary(showTitled, "Hide Title", "Show Title"), "note.text", func(o objc.Object) {
			showTitled = !showTitled
			app.status.Button().SetTitle(utility.Ternary(showTitled, "Proxies", ""))
		}))
		menu.AddItem(utility.MenuItem("Open Window", "text.and.command.macwindow", app.launchWindow))
		menu.AddItem(appkit.MenuItem_SeparatorItem())
		menu.AddItem(utility.MenuItem("Hide", "eye.fill", utility.Ternary(Window == nil, nil, func(objc.Object) { app.Hide(nil) })))
		menu.AddItem(utility.MenuItem("Quit", "dot.circle.and.hand.point.up.left.fill", func(sender objc.Object) { app.Terminate(nil) }))
	})

	menu := appkit.NewMenuWithTitle("main")
	menu.SetDelegate(delegate)
	app.status.SetMenu(menu)
}

func (app *App) setMainMenu(foundation.Notification) {
	menu := appkit.NewMenuWithTitle("main")
	app.SetMainMenu(menu)

	mainMenuItem := appkit.NewMenuItemWithSelector("", "", objc.Selector{})
	mainMenuMenu := appkit.NewMenuWithTitle("App")
	mainMenuMenu.AddItem(appkit.NewMenuItemWithAction("Hide Proxies", "h", func(_ objc.Object) { app.Hide(nil) }))
	mainMenuMenu.AddItem(appkit.NewMenuItemWithAction("Quit Proxies", "q", func(_ objc.Object) { app.Terminate(nil) }))
	mainMenuItem.SetSubmenu(mainMenuMenu)
	menu.AddItem(mainMenuItem)

	editItem := appkit.NewMenuItemWithSelector("", "", objc.Selector{})
	editMenu := appkit.NewMenuWithTitle("Edit")
	editMenu.AddItem(appkit.NewMenuItemWithSelector("Select All", "a", objc.Sel("selectAll:")))
	editMenu.AddItem(appkit.MenuItem_SeparatorItem())
	editMenu.AddItem(appkit.NewMenuItemWithSelector("Copy", "c", objc.Sel("copy:")))
	editMenu.AddItem(appkit.NewMenuItemWithSelector("Paste", "v", objc.Sel("paste:")))
	editMenu.AddItem(appkit.NewMenuItemWithSelector("Cut", "x", objc.Sel("cut:")))
	editMenu.AddItem(appkit.NewMenuItemWithSelector("Undo", "z", objc.Sel("undo:")))
	editMenu.AddItem(appkit.NewMenuItemWithSelector("Redo", "Z", objc.Sel("redo:")))
	editItem.SetSubmenu(editMenu)
	menu.AddItem(editItem)
}

func (app *App) shouldTerminateAfterClose(appkit.Application) bool {
	return false
}

func (app *App) waillTerminate(foundation.Notification) {

}

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	app = newApp()
	delegate := new(appkit.ApplicationDelegate)
	delegate.SetApplicationDidFinishLaunching(app.launching)
	delegate.SetApplicationWillFinishLaunching(app.setMainMenu)
	delegate.SetApplicationShouldTerminateAfterLastWindowClosed(app.shouldTerminateAfterClose)
	delegate.SetApplicationWillTerminate(app.waillTerminate)
	app.SetDelegate(delegate)
	app.Run()
}
