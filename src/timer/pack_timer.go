package timer

import (
	"GoFrame/src/components/comm"
	"os"
	"path/filepath"
	"time"
)

type PackTimer struct {
	packPath string
}

func (s *PackTimer) Init() error {
	s.packPath = filepath.Join(comm.Pwd(), "pack")
	os.MkdirAll(s.packPath, 0755)
	return nil
}
func (s *PackTimer) Uninit() {
}

func (s *PackTimer) Name() string {
	return "打包"
}

func (s *PackTimer) Proc() time.Duration {
	for {
		time.Sleep(1 * time.Second)
	}
	return 1 * time.Second
}
