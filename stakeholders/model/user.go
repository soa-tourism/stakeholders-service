package model

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole int32

const (
	Administrator UserRole = iota
	Author
	Tourist
)

type User struct {
	Id       int64    `gorm:"primaryKey" json:"id"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Role     UserRole `json:"role"`
	IsActive bool     `json:"isActive"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if err := u.Validate(); err != nil {
		return err
	}

	uid := uuid.New()
	u.Id = int64(uid.ID())

	return nil
}

func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("invalid username")
	}
	if u.Password == "" {
		return errors.New("invalid password")
	}

	return nil
}
