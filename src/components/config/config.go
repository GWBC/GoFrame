package config

import (
	"GoFrame/src/components/comm"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const ZIPPassword = "(!@#$%^&*)"

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

type FTPInfo struct {
	Addr     string `yaml:"Addr"`     //服务器地址
	User     string `yaml:"User"`     //用户名
	Password string `yaml:"Password"` //密码
	RootPath string `yaml:"RootPath"` //根路径
}

type UpLoad struct {
	Path              string   `yaml:"Path"`              //同步目录
	PackFilter        []string `yaml:"PackFilter"`        //打包过滤器，文件后缀，例子：.mp4
	PackCount         int      `yaml:"PackCount"`         //打包个数
	PackMaxFile       int      `yaml:"PackMaxFile"`       //打包目录下最大文件数
	ISDelFile         bool     `yaml:"ISDelFile"`         //打包完成后是否删除文件
	FileRetentionTime int      `yaml:"FileRetentionTime"` //删除开启后，保留的时间，单位小时
}

type DownLoad struct {
	Path        string `yaml:"Path"`        //下载写入目录
	DownMaxFile int    `yaml:"PackMaxFile"` //下载目录下最大文件数
}

type Config struct {
	System     System   `yaml:"System"`
	Log        Log      `yaml:"Log"`
	PackPrefix string   `yaml:"PackPrefix"`
	FTPInfo    FTPInfo  `yaml:"FTPInfo"`
	UpLoad     UpLoad   `yaml:"UpLoad"`
	DownLoad   DownLoad `yaml:"DownLoad"`
}

var configPath = filepath.Join(comm.Pwd(), "data", "config.yml")
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
			//不存在，则创建配置文件
			c.Log.MaxSize = 200
			c.Log.MaxAge = 7
			c.Log.MaxBackups = 10
			c.Log.Level = logrus.DebugLevel
			c.Log.IsOutputStd = true

			c.PackPrefix = "file"

			c.FTPInfo.Addr = "172.16.100.223:21"
			c.FTPInfo.User = "zhang"
			c.FTPInfo.Password = "zhang"
			c.FTPInfo.RootPath = "/"

			c.UpLoad.Path = "/root/code/HiSTBLinuxV100R005C00SPC050"
			c.UpLoad.PackFilter = []string{}
			c.UpLoad.PackCount = 40
			c.UpLoad.PackMaxFile = 200
			c.UpLoad.ISDelFile = false
			c.UpLoad.FileRetentionTime = 1

			c.DownLoad.Path = "/root/code/HiSTBLinuxV100R005C00SPC050_back"
			c.DownLoad.DownMaxFile = 200

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
		panic("初始化配置失败，原因：" + err.Error())
	}

	return &obj
})
