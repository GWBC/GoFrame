package service

import (
	"GoFrame/src/components/config"
	"GoFrame/src/components/db"
	"GoFrame/src/components/log"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

type ListenService struct {
	watcher *fsnotify.Watcher
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
	if len(config.Instance.Sync.Path) == 0 {
		return
	}

	s.watcher.Add(config.Instance.Sync.Path)

	for {
		select {
		case event, ok := <-s.watcher.Events:
			if !ok {
				return
			}

			info := &db.FileInfo{}
			info.Path = event.Name
			info.Name = filepath.Base(event.Name)
			info.Ext = filepath.Ext(event.Name)
			info.ModifyAt = time.Now()
			info.ISUpload = 0

			if event.Has(fsnotify.Create) {
				db.Instance.Save(info)
				log.Sys.Debugf("新增文件:%s", event.Name)
			} else if event.Has(fsnotify.Remove) ||
				event.Has(fsnotify.Rename) {
				db.Instance.Delete(info)
				log.Sys.Debugf("删除文件:%s", event.Name)
			}
		case err, ok := <-s.watcher.Errors:
			if !ok {
				return
			}

			log.Sys.Errorf("监听文件服务发生错误，Err：%s", err.Error())
		}
	}
}

func (s *ListenService) SubMessage() []int {
	return []int{}
}

func (s *ListenService) ProcMessage(id int, args ...any) {

}
