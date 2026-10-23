package programminglang

import (
	"gitlab.com/cikadev/ketide/repository"
)

type ProgrammingLang struct {
	LanguageID string `xorm:"unique pk not null"`
	Name       string `xorm:"not null"`
}

func (p *ProgrammingLang) Create() error {
	_, err := repository.DB.InsertOne(p)
	if err != nil {
		return err
	}
	return nil
}

func Total() int64 {
	total, _ := repository.DB.Count(ProgrammingLang{})
	return total
}

func AllList() []ProgrammingLang {
	programmingLangs := make([]ProgrammingLang, 0)

	if err := repository.DB.Find(&programmingLangs, &ProgrammingLang{}); err != nil {
		return nil
	}

	return programmingLangs
}
