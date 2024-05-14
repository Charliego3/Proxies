package utility

import "github.com/progrium/macdriver/macos/appkit"

func SymbolImage(name string, cfg ...appkit.ImageSymbolConfiguration) appkit.Image {
	image := appkit.Image_ImageWithSystemSymbolNameAccessibilityDescription(name, name)
	if len(cfg) > 0 {
		image = image.ImageWithSymbolConfiguration(cfg[0])
	}
	return image
}

var ImageLarge = appkit.ImageSymbolConfiguration_ConfigurationWithScale(appkit.ImageSymbolScaleLarge)
