package comm

import (
	"archive/tar"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
)

func Tar(dirPath, password, tarFileName string) error {
	tarFile, err := os.Create(tarFileName)
	if err != nil {
		return err
	}
	defer tarFile.Close()

	var writer io.Writer = tarFile

	if len(password) != 0 {
		iv := MakeBytes(aes.BlockSize, 0x1f)

		pwd := md5.Sum([]byte(password))
		key := MakeBytes(32, 0xff)
		copy(key, pwd[:])

		block, err := aes.NewCipher(key)
		if err != nil {
			return err
		}
		stream := cipher.NewCTR(block, iv)

		writer = cipher.StreamWriter{
			S: stream,
			W: tarFile,
		}
	}

	gzWriter := gzip.NewWriter(writer)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		linkName := ""
		if info.Mode()&fs.ModeSymlink != 0 {
			linkName, err = os.Readlink(path)
			if err != nil {
				return err
			}
		}

		header, err := tar.FileInfoHeader(info, linkName)
		if err != nil {
			return err
		}

		header.Name = relPath

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		//忽略掉非普通文件，如：目录，软连接
		if !info.Mode().IsRegular() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(tarWriter, file)
		return err
	})
}

func UnTar(tarFileName, password, dirPath string) error {
	file, err := os.Open(tarFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	var reader io.Reader = file

	if len(password) != 0 {
		iv := MakeBytes(aes.BlockSize, 0x1f)

		pwd := md5.Sum([]byte(password))
		key := MakeBytes(32, 0xff)
		copy(key, pwd[:])

		block, err := aes.NewCipher(key)
		if err != nil {
			return err
		}
		stream := cipher.NewCTR(block, iv)

		reader = cipher.StreamReader{
			S: stream,
			R: file,
		}
	}

	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		err = extractFile(tarReader, header, filepath.Join(dirPath, header.Name))
		if err != nil {
			return err
		}
	}

	return nil
}

func extractFile(reader *tar.Reader, header *tar.Header, targetPath string) error {
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	switch header.Typeflag {
	case tar.TypeDir:
		return os.MkdirAll(targetPath, os.FileMode(header.Mode))
	case tar.TypeReg:
		return createRegularFile(reader, targetPath, header)
	case tar.TypeSymlink:
		return os.Symlink(header.Linkname, targetPath)
	case tar.TypeLink:
		return createHardLink(targetPath, filepath.Join(filepath.Dir(targetPath), header.Linkname))
	case tar.TypeChar, tar.TypeBlock, tar.TypeFifo:
		return createDeviceFile(targetPath, header)
	default:
		return nil
	}
}

func createRegularFile(reader io.Reader, path string, header *tar.Header) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.CopyN(file, reader, header.Size); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			return err
		}
	}

	if err := os.Chtimes(path, header.AccessTime, header.ModTime); err != nil {
		return err
	}

	return nil
}

func createHardLink(linkPath, target string) error {
	if _, err := os.Stat(target); os.IsNotExist(err) {
		return err
	}

	return os.Link(target, linkPath)
}

func createDeviceFile(path string, header *tar.Header) error {
	mode := os.FileMode(header.Mode)
	dev := uint32(header.Devmajor)<<8 | uint32(header.Devminor)

	return syscall.Mknod(path, uint32(mode), int(dev))
}
