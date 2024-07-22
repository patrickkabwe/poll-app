package user

import (
	"poll-app/data"

	"gorm.io/gorm"
)

type UserStore interface {
	CreateUser(user *data.User) error
	GetByEmail(email string) (*data.User, error)
}

type userStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) UserStore {
	return &userStore{
		db: db,
	}
}

func (s *userStore) CreateUser(user *data.User) error {
	return s.db.Create(user).Error
}

func (s *userStore) GetByEmail(email string) (*data.User, error) {
	var user data.User
	err := s.db.Where("email = ?", email).First(&user).Error
	return &user, err
}
