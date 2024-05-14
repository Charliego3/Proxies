package main

import (
	"runtime"

	"github.com/charliego3/proxies/utility"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
	"github.com/progrium/macdriver/objc"
)

var windowFrame = utility.SizeOf(700, 500)

type App struct {
	appkit.Application
}

func newApp() *App {
	app := appkit.Application_SharedApplication()
	return &App{app}
}

func (app *App) launching(foundation.Notification) {
	w := appkit.NewWindowWithSizeAndStyle(
		windowFrame.Width, windowFrame.Height,
		appkit.WindowStyleMaskTitled|
			appkit.WindowStyleMaskClosable|
			appkit.WindowStyleMaskMiniaturizable|
			appkit.WindowStyleMaskResizable|
			appkit.WindowStyleMaskFullSizeContentView,
	)
	objc.Retain(&w)

	sidebarController := NewSidebarController(w)
	controller := appkit.NewSplitViewController()
	controller.SetSplitViewItems([]appkit.ISplitViewItem{
		appkit.SplitViewItem_SidebarWithViewController(sidebarController),
		appkit.SplitViewItem_SplitViewItemWithViewController(
			utility.Controller(
				appkit.NewViewWithFrame(
					utility.RectOf(utility.SizeOf(
						340, 400))))),
	})

	delegate := new(appkit.WindowDelegate)
	delegate.SetWindowDidEndLiveResize(func(notification foundation.Notification) {
		sidebarController.SetSidebarMaxWidth()
	})

	w.SetDelegate(delegate)
	w.SetTitlebarSeparatorStyle(appkit.TitlebarSeparatorStyleShadow)
	w.SetToolbarStyle(appkit.WindowToolbarStyleUnifiedCompact)
	w.SetTitlebarAppearsTransparent(false)
	w.SetToolbar(createToolbar(w, controller))
	w.SetContentViewController(controller)
	w.SetContentSize(windowFrame)
	w.SetContentMinSize(windowFrame)
	w.SetTitle("Proxy Tools")
	w.Center()
	w.MakeKeyAndOrderFront(nil)
	app.SetActivationPolicy(appkit.ApplicationActivationPolicyRegular)
	app.ActivateIgnoringOtherApps(true)
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
