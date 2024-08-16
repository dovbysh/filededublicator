package storage

import "github.com/dovbysh/filededublicator.git/pkg/model"

type IFile interface {
	Create(f *model.File) error
	Get(fileFullName string) (*model.File, error)
}
