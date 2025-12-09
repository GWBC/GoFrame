package comm

import (
	"os"
	"path/filepath"
	"strings"
)

func Pwd() string {
	return filepath.Dir(os.Args[0])
}

func FileName(path string) string {
	path = strings.TrimSpace(path)
	if path[len(path)-1] == '\\' || path[len(path)-1] == '/' {
		return ""
	}

	return filepath.Base(path)
}
