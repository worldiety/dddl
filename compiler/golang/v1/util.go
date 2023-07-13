package gen

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var regexModName = regexp.MustCompile(`module\s+(.+)`)

// ModuleName tries to grab the nearest module name from the given directory or returns the empty string.
func ModuleName(dir string) (modname string, moddir string) {
	if dir == "" || dir == "." || dir == "/" || dir == "\\" || dir == ".." {
		return "", ""
	}

	buf, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return ModuleName(filepath.Dir(dir))
	}

	match := regexModName.Find(buf)
	if match == nil {
		return "", ""
	}

	return strings.TrimSpace(strings.TrimPrefix(string(match), "module")), dir
}
