package version

import (
	"fmt"
	"runtime/debug"
	"strconv"
)

var (
	Commit    string
	Dirty     bool
	GoVersion string

	// Printable contains all the version information as a new-line separated list
	Printable string
)

func init() {
	info, ok := debug.ReadBuildInfo()
	if ok {
		GoVersion = info.GoVersion
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				Commit = setting.Value[:7]
			case "vcs.modified":
				Dirty, _ = strconv.ParseBool(setting.Value)
			}
		}
	} else {
		Commit = "unknown"
		Dirty = true
		GoVersion = "unknown"
	}

	Printable = fmt.Sprintf(`commit: %s
dirty:  %t
go:     %s`, Commit, Dirty, GoVersion)
}
