package storage

import (
	"github.com/dovbysh/filededublicator.git/pkg/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type SqliteFile struct {
	db     *gorm.DB
	dbPath string
}

func NewSqliteFile(dbPath string) (IFile, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	r := &SqliteFile{
		db:     db,
		dbPath: dbPath,
	}
	err = r.init()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *SqliteFile) init() error {
	return s.db.AutoMigrate(&model.File{})
}

func (s *SqliteFile) Create(f *model.File) error {
	tx := s.db.Create(f)
	return tx.Error
}

func (s *SqliteFile) Get(fileFullName string) (*model.File, error) {
	var f model.File
	tx := s.db.Where("full_name =?", fileFullName).First(&f)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &f, nil
}
