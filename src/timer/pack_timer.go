package timer

import (
	"GoFrame/src/components/comm"
	"GoFrame/src/components/config"
	"GoFrame/src/components/db"
	"GoFrame/src/components/log"
	"context"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type PackTimer struct {
	uploadPath      string
	uploadPackPath  string
	uploadFilesPath string
}

func (s *PackTimer) Enable() bool {
	return len(config.Instance.UpLoad.Path) != 0
}

func (s *PackTimer) Init() error {
	s.uploadPath = filepath.Join(comm.Pwd(), "upload")
	s.uploadPackPath = filepath.Join(s.uploadPath, "pack")
	s.uploadFilesPath = filepath.Join(s.uploadPath, "files")
	os.MkdirAll(s.uploadPath, 0755)
	return nil
}

func (s *PackTimer) Uninit() {
}

func (s *PackTimer) Name() string {
	return "打包"
}

func (s *PackTimer) Proc(ctx context.Context) time.Duration {
	for {
		select {
		case <-ctx.Done():
			return config.ProcInterval(1)
		default:
			fcount, err := comm.FileCount(s.uploadPath)
			if err != nil {
				log.Sys.Errorf("获取打包目录中的文件数失败，原因：%s", err.Error())
				return config.ProcInterval(1)
			}

			if fcount >= config.Instance.UpLoad.PackMaxFile {
				log.Sys.Errorf("打包目录中的文件数>=%d，暂停打包", config.Instance.UpLoad.PackMaxFile)
				return config.ProcInterval(2)
			}

			//清理打包文件
			os.RemoveAll(s.uploadPackPath)
			os.RemoveAll(s.uploadFilesPath)
			os.MkdirAll(s.uploadPackPath, 0755)
			os.MkdirAll(s.uploadFilesPath, 0755)

			//获取打包文件
			flist := []db.FileInfo{}
			result := db.Instance.Limit(config.Instance.UpLoad.PackCount).Where("flag=0")
			if len(config.Instance.UpLoad.PackFilter) != 0 {
				result = result.Where("ext not in ?", config.Instance.UpLoad.PackFilter)
			}
			result = result.Order("modify_at").Find(&flist)
			if result.Error != nil {
				log.Sys.Errorf("获取打包数失败，原因：%s", result.Error)
				continue
			}

			selectFiles := []string{}

			for _, file := range flist {
				rel, err := filepath.Rel(config.Instance.UpLoad.Path, file.Path)
				if err != nil {
					continue
				}

				info, err := os.Lstat(file.Path)
				if err != nil {
					continue
				}

				//目录
				if info.IsDir() {
					continue
				}

				dstPath := filepath.Join(s.uploadFilesPath, rel)

				//软连接
				if info.Mode()&os.ModeSymlink != 0 {
					err = comm.CopySymlink(file.Path, dstPath)
				} else {
					_, err = comm.CopyFileAndMeta(file.Path, dstPath)
				}

				if err != nil {
					log.Sys.Errorf("拷贝文件发生错误，原因：%s", err.Error())
					continue
				}

				selectFiles = append(selectFiles, file.Path)
			}

			if len(selectFiles) == 0 {
				return config.ProcInterval(1)
			}

			//打包
			baseName := config.Instance.PackPrefix + "_" + strconv.FormatInt(time.Now().UnixMilli(), 10)
			tmpName := filepath.Join(s.uploadPackPath, baseName)

			err = comm.Tar(s.uploadFilesPath, config.PackPassword, tmpName)
			if err != nil {
				log.Sys.Errorf("打包文件失败，原因：%s", err.Error())
				return config.ProcInterval(1)
			}

			info, err := os.Stat(tmpName)
			if err != nil {
				log.Sys.Errorf("获取打包文件信息失败，原因：%s", err.Error())
				return config.ProcInterval(1)
			}

			name := filepath.Join(s.uploadPath, baseName+"_"+strconv.FormatInt(info.Size(), 10))

			err = os.Rename(tmpName, name)
			if err != nil {
				log.Sys.Errorf("修改打包文件名失败，原因：%s", err.Error())
				return config.ProcInterval(1)
			}

			log.Sys.Debugf("生成打包文件：%s", name)

			//修改打包标识
			result = db.Instance.Model(&db.FileInfo{}).Where("path in ?", selectFiles).Update("flag", 1)
			if result.Error != nil {
				log.Sys.Errorf("更新打包标识失败，原因：%s", result.Error.Error())
				return config.ProcInterval(1)
			}
		}
	}
}
