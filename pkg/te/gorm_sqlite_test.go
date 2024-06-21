package te

import (
	"os"
	"testing"

	"github.com/dovbysh/filededublicator.git/pkg/model"
	"github.com/glebarez/sqlite"
	_ "github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	intVal = map[bool]int{true: 1}
)

func TestGormSqlight(t *testing.T) {
	d, err := os.Getwd()
	assert.NoError(t, err)

	db, err := gorm.Open(sqlite.Open("te.db"), &gorm.Config{})
	assert.NoError(t, err)
	assert.NotNil(t, db)

	db.AutoMigrate(&model.File{})

	de, err := os.ReadDir(d)
	assert.NoError(t, err)
	var dirs, symlinks, files int
	for _, v := range de {
		i, err := v.Info()
		assert.NoError(t, err)
		assert.NotNil(t, i)

		t.Log(v.Name(), ":", v.IsDir(), v.Type(), v.Type().IsRegular(), i.Mode().IsRegular(), v.Type().IsDir())
		files += intVal[v.Type().IsRegular()]
		dirs += intVal[v.Type().IsDir()]
		symlinks += intVal[((i.Mode() & os.ModeSymlink) != 0)]

		if v.Type().IsRegular() {
			f := model.File{
				ID:       uuid.New(),
				FullName: d + "/" + v.Name(),
				Size:     i.Size(),
				ModTime:  i.ModTime(),
				Status:   model.FileStatusNew,
			}
			assert.NoError(t, f.CalcSha256())
			db.Create(&f)
		}
	}
}
