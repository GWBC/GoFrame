package timer

import (
	"GoFrame/src/components/config"
	"time"
)

type DeleteFileTimer struct {
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
	if !config.Instance.UpLoad.ISDelFile {
		return -1
	}

	return 30 * time.Second
}
