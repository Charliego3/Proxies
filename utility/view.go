package utility

import (
	"github.com/progrium/macdriver/helper/layout"
	"github.com/progrium/macdriver/macos/appkit"
)

func Controller(view appkit.IView) appkit.ViewController {
	controller := appkit.NewViewController()
	controller.SetView(view)
	return controller
}

func Active(constraints ...appkit.LayoutConstraint) {
	for _, constraint := range constraints {
		constraint.SetActive(true)
	}
}

type SeparatorOption struct {
	Super  appkit.IView
	Color  appkit.Color
	Dark   appkit.Color
	Light  appkit.Color
	Width  float64
	Height float64
}

func SeparatorLine(opt SeparatorOption) appkit.Box {
	box := appkit.NewBox()
	box.SetBoxType(appkit.BoxCustom)
	box.SetBorderWidth(0)
	box.SetTranslatesAutoresizingMaskIntoConstraints(false)
	if opt.Color.IsNil() {
		dark, light := opt.Dark, opt.Light
		if light.IsNil() {
			light = ColorHex("#E2E2E2")
		}
		if dark.IsNil() {
			dark = ColorHex("#000000")
		}
		AddAppearanceObserver(func() {
			box.SetFillColor(ColorWithAppearance(light, dark))
		})
	} else {
		box.SetFillColor(opt.Color)
	}
	if !opt.Super.IsNil() {
		opt.Super.AddSubview(box)
	}
	if opt.Width > 0 {
		layout.SetWidth(box, opt.Width)
	}
	if opt.Height > 0 {
		layout.SetHeight(box, opt.Height)
	}
	return box
}

func SymbolButton(symbol string, super appkit.IView) appkit.Button {
	config := appkit.ImageSymbolConfiguration_ConfigurationWithScale(appkit.ImageSymbolScaleLarge)
	button := appkit.NewButtonWithImage(SymbolImage(symbol, config))
	button.SetTranslatesAutoresizingMaskIntoConstraints(false)
	button.SetBordered(false)
	button.SetFocusRingType(appkit.FocusRingTypeNone)
	if super != nil && !super.IsNil() {
		super.AddSubview(button)
	}
	return button
}

func TableColumn(identifier appkit.UserInterfaceItemIdentifier, title string) appkit.TableColumn {
	column := appkit.NewTableColumn()
	column.SetIdentifier(identifier)
	column.SetTitle(title)
	return column
}

// Alert Options

type alertOpts struct {
	style           appkit.AlertStyle
	title           string
	message         string
	showSuppression bool
	accessoryView   appkit.IView
	icon            appkit.IImage
	showHelp        bool
	helpAnchor      appkit.HelpAnchorName
	buttons         []string
	onhelp          func(appkit.Alert) bool
}

type AlertOption func(*alertOpts)

func WithAlertStyle(style appkit.AlertStyle) AlertOption {
	return func(ao *alertOpts) {
		ao.style = style
	}
}

func WithAlertTitle(title string) AlertOption {
	return func(ao *alertOpts) {
		ao.title = title
	}
}

func WithAlertMessage(message string) AlertOption {
	return func(ao *alertOpts) {
		ao.message = message
	}
}

func WithAlertShowSuppression(show bool) AlertOption {
	return func(ao *alertOpts) {
		ao.showSuppression = show
	}
}

func WithAlertAccessoryView(view appkit.IView) AlertOption {
	return func(ao *alertOpts) {
		ao.accessoryView = view
	}
}

func WithAlertIcon(icon appkit.IImage) AlertOption {
	return func(ao *alertOpts) {
		ao.icon = icon
	}
}

func WithAlertShowHelp(show bool) AlertOption {
	return func(ao *alertOpts) {
		ao.showHelp = show
	}
}

func WithAlertHelpAnchor(anchor appkit.HelpAnchorName) AlertOption {
	return func(ao *alertOpts) {
		ao.helpAnchor = anchor
	}
}

func WithAlertOnHelpClicked(clicked func(appkit.Alert) bool) AlertOption {
	return func(ao *alertOpts) {
		ao.onhelp = clicked
	}
}

func WithAlertButtons(buttons ...string) AlertOption {
	return func(ao *alertOpts) {
		ao.buttons = buttons
	}
}

func ShowAlert(opts ...AlertOption) appkit.ModalResponse {
	opt := new(alertOpts)
	opt.style = appkit.AlertStyleInformational
	for _, f := range opts {
		f(opt)
	}
	dialog := appkit.NewAlert()
	dialog.SetAlertStyle(opt.style)
	dialog.SetMessageText(opt.title)
	dialog.SetInformativeText(opt.message)
	dialog.SetShowsSuppressionButton(opt.showSuppression)
	dialog.SetAccessoryView(opt.accessoryView)
	dialog.SetIcon(opt.icon)
	dialog.SetShowsHelp(opt.showHelp)
	dialog.SetHelpAnchor(opt.helpAnchor)
	if opt.onhelp != nil {
		delegate := new(appkit.AlertDelegate)
		delegate.SetAlertShowHelp(opt.onhelp)
		dialog.SetDelegate(delegate)
	}
	if len(opt.buttons) == 0 {
		opt.buttons = append(opt.buttons, "OK")
	}
	for _, title := range opt.buttons {
		dialog.AddButtonWithTitle(title)
	}
	return dialog.RunModal()
}
