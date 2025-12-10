package timer

import (
	"GoFrame/src/components/comm"
	"GoFrame/src/components/config"
	"GoFrame/src/components/db"
	"GoFrame/src/components/log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type PackTimer struct {
	uploadPath      string
	uploadZipPath   string
	uploadFilesPath string
}

func (s *PackTimer) Enable() bool {
	return len(config.Instance.UpLoad.Path) != 0
}

func (s *PackTimer) Init() error {
	s.uploadPath = filepath.Join(comm.Pwd(), "upload")
	s.uploadZipPath = filepath.Join(s.uploadPath, "zip")
	s.uploadFilesPath = filepath.Join(s.uploadPath, "files")
	os.MkdirAll(s.uploadPath, 0755)
	return nil
}

func (s *PackTimer) Uninit() {
}

func (s *PackTimer) Name() string {
	return "打包"
}

func (s *PackTimer) Proc() time.Duration {
	for {
		fcount, err := comm.FileCount(s.uploadPath)
		if err != nil {
			log.Sys.Errorf("获取打包目录中的文件数失败，原因：%s", err.Error())
			break
		}

		if fcount >= config.Instance.UpLoad.PackMaxFile {
			log.Sys.Errorf("打包目录中的文件数大于%d，暂停打包", config.Instance.UpLoad.PackMaxFile)
			return config.ProcInterval(2)
		}

		os.RemoveAll(s.uploadZipPath)
		os.RemoveAll(s.uploadFilesPath)

		os.MkdirAll(s.uploadPath, 0755)
		os.MkdirAll(s.uploadZipPath, 0755)

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

			dstPath := filepath.Join(s.uploadFilesPath, rel)

			_, err = comm.CopyFile(file.Path, dstPath)
			if err != nil {
				continue
			}

			selectFiles = append(selectFiles, file.Path)
		}

		if len(selectFiles) == 0 {
			break
		}

		//打包成zip
		baseName := config.Instance.PackPrefix + "_" + strconv.FormatInt(time.Now().UnixMilli(), 10)
		tmpName := filepath.Join(s.uploadZipPath, baseName)

		err = comm.Zip(s.uploadFilesPath, config.ZIPPassword, tmpName)
		if err != nil {
			log.Sys.Errorf("打包文件失败，原因：%s", err.Error())
			break
		}

		info, err := os.Stat(tmpName)
		if err != nil {
			log.Sys.Errorf("获取打包文件信息失败，原因：%s", err.Error())
			break
		}

		name := filepath.Join(s.uploadPath, baseName+"_"+strconv.FormatInt(info.Size(), 10))

		err = os.Rename(tmpName, name)
		if err != nil {
			log.Sys.Errorf("修改打包文件名失败，原因：%s", err.Error())
			break
		}

		log.Sys.Debugf("生成打包文件：%s", name)

		//修改打包标识
		result = db.Instance.Model(&db.FileInfo{}).Where("path in ?", selectFiles).Update("flag", 1)
		if result.Error != nil {
			log.Sys.Errorf("更新打包标识失败，原因：%s", result.Error.Error())
			break
		}
	}

	return config.ProcInterval(1)
}
