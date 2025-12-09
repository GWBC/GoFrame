package timer

import (
	"GoFrame/src/components/config"
	"time"
)

type DeleteFileTimer struct {
}

func (s *DeleteFileTimer) Enable() bool {
	return config.Instance.UpLoad.ISDelFile
}

func (s *DeleteFileTimer) Init() error {
	return nil
}

func (s *DeleteFileTimer) Uninit() {
}

func (s *DeleteFileTimer) Name() string {
	return "删除文件"
}

func (s *DeleteFileTimer) Proc() time.Duration {
	return 30 * time.Second
}
