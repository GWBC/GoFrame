package db

import (
	"GoFrame/src/components/comm"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var instance = comm.Single[gorm.DB]{}
var Instance = instance.Instance(func() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("system.db"), &gorm.Config{})
	if err != nil {
		panic("数据库打开失败，Err:" + err.Error())
	}

	//sqlite不支持多连接，这里设置后可以让db在多协程中正常运行
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
})
