package timer

import (
	"GoFrame/src/components/comm"
	"GoFrame/src/components/config"
	"GoFrame/src/components/log"
	"os"
	"path/filepath"
	"time"
)

type UploadFTPTimer struct {
	packPath string
}

func (s *UploadFTPTimer) Enable() bool {
	return len(config.Instance.FTPInfo.Addr) != 0
}

func (s *UploadFTPTimer) Init() error {
	s.packPath = filepath.Join(comm.Pwd(), "pack")
	return nil
}

func (s *UploadFTPTimer) Uninit() {
}

func (s *UploadFTPTimer) Name() string {
	return "上传文件"
}

func (s *UploadFTPTimer) Proc() time.Duration {
	for {
		fs, err := os.ReadDir(s.packPath)
		if err != nil {
			log.Sys.Errorf("获取打包目录下文件失败，原因：%s", err.Error())
			break
		}

		if len(fs) == 0 {
			break
		}

		ftp := comm.FTP{Addr: config.Instance.FTPInfo.Addr,
			User:     config.Instance.FTPInfo.User,
			Password: config.Instance.FTPInfo.Password}

		for _, file := range fs {
			if file.IsDir() {
				continue
			}

			upFile := filepath.Join(s.packPath, file.Name())

			for range 3 {
				err = ftp.UpLoad(upFile, config.Instance.FTPInfo.RootPath)
				if err == nil {
					break
				}
			}

			if err != nil {
				log.Sys.Errorf("上传打包文件失败，原因：%s", err.Error())
				return 10 * time.Second
			}

			err := os.Remove(upFile)
			if err != nil {
				log.Sys.Errorf("删除打包文件失败，原因：%s", err.Error())
				return 10 * time.Second
			}
		}
	}

	return 30 * time.Second
}
