package timer

import (
	"GoFrame/src/components/log"
	"math/rand"
	"time"
)

type TestTimer struct {
}

func (t *TestTimer) Init() error {
	return nil
}

func (t *TestTimer) Uninit() {

}

func (t *TestTimer) Name() string {
	return "测试"
}

func (t *TestTimer) Proc() time.Duration {
	s := rand.Intn(3) + 1
	log.Sys.Info("定时器测试:", s)
	return time.Duration(s * int(time.Second))
}
