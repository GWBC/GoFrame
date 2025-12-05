package db

import "time"

type FileInfo struct {
	Path     string `gorm:"primaryKey"`
	Name     string `gorm:"index:idx_name"`
	ISUpload int    `gorm:"index:idx_is_upload"`
	ModifyAt time.Time
}
