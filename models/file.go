package models

import (
	"time"
)

type File struct {
	ID             string `gorm:"primaryKey"`
	Filename       string
	FolderID       string
	ThumbnailImg   string
	CreatedOnInUTC time.Time
	CreatedBy      string
	UpdatedOnInUTC time.Time
	UpdatedBy      string
	FilePath       string
}
