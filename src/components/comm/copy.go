package comm

import (
	"io"
	"os"
	"path/filepath"
	"syscall"
	"time"
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

	n, err := io.Copy(destination, source)
	if err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			return 0, err
		}
	}

	return n, nil
}

func CopyFileAndMeta(src string, dst string) (int64, error) {
	err := os.MkdirAll(filepath.Dir(dst), 0755)
	if err != nil {
		return 0, err
	}

	info, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return 0, err
	}
	defer dstFile.Close()

	n, err := io.Copy(dstFile, srcFile)
	if err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			return 0, err
		}
	}

	accessTime := info.ModTime()
	if sys, ok := info.Sys().(*syscall.Stat_t); ok {
		accessTime = time.Unix(int64(sys.Atim.Sec), int64(sys.Atim.Nsec))
	}

	os.Chtimes(dst, accessTime, info.ModTime())

	return n, err
}
