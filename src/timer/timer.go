package timer

import (
	"GoFrame/src/components/log"
	"context"
	"sync"
	"time"
)

var timers = []Timer{
	&PackTimer{},
	&DeleteFileTimer{},
	&UploadFTPTimer{},
}

var timerGroup = sync.WaitGroup{}
var timerCtx, timerCancelFun = context.WithCancel(context.Background())

type Timer interface {
	Enable() bool        //是否激活
	Init() error         //初始化
	Uninit()             //卸载
	Name() string        //名称
	Proc() time.Duration //处理
}

func Start() error {
	for _, svr := range timers {
		if !svr.Enable() {
			continue
		}

		err := svr.Init()
		if err != nil {
			log.Sys.Errorf("%s定时服务初始化失败，原因：%s", svr.Name(), err.Error())
			return err
		}
	}

	for _, svr := range timers {
		if !svr.Enable() {
			continue
		}

		timerGroup.Add(1)

		go func(svr Timer) {
			log.Sys.Infof("%s定时服务启动", svr.Name())
			defer func() {
				log.Sys.Infof("%s定时服务退出", svr.Name())
				timerGroup.Done()
			}()

			d := svr.Proc()
			if d < 0 {
				return
			}

			tr := time.NewTimer(d)
			defer tr.Stop()

			for {
				select {
				case <-timerCtx.Done():
					return
				case <-tr.C:
					d = svr.Proc()
					if d < 0 {
						return
					}

					tr.Reset(d)
				}
			}
		}(svr)
	}

	return nil
}

func Stop() {
	timerCancelFun()

	for _, svr := range timers {
		if !svr.Enable() {
			continue
		}

		svr.Uninit()
	}

	timerGroup.Wait()
}
