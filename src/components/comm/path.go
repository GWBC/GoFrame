package comm

import (
	"os"
	"path/filepath"
)

func Pwd() string {
	return filepath.Dir(os.Args[0])
}
