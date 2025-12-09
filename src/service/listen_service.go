package service

import (
	"GoFrame/src/components/config"
	"GoFrame/src/components/db"
	"GoFrame/src/components/log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

type ListenService struct {
	watcher *fsnotify.Watcher
}

func (s *ListenService) Enable() bool {
	return len(config.Instance.UpLoad.Path) != 0
}

func (s *ListenService) Init() error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	s.watcher = w
	return nil
}
func (s *ListenService) Uninit() {
	s.watcher.Close()
}

func (s *ListenService) Name() string {
	return "文件监听"
}

func (s *ListenService) Proc() {
	d := 10 * time.Second
	t := time.NewTicker(d)
	defer t.Stop()

	delFiles := []string{}
	files := map[string]time.Time{}

	s.watcher.Add(config.Instance.UpLoad.Path)

	for {
		select {
		case event, ok := <-s.watcher.Events:
			if !ok {
				return
			}

			nowTime := time.Now()

			if event.Has(fsnotify.Create) {
				files[event.Name] = nowTime
			} else if event.Has(fsnotify.Write) {
				files[event.Name] = nowTime
			} else if event.Has(fsnotify.Remove) ||
				event.Has(fsnotify.Rename) {
				delFiles = append(delFiles, event.Name)
			}
		case err, ok := <-s.watcher.Errors:
			if !ok {
				return
			}

			log.Sys.Errorf("监听文件服务发生错误，Err：%s", err.Error())
		case <-t.C:
			//执行删除
			if len(delFiles) != 0 {
				result := db.Instance.Delete(&db.FileInfo{}, delFiles)
				if result.Error != nil {
					log.Sys.Debugf("删除文件信息失败，原因：%s", result.Error)
					continue
				}

				//将创建中的文件删除
				for _, file := range delFiles {
					delete(files, file)
				}

				delFiles = []string{}
				log.Sys.Debugf("删除文件信息，个数：%d", result.RowsAffected)
			}

			t := time.Now()
			for file, tm := range files {
				//一段时间没有写入，则认为文件已写完
				if t.Sub(tm) < 60*time.Second {
					continue
				}

				finfo, err := os.Stat(file)
				if err != nil {
					continue
				}

				info := &db.FileInfo{}
				info.Path = file
				info.Name = filepath.Base(info.Path)
				info.Ext = filepath.Ext(info.Path)
				info.ModifyAt = finfo.ModTime()
				info.Flag = 0

				result := db.Instance.Save(info)
				if result.Error != nil {
					log.Sys.Errorf("写入文件信息失败，原因：%s", result.Error.Error())
					continue
				}

				log.Sys.Debugf("新增文件:%s", file)
				delete(files, file)
			}

		}
	}
}

func (s *ListenService) SubMessage() []int {
	return []int{}
}

func (s *ListenService) ProcMessage(id int, args ...any) {

}
