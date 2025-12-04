package service

import (
	"GoFrame/src/components/log"
	"GoFrame/src/message"
	"time"
)

type TestService struct {
	isRun bool
}

func (s *TestService) Init() error {
	s.isRun = true
	return nil
}
func (s *TestService) Uninit() {
	s.isRun = false
}

func (s *TestService) Name() string {
	return "test1"
}

func (s *TestService) Proc() {
	for s.isRun {
		log.Sys.Info("test")
		message.Pub(message.MSG_TEST)
		time.Sleep(1 * time.Second)
	}
}

func (s *TestService) SubMessage() []int {
	return []int{message.MSG_TEST}
}

func (s *TestService) ProcMessage(id int, args ...any) {
	log.Sys.Info("收到消息")
}
