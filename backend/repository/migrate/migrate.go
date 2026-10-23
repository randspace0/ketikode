package migrate

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/randspace0/ketikode/backend/repository/codes"
	"github.com/randspace0/ketikode/backend/repository/programminglang"
	"github.com/randspace0/ketikode/backend/repository/users"
)

// Languages must match the executor's languages.go.
var languages = []programminglang.ProgrammingLang{
	{LanguageID: "c", Name: "C"},
	{LanguageID: "cpp", Name: "C++"},
	{LanguageID: "ruby", Name: "Ruby"},
	{LanguageID: "python3", Name: "Python 3"},
	{LanguageID: "java", Name: "Java"},
	{LanguageID: "rust", Name: "Rust"},
}

func Run(db *gorm.DB) error {
	if err := db.AutoMigrate(&users.Users{}, &codes.Codes{}, &programminglang.ProgrammingLang{}); err != nil {
		return err
	}
	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&languages).Error
}
