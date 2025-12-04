package db

import "time"

type FileInfo struct {
	AbsPath      string `gorm:"primaryKey"`
	RelativePath string
	ISUpload     int `gorm:"index:idx_is_upload"`
	Size         int64
	UpdatedAt    time.Time
}
