package main

import (
	"GoFrame/src/components/config"
	"GoFrame/src/process"
	svr "GoFrame/src/service"
	"GoFrame/src/timer"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if config.Instance.System.Services {
		//以服务方式运行
		RunService()
		return
	}

	//前置处理
	err := process.PreProcess()
	if err != nil {
		return
	}

	//启动后台服务
	err = svr.Start()
	if err != nil {
		return
	}

	//启动定时器服务
	err = timer.Start()
	if err != nil {
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	timer.Stop()
	svr.Stop()

	//后置处理
	process.PostProcess()
}
