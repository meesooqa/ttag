package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
)

type Finder struct {
	log *zap.Logger
}

func NewFinder(log *zap.Logger) *Finder {
	return &Finder{log: log}
}

func (f *Finder) FindFiles(fileOrDirPath string, filesChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(filesChan)

	i, e := os.Stat(fileOrDirPath)
	if e != nil {
		f.log.Error("Can't stat file", zap.String("fileOrDirPath", fileOrDirPath), zap.Error(e))
		return
	}

	if i.IsDir() {
		f.findFilesInDir(fileOrDirPath, filesChan)
	} else {
		// is file
		filesChan <- fileOrDirPath
	}
}

func (f *Finder) findFilesInDir(root string, filesChan chan<- string) {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			f.log.Error(fmt.Sprintf("Error while walking %q: %v\n", path, err))
		}

		if !info.IsDir() {
			filesChan <- path
		}
		return nil
	})
	if err != nil {
		f.log.Error(fmt.Sprintf("Directory walk error: %v\n", err))
	}
}
