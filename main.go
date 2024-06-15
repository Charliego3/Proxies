package main

import (
	"runtime"

	"github.com/charliego3/proxies/utility"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
	"github.com/progrium/macdriver/objc"
)

var windowFrame = utility.SizeOf(600, 500)

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
	Rules   *Rules
}

var MainWindow ProxyWindow

func (app *App) launching(foundation.Notification) {
	LoadProxies()
	MainWindow = ProxyWindow{}
	MainWindow.Window = appkit.NewWindowWithSizeAndStyle(
		windowFrame.Width, windowFrame.Height,
		appkit.WindowStyleMaskTitled|
			appkit.WindowStyleMaskClosable|
			appkit.WindowStyleMaskMiniaturizable|
			appkit.WindowStyleMaskResizable|
			appkit.WindowStyleMaskFullSizeContentView,
	)
	objc.Retain(&MainWindow)

	MainWindow.Sidebar = NewSidebarController()
	MainWindow.Rules = NewRulesViewController()
	controller := appkit.NewSplitViewController()
	controller.SetSplitViewItems([]appkit.ISplitViewItem{
		appkit.SplitViewItem_SidebarWithViewController(MainWindow.Sidebar),
		appkit.SplitViewItem_SplitViewItemWithViewController(MainWindow.Rules),
	})

	delegate := new(appkit.WindowDelegate)
	delegate.SetWindowDidEndLiveResize(func(notification foundation.Notification) {
		MainWindow.Sidebar.SetSidebarMaxWidth()
	})

	MainWindow.SetDelegate(delegate)
	MainWindow.SetTitlebarSeparatorStyle(appkit.TitlebarSeparatorStyleShadow)
	MainWindow.SetToolbarStyle(appkit.WindowToolbarStyleUnifiedCompact)
	MainWindow.SetTitlebarAppearsTransparent(false)
	MainWindow.SetToolbar(createToolbar(controller))
	MainWindow.SetContentViewController(controller)
	MainWindow.SetContentSize(windowFrame)
	MainWindow.SetContentMinSize(windowFrame)
	MainWindow.Center()
	MainWindow.MakeKeyAndOrderFront(nil)
	app.SetActivationPolicy(appkit.ApplicationActivationPolicyRegular)
	app.ActivateIgnoringOtherApps(true)

	if len(FetchProxies()) == 0 {
		MainWindow.SetTitle("Proxies")
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
