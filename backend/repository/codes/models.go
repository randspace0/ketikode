package codes

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/randspace0/ketikode/backend/repository"
)

type Codes struct {
	ID        uint64 `gorm:"primaryKey"`
	Language  string `gorm:"not null"`
	Code      string `gorm:"type:text;not null"`
	Stdin     string `gorm:"type:text"`
	Stdout    string `gorm:"type:text"`
	Stderr    string `gorm:"type:text"`
	Owner     uint64 `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Codes) Create() error {
	return repository.DB.Create(c).Error
}

// FindByID returns nil, nil when the code does not exist.
func (c *Codes) FindByID() (*Codes, error) {
	var found Codes
	err := repository.DB.First(&found, c.ID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &found, nil
}

func (c *Codes) FindAllOwnedCodesByUserID() ([]Codes, error) {
	found := []Codes{}
	err := repository.DB.Where("owner = ?", c.Owner).Order("id DESC").Find(&found).Error
	return found, err
}

// List returns one page of submissions, newest first, optionally filtered by language.
func List(language string, limit, offset int) ([]Codes, error) {
	found := []Codes{}
	q := repository.DB.Order("id DESC").Limit(limit).Offset(offset)
	if language != "" {
		q = q.Where("language = ?", language)
	}
	return found, q.Find(&found).Error
}

// ErrNotFound is returned by Delete when no code of that owner matches.
var ErrNotFound = errors.New("code not found")

// Delete removes the code only if it belongs to c.Owner.
func (c *Codes) Delete() error {
	res := repository.DB.Where("id = ? AND owner = ?", c.ID, c.Owner).Delete(&Codes{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// Total counts submissions, filtered by c.Language when set.
func (c *Codes) Total() (int64, error) {
	var total int64
	q := repository.DB.Model(&Codes{})
	if c.Language != "" {
		q = q.Where("language = ?", c.Language)
	}
	return total, q.Count(&total).Error
}
