package codes

import (
	_ "github.com/lib/pq"

	"time"

	"gitlab.com/cikadev/ketide/repository"
)

type Codes struct {
	ID        uint64 `xorm:"unique pk not null autoincr"`
	Language  string `xorm:"not null"`
	Code      string `xorm:"not null"`
	Stdin     string
	Stdout    string
	Stderr    string
	Owner     uint64    `xorm:"not null"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

func (c *Codes) Create() error {
	_, err := repository.DB.InsertOne(c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Codes) FindByID() (*Codes, error) {
	codes := make([]Codes, 0)

	err := repository.DB.Find(&codes, &Codes{
		ID: c.ID,
	})

	if len(codes) == 0 {
		return nil, err
	}

	if err != nil {
		return &codes[0], err
	}

	return &codes[0], nil
}

func (c *Codes) FindAllOwnedCodesByUserID() ([]Codes, error) {
	codes := []Codes{}

	err := repository.DB.Find(&codes, &Codes{
		Owner: c.Owner,
	})

	if err != nil {
		return codes, err
	}

	return codes, nil
}

func (c *Codes) Delete() error {
	_, err := repository.DB.Delete(c)

	return err
}

func (c *Codes) Total() (int64, error) {
	total, err := repository.DB.Count(c)

	if err != nil {
		return total, err
	}

	return total, nil
}
