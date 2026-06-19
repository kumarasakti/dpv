package main

import (
	"runtime/debug"

	"github.com/kumarasakti/dpv/cmd"
)

// version is set by -ldflags at build time (make build / CI release).
// When installed via "go install", it stays empty and we fall back to
// the module version embedded by the Go toolchain in the binary.
var version = ""

func main() {
	cmd.Execute(resolveVersion())
}

func resolveVersion() string {
	if version != "" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "dev"
}
