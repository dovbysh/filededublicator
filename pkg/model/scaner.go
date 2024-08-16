package model

import (
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
)

var (
	ErrReadDirError error = fmt.Errorf("read dir error")
)

type Scaner struct {
	StartDir     string
	ReadDirFiles chan File
	ReadDirReady chan struct{}
}

func NewScaner(startDir string) *Scaner {
	return &Scaner{
		StartDir:     startDir,
		ReadDirReady: make(chan struct{}, 1),
	}
}

func (scaner *Scaner) ReadDir() error {
	scaner.ReadDirFiles = make(chan File, 1000)
	defer close(scaner.ReadDirFiles)

	scaner.ReadDirReady <- struct{}{}

	return scaner.readDir(scaner.StartDir)
}

func (scaner *Scaner) readDir(dir string) error {
	de, err := os.ReadDir(dir)
	if err != nil {
		return errors.Join(ErrReadDirError, err)
	}
	for _, v := range de {
		i, err := v.Info()
		if err != nil {
			return errors.Join(ErrReadDirError, err)
		}

		if v.Type().IsRegular() {
			f := File{
				ID:       uuid.New(),
				FullName: dir + "/" + v.Name(),
				Size:     i.Size(),
				ModTime:  i.ModTime(),
				Status:   FileStatusNew,
			}
			scaner.ReadDirFiles <- f
		}
		if v.Type().IsDir() {
			err = scaner.readDir(dir + "/" + v.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}
