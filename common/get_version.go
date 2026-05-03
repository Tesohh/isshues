package common

import "runtime/debug"

// GetVersion returns this build's version.
// if it is impossible to read it, it will return "unknown".
// if this is a development build, it will say "(devel)"
func GetVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	return info.Main.Version
}
