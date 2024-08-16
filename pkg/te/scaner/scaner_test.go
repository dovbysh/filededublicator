package scaner

import (
	"os"
	"sync"
	"testing"

	"github.com/dovbysh/filededublicator.git/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestScanerReadDir(t *testing.T) {
	d, err := os.Getwd()
	assert.NoError(t, err)

	scaner := model.NewScaner(d)

	files := make(map[string]bool)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		<-scaner.ReadDirReady
		for f := range scaner.ReadDirFiles {
			t.Log(f.FullName, f.ModTime)
			files[f.FullName] = true
		}
		wg.Done()
	}()
	err = scaner.ReadDir()

	assert.NoError(t, err)
	wg.Wait()
	assert.True(t, files[d+"/d1/f1.txt"])
	assert.True(t, files[d+"/d1/d2/f2.txt"])
	assert.True(t, files[d+"/scaner_test.go"])
}
