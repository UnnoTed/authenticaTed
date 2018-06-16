package utils

import (
	"path/filepath"
	"strings"
)

func GetCurrentPackage(b string) string {
	basepath := filepath.Dir(b)
	idx := strings.Index(basepath, "GoPath/src/")
	return strings.Replace(basepath[idx:], "GoPath/src/", "", -1)
}
