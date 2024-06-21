package model

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileStatus string

const (
	FileStatusNew FileStatus = "NEW"
)

type File struct {
	gorm.Model
	ID       uuid.UUID
	FullName string
	Size     int64     // length in bytes for regular files; system-dependent for others
	ModTime  time.Time // modification time
	Status   FileStatus
	Sha256   *string // 32 Bytes
}

func (ff *File) CalcSha256() error {

	f, err := os.Open(ff.FullName)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	sh := hex.EncodeToString(h.Sum(nil))
	ff.Sha256 = &sh
	return nil
}
