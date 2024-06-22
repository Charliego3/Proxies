package main

import (
	"runtime"

	"github.com/charliego3/proxies/utility"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
	"github.com/progrium/macdriver/objc"
)

var (
	windowFrame = utility.SizeOf(600, 500)
	defaults    = appkit.UserDefaultsController_SharedUserDefaultsController().Defaults()
)

type App struct {
	appkit.Application
}

func newApp() *App {
	app := appkit.Application_SharedApplication()
	return &App{app}
}

type ProxyWindow struct {
	appkit.Window

	Sidebar *Sidebar
	Rules   *RuleView
}

var Window ProxyWindow

func (app *App) launching(foundation.Notification) {
	Window = ProxyWindow{}
	Window.Window = appkit.NewWindowWithSizeAndStyle(
		windowFrame.Width, windowFrame.Height,
		appkit.WindowStyleMaskTitled|
			appkit.WindowStyleMaskClosable|
			appkit.WindowStyleMaskMiniaturizable|
			appkit.WindowStyleMaskResizable|
			appkit.WindowStyleMaskFullSizeContentView,
	)
	objc.Retain(&Window)

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
	app.SetActivationPolicy(appkit.ApplicationActivationPolicyRegular)
	app.ActivateIgnoringOtherApps(true)

	if proxies.Length() == 0 {
		Window.SetTitle("Proxies")
		OpenNewProxySheet()
	}
}

func (app *App) setMainMenu(foundation.Notification) {

}

func (app *App) shouldTerminateAfterClose(appkit.Application) bool {
	return true
}

func (app *App) waillTerminate(foundation.Notification) {

}

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	app := newApp()
	delegate := new(appkit.ApplicationDelegate)
	delegate.SetApplicationDidFinishLaunching(app.launching)
	delegate.SetApplicationWillFinishLaunching(app.setMainMenu)
	delegate.SetApplicationShouldTerminateAfterLastWindowClosed(app.shouldTerminateAfterClose)
	delegate.SetApplicationWillTerminate(app.waillTerminate)
	app.SetDelegate(delegate)
	app.Run()
}
