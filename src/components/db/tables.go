package db

import "time"

type FileInfo struct {
	Path     string `gorm:"primaryKey"`
	Name     string `gorm:"index:idx_name"`
	Ext      string `gorm:"index:idx_ext"`
	ISUpload int    `gorm:"index:idx_is_upload"`
	ModifyAt time.Time
}
