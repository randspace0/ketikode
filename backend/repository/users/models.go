package users

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/randspace0/ketikode/backend/repository"
)

type Users struct {
	ID        uint64 `gorm:"primaryKey"`
	Username  string `gorm:"uniqueIndex;not null"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// FindByID returns nil, nil when the user does not exist.
func (u *Users) FindByID() (*Users, error) {
	var found Users
	err := repository.DB.First(&found, u.ID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &found, nil
}

func (u *Users) Create() error {
	return repository.DB.Create(u).Error
}

// MaxPasswordLength is bcrypt's input limit; longer passwords are silently truncated.
const MaxPasswordLength = 72

// HashPassword returns the bcrypt hash to store in Users.Password.
func HashPassword(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(hash), err
}

// FindByUsername returns nil, nil when the user does not exist.
func FindByUsername(username string) (*Users, error) {
	var found Users
	err := repository.DB.Where("username = ?", username).First(&found).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &found, nil
}

func (u *Users) CheckPassword(plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain)) == nil
}

// UpdatePassword stores an already hashed password.
func UpdatePassword(id uint64, hash string) error {
	return repository.DB.Model(&Users{}).Where("id = ?", id).Update("password", hash).Error
}

// UsernamesByID maps user IDs to usernames in a single query.
func UsernamesByID(ids []uint64) (map[uint64]string, error) {
	names := map[uint64]string{}
	if len(ids) == 0 {
		return names, nil
	}

	var found []Users
	if err := repository.DB.Where("id IN ?", ids).Find(&found).Error; err != nil {
		return nil, err
	}
	for _, u := range found {
		names[u.ID] = u.Username
	}
	return names, nil
}
