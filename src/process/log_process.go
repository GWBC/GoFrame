package process

import (
	"GoFrame/src/components/config"
	"GoFrame/src/components/log"
)

type LogProcess struct {
}

func (l *LogProcess) PreProcess() error {
	log.Sys.Infof("系统版本：%s", config.Instance.System.Version)
	return nil
}

func (l *LogProcess) PostProcess() {
	log.Instance.Close()
}
