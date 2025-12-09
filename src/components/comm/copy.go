package comm

import (
	"io"
	"os"
	"path/filepath"
)

func CopyFile(src string, dst string) (int64, error) {
	err := os.MkdirAll(filepath.Dir(dst), 0755)
	if err != nil {
		return 0, err
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	return io.Copy(destination, source)
}
