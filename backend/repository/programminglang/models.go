package programminglang

import (
	"github.com/randspace0/ketikode/backend/repository"
)

type ProgrammingLang struct {
	LanguageID string `gorm:"primaryKey"`
	Name       string `gorm:"not null"`
}

func AllList() []ProgrammingLang {
	var langs []ProgrammingLang
	if err := repository.DB.Order("language_id").Find(&langs).Error; err != nil {
		return nil
	}
	return langs
}
