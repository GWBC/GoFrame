package timer

import (
	"GoFrame/src/components/comm"
	"GoFrame/src/components/config"
	"GoFrame/src/components/log"
	"context"
	"os"
	"path/filepath"
	"time"
)

type UploadFTPTimer struct {
	uploadPath string
}

func (s *UploadFTPTimer) Enable() bool {
	return len(config.Instance.FTPInfo.Addr) != 0
}

func (s *UploadFTPTimer) Init() error {
	s.uploadPath = filepath.Join(comm.Pwd(), "upload")
	return nil
}

func (s *UploadFTPTimer) Uninit() {
}

func (s *UploadFTPTimer) Name() string {
	return "上传文件"
}

func (s *UploadFTPTimer) Proc(ctx context.Context) time.Duration {
	ftp := comm.FTP{Addr: config.Instance.FTPInfo.Addr,
		User:     config.Instance.FTPInfo.User,
		Password: config.Instance.FTPInfo.Password}

	for {
		select {
		case <-ctx.Done():
			return config.ProcInterval(1)
		default:
			fs, err := os.ReadDir(s.uploadPath)
			if err != nil {
				log.Sys.Errorf("获取打包目录下文件失败，原因：%s", err.Error())
				return config.ProcInterval(1)
			}

			if len(fs) == 0 {
				return config.ProcInterval(1)
			}

			for _, file := range fs {
				if file.IsDir() {
					continue
				}

				upFile := filepath.Join(s.uploadPath, file.Name())

				for range 3 {
					err = ftp.UpLoad(upFile, config.Instance.FTPInfo.RootPath)
					if err == nil {
						break
					}
				}

				if err != nil {
					log.Sys.Errorf("上传打包文件失败，原因：%s", err.Error())
					return config.ProcInterval(1)
				}

				err := os.Remove(upFile)
				if err != nil {
					log.Sys.Errorf("删除打包文件失败，原因：%s", err.Error())
					return config.ProcInterval(1)
				}
			}
		}
	}
}
