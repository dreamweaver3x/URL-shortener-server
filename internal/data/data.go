package data

import (
	"SOKR/internal/shorturl"
	"fmt"
	"gorm.io/gorm"
)

type Data struct {
	ID       int  `gorm:"primaryKey"`
	FullUrl    string
	ShortUrl   string
	NumsOfRedirects uint32
	Accessible bool
}

func InitModels(db *gorm.DB) error {
	err := db.AutoMigrate(&Data{})
	if err != nil {
		return err
	}

	return nil
}
func Create(db *gorm.DB) {
	u := Data{
		FullUrl:    "http://google.com",
		ShortUrl:   "aaaaa",
		Accessible: true,
	}
	result := db.Create(&u)
	if result.Error != nil {
		panic(result.Error)
	}
	fmt.Println("1st created")
}

func AddElem(db *gorm.DB, longUrl string) string {
	u := &Data{}
	db.Last(&u)
	u.ID++
	u.ShortUrl = shorturl.Encode(u.ID)
	u.FullUrl = longUrl
	u.NumsOfRedirects = 0
	result := db.Create(u)
	if result.Error != nil {
		return "can't make short url"
	}
	return "http://localhost:8080/" + u.ShortUrl
}

func CheckForElemLong(db *gorm.DB, longUrl string) bool {
	u := &Data{}
	result := db.Where("full_url = ?", longUrl).First(&u)
	if result.Error != nil {
		return false
	}
	return true
}

func CheckForElemShort(db *gorm.DB, shortUrl string) bool {
	u := &Data{}
	println("checkforelemshort = ", shortUrl)
	result := db.Where("short_url = ?", shortUrl).First(&u)
	if result.Error != nil {
		return false
	}
	if u.ShortUrl == "" {
		return false
	}
	return true
}

func GetLongUrl(db *gorm.DB, shortUrl string) string {
	u := &Data{}
	result := db.Where("short_url = ?", shortUrl).First(&u)
	if result.Error != nil {
		return ""
	}
	fullUrl := u.FullUrl
	u.NumsOfRedirects++
	db.Save(u)
	return fullUrl
}
