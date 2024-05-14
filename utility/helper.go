package utility

import (
	"path/filepath"

	"github.com/progrium/macdriver/macos/foundation"
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
