package version

import (
	"fmt"
	"runtime"
)

var (
	Build     = ""
	Branch    = ""
	Version   = ""
	BuildTime = ""
)

func GetVersion() string {
	return fmt.Sprintf("%s-%s-%s-%v-%v, build time:%s.", Branch, Build, Version, runtime.GOOS, runtime.GOARCH, BuildTime)
}
