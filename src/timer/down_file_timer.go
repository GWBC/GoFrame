package timer

import (
	"GoFrame/src/components/comm"
	"GoFrame/src/components/config"
	"GoFrame/src/components/log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type DownFileTimer struct {
	downPath string
}

func (s *DownFileTimer) Enable() bool {
	return len(config.Instance.FTPInfo.Addr) != 0 && len(config.Instance.DownLoad.Path) != 0
}

func (s *DownFileTimer) Init() error {
	s.downPath = filepath.Join(comm.Pwd(), "down")
	os.RemoveAll(s.downPath)
	os.MkdirAll(s.downPath, 0755)
	return nil
}

func (s *DownFileTimer) Uninit() {
}

func (s *DownFileTimer) Name() string {
	return "下载文件"
}

func (s *DownFileTimer) Proc() time.Duration {
	ftp := comm.FTP{
		Addr:     config.Instance.FTPInfo.Addr,
		User:     config.Instance.FTPInfo.User,
		Password: config.Instance.FTPInfo.Password,
	}

	for {
		if s.isStopDown() {
			return 1 * time.Minute
		}

		fs, err := ftp.FileList(config.Instance.FTPInfo.RootPath)
		if err != nil {
			log.Sys.Errorf("获取文件列表失败，原因：%s", err.Error())
			break
		}

		for _, file := range fs {
			if s.isStopDown() {
				return 1 * time.Minute
			}

			v := strings.Split(file.Name, "_")
			if len(v) != 3 {
				continue
			}

			//判断文件前缀
			if !strings.EqualFold(v[0], config.Instance.PackPrefix) {
				continue
			}

			//判断文件大小是否完整
			fsize, err := strconv.ParseUint(v[2], 10, 64)
			if err != nil {
				log.Sys.Errorf("字符串文件大小转换失败，文件：%s，原因：%s", file.Name, err.Error())
				continue
			}

			if fsize != file.Size {
				log.Sys.Errorf("文件不完整：%s，完整大小：%s，实际大小：%d", file.Name, v[2], file.Size)
				continue
			}

			//下载文件
			downFile := filepath.Join(config.Instance.FTPInfo.RootPath, file.Name)
			downFile = filepath.ToSlash(downFile)

			err = ftp.Down(downFile, filepath.Join(s.downPath, filepath.Base(downFile)))
			if err != nil {
				log.Sys.Errorf("下载文件失败：%s，原因：%s", file.Name, err.Error())
				continue
			}

			//删除ftp上的文件
			ftp.Delete(file.Name)
		}

		break
	}

	return 30 * time.Second
}

func (s *DownFileTimer) isStopDown() bool {
	fcount, err := comm.FileCount(s.downPath)
	if err != nil {
		log.Sys.Errorf("获取下载目录中的文件数失败，原因：%s", err.Error())
		return true
	}

	if fcount >= config.Instance.DownLoad.DownMaxFile {
		log.Sys.Errorf("下载目录中的文件数大于%d，停止下载", config.Instance.DownLoad.DownMaxFile)
		return true
	}

	return false
}
