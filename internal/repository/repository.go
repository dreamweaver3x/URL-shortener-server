package repository

import (
	"SOKR/internal/models"
	"SOKR/internal/shorturl"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
)

type LinksRepository struct {
	db *gorm.DB
}

func NewLinksRepository(db *gorm.DB) *LinksRepository {
	return &LinksRepository{db: db}
}

func (l *LinksRepository) Create(u *models.Link) (*models.Link, error) {
	tx := l.db.Begin()
	result := tx.Create(u)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	u.ShortUrl = shorturl.Encode(int(u.ID))

	result = tx.Save(u)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	result = tx.Commit()
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	return u, nil
}

func (l *LinksRepository) CheckForElemLong(u *models.Link) bool {
	u.FullUrl = strings.TrimPrefix(u.FullUrl, "http://")
	u.FullUrl = strings.TrimPrefix(u.FullUrl, "https://")
	u.FullUrl = strings.TrimPrefix(u.FullUrl, "www.")
	result := l.db.Where("full_url = ?", u.FullUrl).First(&u)
	if result.Error != nil {
		return false
	}
	return true
}

func (l *LinksRepository) CheckForElemShort(u *models.Link) bool {
	result := l.db.Where("short_url = ?", u.ShortUrl).First(&u)
	if result.Error != nil {
		return false
	}
	if u.ShortUrl == "" {
		return false
	}
	return true
}

func (l *LinksRepository) GetLongUrl(u *models.Link) (*models.Link, error) {
	result := l.db.Where("short_url = ?", u.ShortUrl).First(&u)
	if result.Error != nil {
		return nil, result.Error
	}
	u.NumsOfRedirects++
	result = l.db.Model(u).Where("short_url = ?", u.ShortUrl).Update("nums_of_redirects", u.NumsOfRedirects)
	if result.Error != nil {
		return nil, result.Error
	}
	return u, nil
}

func (l *LinksRepository) GetAllLinks() []models.Link {
	allLinks := make([]models.Link, 0)
	l.db.Select("id","full_url", "accessible").Find(&allLinks)
	return allLinks
}
func (l *LinksRepository) UpdateAccess(inaccessibleLink, accessibleLink []uint) {
	u := models.Link{}
	result := l.db.Model(u).Where("id in ?", accessibleLink).Update("accessible", true)
	if result.Error != nil {
		log.Println(result.Error)
	}
	result = l.db.Model(u).Where("id in ?", inaccessibleLink).Update("accessible", false)
	if result.Error != nil {
		log.Println(result.Error)
	}
}

func convSliceToStr(slice []uint) string {
	var strSlice []string
	for _, el := range slice {
		strSlice = append(strSlice, strconv.Itoa(int(el)))
	}
	return strings.Join(strSlice, ",")
}

//func (l *LinksRepository) CheckUrlsStatus() string {
//	u := &models.Link{}
//
//	for {
//		l.db.Last(u)
//		for i := 1; i <= int(u.ID)/10+1; i++ {
//			go func(id int) {
//				for k := id - 10; k < id; k++ {
//					link := models.Link{}
//					result := l.db.Where("id = ?", k).First(&link)
//					if result.Error != nil {
//						log.Println(result.Error)
//						continue
//					}
//					resp, err := http.Get("http://" + link.FullUrl)
//
//					if err != nil || resp.StatusCode != 200 {
//						log.Println(err)
//						l.db.Table("links").Where("id = ?", k).Update("accessible", false)
//						log.Println(link.FullUrl, " SSILKA NE RABOTAET")
//						continue
//
//					}
//					l.db.Table("links").Where("id = ?", k).Update("accessible", true)
//					err = resp.Body.Close()
//					if err != nil {
//						log.Println(err)
//					}
//				}
//			}(i * 10)
//		}
//		time.Sleep(time.Minute * 10)
//	}
//}
