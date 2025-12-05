package db

import (
	"GoFrame/src/components/comm"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var instance = comm.Single[gorm.DB]{}
var Instance = instance.Instance(func() *gorm.DB {
	fileName := filepath.Join(comm.Pwd(), "data", "system.db")
	os.MkdirAll(filepath.Dir(fileName), 0755)

	db, err := gorm.Open(sqlite.Open(fileName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("数据库打开失败，原因：" + err.Error())
	}

	//sqlite不支持多连接，这里设置后可以让db在多协程中正常运行
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
})
