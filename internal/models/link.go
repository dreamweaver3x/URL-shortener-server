package models

import (
	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	FullUrl         string
	ShortUrl        string
	NumsOfRedirects uint32 `gorm:"default:0"`
	Accessible      bool   `gorm:"default:true"`
}

func InitModels(db *gorm.DB) error {
	err := db.AutoMigrate(&Link{})
	if err != nil {
		return err
	}

	return nil
}
