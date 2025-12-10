package comm

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/jlaffaye/ftp"
)

type FTP struct {
	Addr     string
	User     string
	Password string
}

func (f *FTP) FileCount(path string) (int, error) {
	c, err := f.login()
	if err != nil {
		return 0, err
	}
	defer f.close(c)

	names, err := c.NameList(path)
	if err != nil {
		return 0, err
	}

	return len(names), nil
}

func (f *FTP) FileNames(path string) ([]string, error) {
	c, err := f.login()
	if err != nil {
		return nil, err
	}
	defer f.close(c)

	return c.NameList(path)
}

func (f *FTP) FileList(path string) ([]*ftp.Entry, error) {
	c, err := f.login()
	if err != nil {
		return nil, err
	}
	defer f.close(c)

	return c.List(path)
}

func (f *FTP) UpLoad(localPath string, remotePath string) error {
	c, err := f.login()
	if err != nil {
		return err
	}
	defer f.close(c)

	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = c.Type(ftp.TransferTypeBinary)
	if err != nil {
		return err
	}

	if len(FileName(remotePath)) == 0 {
		remotePath = filepath.Join(remotePath, filepath.Base(localPath))
		remotePath = filepath.ToSlash(remotePath)
	}

	return c.Stor(remotePath, file)
}

func (f *FTP) Down(remotePath string, localPath string) error {
	c, err := f.login()
	if err != nil {
		return err
	}
	defer f.close(c)

	name := FileName(localPath)
	if len(name) == 0 {
		localPath = filepath.Join(localPath, filepath.Base(remotePath))
	}

	err = os.MkdirAll(filepath.Dir(localPath), 0755)
	if err != nil {
		return err
	}

	resp, err := c.Retr(remotePath)
	if err != nil {
		return err
	}
	defer resp.Close()

	file, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp)

	return err
}

func (f *FTP) Delete(path string) error {
	c, err := f.login()
	if err != nil {
		return err
	}
	defer f.close(c)

	return c.Delete(path)
}

func (f *FTP) login() (*ftp.ServerConn, error) {
	c, err := ftp.Dial(f.Addr, ftp.DialWithTimeout(20*time.Second))
	if err != nil {
		return nil, err
	}

	err = c.Login(f.User, f.Password)
	if err != nil {
		c.Quit()
		return nil, err
	}

	return c, nil
}

func (f *FTP) close(c *ftp.ServerConn) {
	c.Logout()
	c.Quit()
}
