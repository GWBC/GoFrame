package comm

import (
	"os"
	"path/filepath"
)

func FileCount(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, v := range entries {
		if v.IsDir() {
			continue
		}

		count++
	}

	return count, nil
}

func CopySymlink(src, dst string) error {
	os.MkdirAll(filepath.Dir(dst), 0755)

	target, err := os.Readlink(src)
	if err != nil {
		return err
	}

	return os.Symlink(target, dst)
}
