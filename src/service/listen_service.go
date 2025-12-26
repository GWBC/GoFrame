package service

import (
	"GoFrame/src/components/config"
	"GoFrame/src/components/db"
	"GoFrame/src/components/log"
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/radovskyb/watcher"
)

type ListenService struct {
	watcher *watcher.Watcher
	wg      sync.WaitGroup
}

func (s *ListenService) Enable() bool {
	return len(config.Instance.UpLoad.Path) != 0
}

func (s *ListenService) Init() error {
	s.watcher = watcher.New()
	s.watcher.AddRecursive(config.Instance.UpLoad.Path)
	return nil
}
func (s *ListenService) Uninit() {
	s.watcher.Close()
}

func (s *ListenService) Name() string {
	return "文件监听"
}

func (s *ListenService) Proc(ctx context.Context) {
	s.work(ctx)
	s.watcher.Start(10 * time.Second)
	s.wg.Wait()
}

func (s *ListenService) SubMessage() []int {
	return []int{}
}

func (s *ListenService) ProcMessage(id int, args ...any) {

}

func (s *ListenService) work(ctx context.Context) {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		d := 10 * time.Second
		t := time.NewTicker(d)
		defer t.Stop()

		delFiles := []string{}
		files := map[string]time.Time{}

		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-s.watcher.Event:
				if !ok {
					return
				}

				nowTime := time.Now()

				if event.IsDir() {
					continue
				}

				switch event.Op {
				case watcher.Create:
					files[event.Path] = nowTime
				case watcher.Write:
					files[event.Path] = nowTime
				case watcher.Remove:
					delFiles = append(delFiles, event.Path)
				case watcher.Rename:
					files[event.Path] = nowTime
					delFiles = append(delFiles, event.OldPath)
				}
			case err := <-s.watcher.Error:
				log.Sys.Errorf("文件监听发生错误，退出监听，原因：%s", err.Error())
				return
			case <-t.C:
				//执行删除
				const batchSize = 2000
				for i := 0; i < len(delFiles); i += batchSize {
					end := i + batchSize
					if end > len(delFiles) {
						end = len(delFiles)
					}

					result := db.Instance.Delete(&db.FileInfo{}, delFiles[i:end])
					if result.Error != nil {
						log.Sys.Debugf("删除文件信息失败，原因：%s", result.Error.Error())
						continue
					}

					//将创建中的文件删除
					for _, file := range delFiles {
						delete(files, file)
					}

					delFiles = []string{}
					if result.RowsAffected != 0 {
						log.Sys.Debugf("删除文件信息，个数：%d", result.RowsAffected)
					}
				}

				t := time.Now()
				for file, tm := range files {
					//一段时间没有写入，则认为文件已写完
					if t.Sub(tm) < 1*time.Minute {
						continue
					}

					finfo, err := os.Stat(file)
					if err != nil {
						continue
					}

					info := &db.FileInfo{}
					info.Path = file
					info.Name = filepath.Base(info.Path)
					info.ModifyAt = finfo.ModTime()
					info.Flag = 0
					info.Ext = filepath.Ext(info.Path)

					//处理忽略后缀
					for _, suffix := range config.Instance.UpLoad.IgnoreSuffix {
						if !strings.HasSuffix(info.Path, suffix) {
							continue
						}

						info.Ext = filepath.Ext(strings.TrimSuffix(info.Path, suffix))
						break
					}

					if len(info.Ext) == 0 {
						info.Ext = "unknown"
					}

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
	}()
}
