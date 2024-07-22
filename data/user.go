package data

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint64     `json:"id" gorm:"primaryKey"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Polls     []Poll     `json:"polls" gorm:"foreignKey:CreatedById"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(*gorm.DB) error {
	now := time.Now()
	u.CreatedAt = &now
	u.UpdatedAt = &now
	u.Password = "123456"
	return nil
}

func (u *User) BeforeUpdate(*gorm.DB) error {
	now := time.Now()
	u.UpdatedAt = &now
	return nil
}

func (u *User) ComparePassword(password string) bool {
	return u.Password == password
}
