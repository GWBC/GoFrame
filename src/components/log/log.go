package log

import (
	"GoFrame/src/components/comm"
	"GoFrame/src/components/config"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Log struct {
	isOutputStd bool
	logger      *logrus.Logger
	wirter      *lumberjack.Logger
}

func (l *Log) Write(p []byte) (n int, err error) {
	if l.isOutputStd {
		fmt.Print(string(p))
	}

	return l.wirter.Write(p)
}

func (l *Log) Init(fileName string, cfg config.Log) error {
	err := os.MkdirAll(filepath.Dir(fileName), 0755)
	if err != nil {
		return err
	}

	l.wirter = &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   true,
	}

	l.logger = logrus.New()
	l.logger.SetOutput(l)
	l.logger.SetLevel(cfg.Level)
	l.isOutputStd = cfg.IsOutputStd

	formatter := &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}

	l.logger.SetFormatter(formatter)

	return err
}

func (l *Log) Logger() *logrus.Logger {
	return l.logger
}

func (l *Log) Close() {
	if l.wirter == nil {
		return
	}

	l.wirter.Close()
}

///////////////////////////////////////////////////////

var instance = comm.Single[Log]{}
var Instance = instance.Instance(func() *Log {
	obj := Log{}
	err := obj.Init(filepath.Join(comm.Pwd(), "log", "system.log"), config.Instance.Log)
	if err != nil {
		panic("创建日志组件失败，" + err.Error())
	}

	return &obj
})

var Sys = Instance.Logger()
