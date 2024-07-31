package vers

import (
	"fmt"
	"runtime/debug"
	"strings"
)

func Version(versionValue string) string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}

	// If nothing is injected at build time, try to read from Main.ComposeVersion.
	if versionValue == "" {
		versionValue = strings.Trim(buildInfo.Main.Version, "()")
	}

	// If Main.ComposeVersion doesn't pick up the module Version, mark it with the vcs info from BuildInfo.
	if versionValue == "" || versionValue == "devel" {
		var rev, time string

		for _, setting := range buildInfo.Settings {
			if setting.Key == "vcs.revision" {
				rev = setting.Value
			} else if setting.Key == "vcs.time" {
				time = setting.Value
			}
		}

		versionValue = fmt.Sprintf("devel-%s@%s", rev, time)
	}

	return versionValue
}
