package main

import (
	"fmt"
	"runtime"

	"github.com/kriechi/git-monitor/cmd"
)

func main() {
	cmd.Execute(buildVersionString())
}

// VersionSuffix is empty for proper releases.
var VersionSuffix = "-dev"

// CurrentVersion represents the current build version.
var CurrentVersion = "0.0.2" + VersionSuffix

// BuildDate is a human-readable timestamp of when this binary was built.
var BuildDate = "2006-01-02T15:04:05Z-0700" // dummy timestamp

func buildVersionString() string {
	version := "v" + CurrentVersion

	date := BuildDate
	if date == "" {
		date = "unknown"
	}

	return fmt.Sprintf(
		"%s goos:%s goarch:%s runtime:%s built:%s",
		version,
		runtime.GOOS,
		runtime.GOARCH,
		runtime.Version(),
		date,
	)
}
