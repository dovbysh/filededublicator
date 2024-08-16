package storage

import (
	"errors"
	"os"
	"testing"

	"github.com/dovbysh/filededublicator.git/pkg/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewSqliteFile(t *testing.T) {
	d, err := os.Getwd()
	assert.NoError(t, err)

	dbPath := d + "/test.db"
	err = os.Remove(dbPath)
	t.Log(err)
	for i := 0; i < 2; i++ {
		func() {
			storage1, err := NewSqliteFile(dbPath)
			assert.NoError(t, err)
			_ = storage1
		}()
	}
}

func TestCreateGet(t *testing.T) {
	d, err := os.Getwd()
	assert.NoError(t, err)

	dbPath := d + "/create.db"
	err = os.Remove(dbPath)
	t.Log(err)

	storage1, err := NewSqliteFile(dbPath)
	assert.NoError(t, err)

	sha := "234243232144"
	f := model.File{
		ID:       uuid.New(),
		FullName: "/z",
		Size:     123,
		Status:   model.FileStatusNew,
		Sha256:   &sha,
	}
	err = storage1.Create(&f)
	assert.NoError(t, err)

	f2, err := storage1.Get(f.FullName)
	assert.True(t, f.Same(f2))
	assert.NoError(t, err)

	f3, err := storage1.Get("/not-exists.txt")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	assert.Nil(t, f3)
}
