package config

import (
	"GoFrame/src/components/comm"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type System struct {
	Services bool `yaml:"Services"`
}

type Log struct {
	MaxSize     int          `yaml:"MaxSize"`     //单个文件最大大小，单位MB
	MaxBackups  int          `yaml:"MaxBackups"`  //切割文件后，最大备份文件数，单位个数
	MaxAge      int          `yaml:"MaxAge"`      //文件最大保留天数，单位天
	Level       logrus.Level `yaml:"Level"`       //日志等级
	IsOutputStd bool         `yaml:"IsOutputStd"` //是否输出到终端
}

type Config struct {
	System System `yaml:"System"`
	Log    Log    `yaml:"Log"`
}

var configPath = filepath.Join(comm.Pwd(), "config", "config.yml")
var configBackPath = configPath + ".back"

func (c *Config) initConfig() error {
	err := os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		return err
	}

	//处理备份文件
	_, err = os.Stat(configBackPath)
	if err == nil {
		os.Remove(configPath)
		err = os.Rename(configBackPath, configPath)
		if err != nil {
			return err
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			log := Log{}
			log.MaxSize = 200
			log.MaxAge = 7
			log.MaxBackups = 10
			log.Level = logrus.DebugLevel
			log.IsOutputStd = true

			//不存在，则创建配置文件
			c.Log = log
			return c.Save()
		}

		return err
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		//解析失败，重置配置文件
		return c.Save()
	}

	return err
}

func (c *Config) Save() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(configBackPath, data, 0644)
	if err != nil {
		return err
	}

	os.Remove(configPath)
	return os.Rename(configBackPath, configPath)
}

///////////////////////////////////////////////////

var instance = comm.Single[Config]{}
var Instance = instance.Instance(func() *Config {
	obj := Config{}
	err := obj.initConfig()
	if err != nil {
		panic("初始化配置失败，Err:" + err.Error())
	}

	return &obj
})
