package timer

import (
	"GoFrame/src/components/log"
	"context"
	"sync"
	"time"
)

var timers = []Timer{
	&PackTimer{},
}

var timerGroup = sync.WaitGroup{}
var timerCtx, timerCancelFun = context.WithCancel(context.Background())

type Timer interface {
	Init() error         //初始化
	Uninit()             //卸载
	Name() string        //名称
	Proc() time.Duration //处理
}

func Start() error {
	for _, svr := range timers {
		err := svr.Init()
		if err != nil {
			log.Sys.Errorf("%s定时服务初始化失败，原因：%s", svr.Name(), err.Error())
			return err
		}
	}

	for _, svr := range timers {
		timerGroup.Add(1)

		go func(svr Timer) {
			log.Sys.Infof("%s定时服务启动", svr.Name())

			d := svr.Proc()
			ticker := time.NewTicker(d)

			defer func() {
				log.Sys.Infof("%s定时服务退出", svr.Name())
				ticker.Stop()
				timerGroup.Done()
			}()

			for {
				select {
				case <-timerCtx.Done():
					return
				case <-ticker.C:
					tmp := svr.Proc()
					if tmp != d {
						ticker.Reset(tmp)
						d = tmp
					}
				}
			}
		}(svr)
	}

	return nil
}

func Stop() {
	timerCancelFun()

	for _, svr := range timers {
		svr.Uninit()
	}

	timerGroup.Wait()
}
