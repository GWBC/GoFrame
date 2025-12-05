package service

import (
	"GoFrame/src/components/log"
	"GoFrame/src/message"
	"sync"
)

var serviceList = []Service{
	&ScanService{},
	&ListenService{},
}

var serviceGroup = sync.WaitGroup{}

type Service interface {
	Init() error                     //初始化
	Uninit()                         //卸载
	Name() string                    //名称
	Proc()                           //处理
	SubMessage() []int               //订阅消息
	ProcMessage(id int, args ...any) //处理消息
}

func Start() error {
	for _, svr := range serviceList {
		err := svr.Init()
		if err != nil {
			log.Sys.Errorf("%s服务初始化失败，错误原因：%s", svr.Name(), err.Error())
			return err
		}
	}

	for _, svr := range serviceList {
		serviceGroup.Add(1)

		go func(svr Service) {
			log.Sys.Infof("%s服务启动", svr.Name())

			subMsg := map[int]int{}

			defer func() {
				log.Sys.Infof("%s服务退出", svr.Name())
				for msgID, subID := range subMsg {
					message.Unsub(msgID, subID)
				}

				serviceGroup.Done()
			}()

			msgs := svr.SubMessage()
			for _, msgID := range msgs {
				subMsg[msgID] = message.Sub(msgID, svr.ProcMessage)
			}

			svr.Proc()
		}(svr)
	}

	return nil
}

func Stop() {
	for _, svr := range serviceList {
		svr.Uninit()
	}

	serviceGroup.Wait()
}
