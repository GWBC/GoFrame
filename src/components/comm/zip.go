package comm

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"

	zip "github.com/mzky/zip"
)

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
	tPath, _ := filepath.Abs(dirPath)

	fz, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer fz.Close()
	zw := zip.NewWriter(fz)
	defer zw.Close()

	err = filepath.Walk(tPath, func(path string, info os.FileInfo, err error) error {
		fr, errA := os.Open(path)
		if errA != nil {
			return errA
		}
		defer fr.Close()

		path = path[len(tPath):]

		if info.IsDir() {
			path += "/"
		}

		// 写入文件的头信息
		var w io.Writer
		var errB error
		if password != "" {
			w, errB = zw.Encrypt(path, password, zip.AES256Encryption)
		} else {
			w, errB = zw.Create(path)
		}

		if errB != nil {
			return errB
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
		return errors.New("zip file format error")
	}

	r, err := zip.OpenReader(zipFileName)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.IsEncrypted() {
			f.SetPassword(password)
		}

		fp := filepath.Join(dirPath, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fp+"/", os.ModePerm)
			continue
		}

		w, err := os.Create(fp)
		if err != nil {
			return err
		}
		defer w.Close()

		fr, err := f.Open()
		if err != nil {
			return err
		}
		defer fr.Close()

		if _, errC := io.Copy(w, fr); errC != nil {
			return errC
		}
	}

	return nil
}
