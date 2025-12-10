package timer

import (
	"GoFrame/src/components/comm"
	"GoFrame/src/components/config"
	"GoFrame/src/components/log"
	"os"
	"path/filepath"
	"time"
)

type UnzipFileTimer struct {
	downPath string
}

func (s *UnzipFileTimer) Enable() bool {
	return len(config.Instance.FTPInfo.Addr) != 0 && len(config.Instance.DownLoad.Path) != 0
}

func (s *UnzipFileTimer) Init() error {
	s.downPath = filepath.Join(comm.Pwd(), "down")
	return nil
}

func (s *UnzipFileTimer) Uninit() {
}

func (s *UnzipFileTimer) Name() string {
	return "解压文件"
}

func (s *UnzipFileTimer) Proc() time.Duration {
	for {
		fs, err := os.ReadDir(s.downPath)
		if err != nil {
			log.Sys.Errorf("获取下载目录中的文件失败，原因：%s", err.Error())
			break
		}

		if len(fs) == 0 {
			break
		}

		for _, file := range fs {
			if file.IsDir() {
				continue
			}

			zipPath := filepath.Join(s.downPath, file.Name())
			err := comm.UnZip(zipPath, config.ZIPPassword, config.Instance.DownLoad.Path)
			if err != nil {
				log.Sys.Errorf("解压zip文件失败：%s，原因：%s", zipPath, err.Error())
			}

			os.Remove(zipPath)

			log.Sys.Debugf("处理压缩文件：%s", zipPath)
		}
	}

	return config.ProcInterval(1)
}
