package timer

import (
	"GoFrame/src/components/config"
	"time"
)

type UploadFTPTimer struct {
}

func (s *UploadFTPTimer) Enable() bool {
	return len(config.Instance.FTPInfo.Addr) != 0
}

func (s *UploadFTPTimer) Init() error {
	return nil
}

func (s *UploadFTPTimer) Uninit() {
}

func (s *UploadFTPTimer) Name() string {
	return "上传文件"
}

func (s *UploadFTPTimer) Proc() time.Duration {
	return 30 * time.Second
}
