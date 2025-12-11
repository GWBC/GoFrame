package timer

import (
	"GoFrame/src/components/config"
	"GoFrame/src/components/db"
	"GoFrame/src/components/log"
	"context"
	"os"
	"time"
)

type DeleteFileTimer struct {
}

func (s *DeleteFileTimer) Enable() bool {
	return config.Instance.UpLoad.ISDelFile
}

func (s *DeleteFileTimer) Init() error {
	return nil
}

func (s *DeleteFileTimer) Uninit() {
}

func (s *DeleteFileTimer) Name() string {
	return "删除文件"
}

func (s *DeleteFileTimer) Proc(ctx context.Context) time.Duration {
	for {
		select {
		case <-ctx.Done():
			return config.ProcInterval(1)
		default:
			//获取打包文件
			flist := []db.FileInfo{}
			result := db.Instance.Limit(100).Where("flag=1").Order("modify_at").Find(&flist)
			if result.Error != nil {
				log.Sys.Errorf("获取删除文件列表失败，原因：%s", result.Error)
				continue
			}

			delFiles := []db.FileInfo{}
			now := time.Now()

			for _, file := range flist {
				info, err := os.Stat(file.Path)
				if err != nil {
					if os.IsNotExist(err) {
						delFiles = append(delFiles, file)
					}

					continue
				}

				diff := now.Sub(info.ModTime())

				//修改过系统时间，忽略
				if diff < 0 {
					continue
				}

				//修改时间和当前时间差异一段时间，认为已经没有写入，才可以删除
				if diff <= 5*time.Minute {
					continue
				}

				//保留时间处理
				if diff < time.Duration(config.Instance.UpLoad.FileRetentionTime)*time.Hour {
					continue
				}

				delFiles = append(delFiles, file)
			}

			if len(delFiles) == 0 {
				return config.ProcInterval(1)
			}

			for _, file := range delFiles {
				err := os.Remove(file.Path)
				if err != nil {
					if !os.IsNotExist(err) {
						continue
					}
				}

				db.Instance.Delete(&file)
			}
		}
	}
}
