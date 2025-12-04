package process

import "GoFrame/src/components/log"

type LogProcess struct {
}

func (l *LogProcess) PreProcess() error {
	return nil
}

func (l *LogProcess) PostProcess() {
	log.Instance.Close()
}
