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

type UnPackeFileTimer struct {
	downPath string
}

func (u *UnPackeFileTimer) Enable() bool {
	return len(config.Instance.FTPInfo.Addr) != 0 && len(config.Instance.DownLoad.Path) != 0
}

func (u *UnPackeFileTimer) Init() error {
	u.downPath = filepath.Join(comm.Pwd(), "down")
	return nil
}

func (u *UnPackeFileTimer) Uninit() {
}

func (u *UnPackeFileTimer) Name() string {
	return "解压文件"
}

func (u *UnPackeFileTimer) Proc(ctx context.Context) time.Duration {
	for {
		select {
		case <-ctx.Done():
			return config.ProcInterval(1)
		default:
			fs, err := os.ReadDir(u.downPath)
			if err != nil {
				log.Sys.Errorf("获取下载目录中的文件失败，原因：%s", err.Error())
				return config.ProcInterval(1)
			}

			if len(fs) == 0 {
				return config.ProcInterval(1)
			}

			for _, file := range fs {
				if file.IsDir() {
					continue
				}

				zipPath := filepath.Join(u.downPath, file.Name())
				err := comm.UnTar(zipPath, config.PackPassword, config.Instance.DownLoad.Path)
				if err != nil {
					log.Sys.Errorf("解压zip文件失败：%s，原因：%s", zipPath, err.Error())
				}

				os.Remove(zipPath)

				log.Sys.Debugf("处理压缩文件：%s", zipPath)
			}
		}
	}
}
