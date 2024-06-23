package utility

import (
	"path/filepath"

	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
	"github.com/progrium/macdriver/objc"
)

const Identifier = "com.charlie.proxies"

func SupportPath(sub ...string) string {
	url := foundation.FileManager_DefaultManager().
		URLForDirectoryInDomainAppropriateForURLCreateError(
			foundation.ApplicationSupportDirectory,
			foundation.UserDomainMask,
			nil,
			true,
			nil,
		)

	return filepath.Join(append([]string{url.Path(), Identifier}, sub...)...)
}

func Ternary[T any](b bool, t, f T) T {
	if b {
		return t
	}
	return f
}

func FontAttribute(size float64) map[foundation.AttributedStringKey]objc.IObject {
	return map[foundation.AttributedStringKey]objc.IObject{
		"NSFont": appkit.Font_MenuFontOfSize(size),
	}
}
