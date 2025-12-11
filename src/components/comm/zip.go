package comm

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	zip "github.com/mzky/zip"
)

//不支持软连接

func IsZip(fileName string) bool {
	f, err := os.Open(fileName)
	if err != nil {
		return false
	}
	defer f.Close()

	buf := make([]byte, 4)
	if n, err := f.Read(buf); err != nil || n < 4 {
		return false
	}

	return bytes.Equal(buf, []byte("PK\x03\x04"))
}

func Zip(dirPath, password, zipFileName string) error {
	tPath, err := filepath.Abs(dirPath)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(tPath, string(filepath.Separator)) {
		tPath += string(filepath.Separator)
	}

	fz, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer fz.Close()
	zw := zip.NewWriter(fz)
	defer zw.Close()

	err = filepath.Walk(tPath, func(path string, info os.FileInfo, err error) error {
		fr, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fr.Close()

		path, err = filepath.Rel(tPath, path)
		if err != nil {
			return err
		}

		if path == "." {
			return nil
		}

		if info.IsDir() {
			path += "/"
		}

		path = filepath.ToSlash(path)

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = path
		header.Flags = 0x800
		if len(password) != 0 {
			header.Method = zip.Deflate
			header.SetPassword(password)
			header.SetEncryptionMethod(zip.AES256Encryption)
		}

		w, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// 写入文件内容
		_, errC := io.Copy(w, fr)
		if errC != nil {
			return errC
		}

		return nil
	})

	if err != nil {
		return err
	}

	return zw.Flush()
}

func UnZip(zipFileName, password, dirPath string) error {
	if !IsZip(zipFileName) {
		return errors.New("无效的ZIP文件")
	}

	r, err := zip.OpenReader(zipFileName)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.IsEncrypted() {
			if len(password) == 0 {
				continue
			}

			f.SetPassword(password)
		}

		if f.FileInfo().IsDir() {
			continue
		}

		fr, err := f.Open()
		if err != nil {
			return err
		}
		defer fr.Close()

		fpath := filepath.Join(dirPath, f.Name)
		os.MkdirAll(filepath.Dir(fpath), os.ModePerm)

		w, err := os.Create(fpath)
		if err != nil {
			return err
		}
		defer w.Close()

		_, err = io.Copy(w, fr)
		if err != nil {
			return err
		}
	}

	return nil
}
