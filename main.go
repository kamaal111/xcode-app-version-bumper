package main

import (
	"flag"

	bumpVersion "github.com/kamaal111/xcode-app-version-bumper/bumpVersion"
)

const (
	DEFAULT_STRING_FLAG      = ""
	DEFAULT_INT_FLAG         = 0
	XCODE_BUILD_NUMBER_KEY   = "CURRENT_PROJECT_VERSION"
	XCODE_VERSION_NUMBER_KEY = "MARKETING_VERSION"
)

func main() {
	buildNumber := flag.Int("number", DEFAULT_INT_FLAG, "build number")
	versionNumber := flag.String("version", DEFAULT_STRING_FLAG, "version number")
	projectPath := flag.String("project", DEFAULT_STRING_FLAG, "xcode project path")
	flag.Parse()

	bumpVersion.BumpVersion(versionNumber, buildNumber, projectPath)
}
