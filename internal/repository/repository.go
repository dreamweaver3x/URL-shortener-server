package repository

import (
	"SOKR/internal/models"
	"SOKR/internal/shorturl"
	"fmt"
	"gorm.io/gorm"
	"log"
)

type LinksRepository struct {
	db *gorm.DB
}

func NewLinksRepository(db *gorm.DB) *LinksRepository {
	return &LinksRepository{db: db}
}

func (l *LinksRepository) Create(u *models.Link)  error {
	tx := l.db.Begin()
	result := tx.Create(u)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	u.ShortUrl = shorturl.Encode(int(u.ID))

	result = tx.Save(u)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	result = tx.Commit()
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	return  nil
}

func (l *LinksRepository) GetLongUrl(u *models.Link) (*models.Link, error) {
	result := l.db.Where("short_url = ?", u.ShortUrl).First(&u)
	if result.Error != nil {
		return nil, result.Error
	}
	result = l.db.Model(u).Where("short_url = ?", u.ShortUrl).Update("nums_of_redirects", gorm.Expr("nums_of_redirects + ?", 1))
	if result.Error != nil {
		return nil, result.Error
	}
	return u, nil
}
func (l *LinksRepository) GetFiveHundredLinks(id uint, ch *chan models.Link)  error {
links := make([]models.Link, 0)
result := l.db.Select("id","full_url", "accessible").Where("id BETWEEN ? AND ?", id*500+1, id*500+500).Find(&links)
log.Println("id raven ", id)
if result.Error != nil {
	log.Println(result.Error)
	return result.Error
}
if result.RowsAffected == 0 {
	return fmt.Errorf("net strok na id=%d", id)
}
for i := 0; i < len(links); i++ {
		*ch <- links[i]
}

return nil
}


/*func (l *LinksRepository) GetAllLinks() []models.Link {
	allLinks := make([]models.Link, 0)
	l.db.Select("id","full_url", "accessible").Find(&allLinks)
	return allLinks
}*/

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
	log.Println("URAAAAAAAA")
}

func (l *LinksRepository) GetStats(u *models.Link)  error {
	result := l.db.Where("short_url = ?", u.ShortUrl).First(&u)
	if result.Error != nil {
		return result.Error
	}
	return nil
}


