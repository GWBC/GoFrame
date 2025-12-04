package main

import (
	"GoFrame/src/process"
	svr "GoFrame/src/service"
	"log"
	"os"
	"strings"

	"github.com/kardianos/service"
)

/*
#安装服务
程序 instal
#卸载服务
程序 uninstal
#启动服务
程序 star
#重启服务
程序 restar
# 停止服务
程序 stop
*/

type program struct {
}

func (p *program) Start(s service.Service) error {
	err := process.PreProcess()
	if err != nil {
		return err
	}

	return p.run() //p.run函数不能阻塞，返回后由Service接管
}

func (p *program) run() error {
	return svr.Start()
}

func (p *program) Stop(s service.Service) error {
	svr.Stop()
	process.PostProcess()
	return nil
}

func RunService() error {
	ctrlCmd := map[string]bool{
		"install":   true,
		"uninstall": true,
		"start":     true,
		"restart":   true,
		"stop":      true,
	}

	svcConfig := &service.Config{
		Name:        "GoFrame",
		DisplayName: "GoFrame",
		Description: "Go Frame",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) == 2 {
		cmd := strings.TrimSpace(os.Args[1])

		if ctrlCmd[cmd] {
			err = service.Control(s, cmd)
			if err != nil {
				log.Fatal(err)
			}

			return nil
		}
	}

	return s.Run()
}
