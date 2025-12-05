package comm

import (
	"errors"
	"io/fs"
	"path/filepath"
	"sync"
)

type ResultFun func(path string, info fs.FileInfo) error

type ScanFile struct {
	wg      sync.WaitGroup
	isRun   bool
	count   int
	dirChan chan string
}

func (s *ScanFile) Start(dir string, fun ResultFun) error {
	s.count = 20
	s.isRun = true
	s.dirChan = make(chan string, s.count*10)

	for range s.count {
		s.wg.Add(1)

		go func() {
			defer s.wg.Done()

			for dirPath := range s.dirChan {
				s.processDir(dirPath, fun)
			}
		}()
	}

	s.dirChan <- dir

	return nil
}

func (s *ScanFile) Stop() {
	s.isRun = false
	close(s.dirChan)
	s.wg.Wait()
}

func (s *ScanFile) Scan(dir string) {
	s.dirChan <- dir
}

func (s *ScanFile) processDir(dirPath string, fun ResultFun) {
	filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if !s.isRun {
			return errors.New("exit")
		}

		if err != nil {
			return err
		}

		if info.IsDir() && path != dirPath {
			select {
			case s.dirChan <- path:
			default:
				//异步队列已满，则使用同步
				s.processDir(path, fun)
			}

			return filepath.SkipDir
		}

		return fun(path, info)
	})
}
