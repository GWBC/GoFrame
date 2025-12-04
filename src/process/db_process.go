package process

import (
	"GoFrame/src/components/db"
)

type DBProcess struct {
}

func (d *DBProcess) PreProcess() error {
	return db.Instance.AutoMigrate(&db.FileInfo{})
}

func (d *DBProcess) PostProcess() {
	sqlDB, err := db.Instance.DB()
	if err != nil {
		return
	}

	sqlDB.Close()
}
