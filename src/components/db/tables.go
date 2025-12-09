package db

import "time"

type FileInfo struct {
	Path     string `gorm:"primaryKey"`     //文件绝对路径
	Name     string `gorm:"index:idx_name"` //文件名称
	Ext      string `gorm:"index:idx_ext"`  //文件类别
	Flag     int    `gorm:"index:idx_flag"` //标识0初始，1打包
	ModifyAt time.Time
}
