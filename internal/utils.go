package internal

import (
	"runtime"
	"strings"
)

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func ToWindowsPath(path string) string {
	return strings.ReplaceAll(path, "/", "\\")
}
