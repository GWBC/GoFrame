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
	packPath    string
	tmpPackPath string
}

func (s *PackTimer) Enable() bool {
	return len(config.Instance.UpLoad.Path) != 0
}

func (s *PackTimer) Init() error {
	s.packPath = filepath.Join(comm.Pwd(), "pack")
	s.tmpPackPath = filepath.Join(comm.Pwd(), "temp")
	os.MkdirAll(s.packPath, 0755)
	return nil
}

func (s *PackTimer) Uninit() {
}

func (s *PackTimer) Name() string {
	return "打包"
}

func (s *PackTimer) Proc() time.Duration {
	for {
		fcount, err := comm.FileCount(s.packPath)
		if err != nil {
			log.Sys.Errorf("获取打包目录下文件数失败，原因：%s", err.Error())
			break
		}

		if fcount >= config.Instance.UpLoad.PackMaxFile {
			log.Sys.Errorf("打包目录下文件数大于%d，停止打包", config.Instance.UpLoad.PackMaxFile)
			return 1 * time.Minute
		}

		os.RemoveAll(s.tmpPackPath)
		os.MkdirAll(s.packPath, 0755)

		//获取打包文件
		flist := []db.FileInfo{}
		tx := db.Instance.Limit(config.Instance.UpLoad.PackCount).Where("flag=0")
		if len(config.Instance.UpLoad.PackFilter) != 0 {
			tx = tx.Where("ext not in ?", config.Instance.UpLoad.PackFilter)
		}
		result := tx.Order("modify_at").Find(&flist)
		if result.Error != nil {
			log.Sys.Errorf("获取打包数失败，原因：%s", result.Error)
			continue
		}

		t := time.Now().Unix()
		selectFiles := []string{}

		for _, file := range flist {
			rel, err := filepath.Rel(config.Instance.UpLoad.Path, file.Path)
			if err != nil {
				continue
			}

			dstPath := filepath.Join(s.tmpPackPath, rel)

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
		baseName := "file_" + strconv.FormatInt(t, 10)
		tmpName := filepath.Join(s.packPath, baseName)

		err = comm.Zip(s.tmpPackPath, config.ZIPPassword, tmpName)
		if err != nil {
			log.Sys.Errorf("打包文件失败，原因：%s", err.Error())
			break
		}

		info, err := os.Stat(tmpName)
		if err != nil {
			log.Sys.Errorf("获取打包文件信息失败，原因：%s", err.Error())
			break
		}

		name := filepath.Join(filepath.Dir(tmpName), baseName+"_"+strconv.FormatInt(info.Size(), 10))

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

	return 30 * time.Second
}
