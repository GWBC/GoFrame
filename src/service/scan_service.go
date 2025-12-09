package service

import (
	"GoFrame/src/components/comm"
	"GoFrame/src/components/config"
	"GoFrame/src/components/db"
	"GoFrame/src/components/log"
	"errors"
	"io/fs"
	"path/filepath"
	"sync"
	"time"

	"gorm.io/gorm/clause"
)

type ScanService struct {
	scan       comm.ScanFile
	lock       sync.Mutex
	batchCount int
	files      []db.FileInfo
	filesChan  chan []db.FileInfo
}

func (s *ScanService) Init() error {
	if len(config.Instance.UpLoad.Path) == 0 {
		return errors.New("扫描路径未配置")
	}

	s.batchCount = 2000
	s.filesChan = make(chan []db.FileInfo, 100)

	return s.scan.Start(config.Instance.UpLoad.Path, func(path string, info fs.FileInfo) error {
		finfo := db.FileInfo{}
		finfo.Path = path
		finfo.Name = filepath.Base(path)
		finfo.Ext = filepath.Ext(path)
		finfo.ModifyAt = info.ModTime()
		finfo.Flag = 0

		if len(s.files) >= s.batchCount {
			s.lock.Lock()
			tmp := s.files
			s.files = []db.FileInfo{}
			s.lock.Unlock()

			s.filesChan <- tmp
		} else {
			s.lock.Lock()
			s.files = append(s.files, finfo)
			s.lock.Unlock()
		}

		return nil
	})
}
func (s *ScanService) Uninit() {
	s.scan.Stop()
	close(s.filesChan)
}

func (s *ScanService) Name() string {
	return "文件扫描"
}

func (s *ScanService) Proc() {
	d := 10 * time.Second
	t := time.NewTimer(d)
	defer t.Stop()

	for {
		select {
		case files := <-s.filesChan:
			result := db.Instance.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "path"}},
				DoUpdates: clause.Assignments(map[string]any{
					"flag": clause.Expr{SQL: `
							CASE 
								WHEN file_infos.modify_at != excluded.modify_at 
								THEN excluded.flag 
								ELSE file_infos.flag
							END,
							modify_at = excluded.modify_at
							where modify_at != excluded.modify_at
						`},
				})}).CreateInBatches(&files, len(files))
			if result.Error != nil {
				log.Sys.Error("写入文件信息失败，原因：", result.Error.Error())
				continue
			}

			log.Sys.Debug("写入文件信息，条数：", result.RowsAffected)
			t.Reset(d)
		case <-t.C:
			s.lock.Lock()
			if len(s.files) != 0 {
				result := db.Instance.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "path"}},
					DoUpdates: clause.Assignments(map[string]any{
						"flag": clause.Expr{SQL: `
							CASE 
								WHEN file_infos.modify_at != excluded.modify_at 
								THEN excluded.flag 
								ELSE file_infos.flag
							END,
							modify_at = excluded.modify_at
							where modify_at != excluded.modify_at
						`},
					})}).CreateInBatches(&s.files, len(s.files))
				if result.Error != nil {
					log.Sys.Error("定时写入文件信息失败，原因：", result.Error.Error())
					continue
				}

				s.files = []db.FileInfo{}
				log.Sys.Debug("定时写入文件信息，条数：", result.RowsAffected)
			}
			s.lock.Unlock()
			t.Reset(d)
		}
	}
}

func (s *ScanService) SubMessage() []int {
	return []int{}
}

func (s *ScanService) ProcMessage(id int, args ...any) {

}
