package utility

import (
	"strings"

	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
)

func IsDark() bool {
	effected := appkit.Application_SharedApplication().EffectiveAppearance()
	return effected.BestMatchFromAppearancesWithNames([]appkit.AppearanceName{
		appkit.AppearanceNameAqua,
		appkit.AppearanceNameDarkAqua,
	}) == appkit.AppearanceNameDarkAqua
}

func ColorWithAppearance(light, dark appkit.Color) appkit.Color {
	if IsDark() {
		return dark
	}
	return light
}

func ColorWithRGBA(r, g, b, a float64) appkit.Color {
	return appkit.Color_ColorWithRedGreenBlueAlpha(r/255, g/255, b/255, a)
}

func ColorHex(hex string) appkit.Color {
	hex = strings.TrimPrefix(hex, "#")
	scanner := foundation.NewScannerWithString(hex)
	var rgb int64
	if !scanner.ScanHexLongLong(&rgb) {
		return appkit.Color_ClearColor()
	}

	return ColorWithRGBA(
		float64((rgb&0xFF0000)>>16),
		float64((rgb&0x00FF00)>>8),
		float64(rgb&0x0000FF),
		1,
	)
}
