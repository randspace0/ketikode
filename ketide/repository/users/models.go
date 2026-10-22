package users

import (
	"time"

	"gitlab.com/cikadev/ketide/repository"
)

type Users struct {
	ID        uint64    `xorm:"unique pk not null autoincr"`
	Username  string    `xorm:"unique not null"`
	Email     string    `xorm:"unique not null"`
	Password  string    `xorm:"not null"`
	CreatedAt time.Time `xorm:"created not null"`
	UpdatedAt time.Time `xorm:"updated not null"`
}

func (u *Users) FindByID() (*Users, error) {
	users := make([]Users, 0)

	err := repository.DB.Find(&users, &Users{
		ID: u.ID,
	})

	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, err
	}

	return &users[0], nil
}

func (u *Users) Create() error {
	_, err := repository.DB.InsertOne(u)
	if err != nil {
		return err
	}
	return nil
}

func (u *Users) GetUserByUsernamePassword() (*Users, error) {
	users := make([]Users, 0)

	err := repository.DB.Find(&users, &Users{
		Username: u.Username,
		Password: u.Password,
	})

	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, err
	}

	return &users[0], nil
}

func (u *Users) UpdateWhereID(id uint64) error {
	_, err := repository.DB.Id(id).Update(u)
	return err
}
