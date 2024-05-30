package model

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/mail"
	"unicode"
)

type Person struct {
	Id                int64  `json:"id"`
	UserId            int64  `json:"userId"`
	Name              string `json:"name"`
	Surname           string `json:"surname"`
	ProfilePictureUrl string `json:"profilePictureUrl"`
	Biography         string `json:"biography"`
	Motto             string `json:"motto"`
	Email             string `json:"email"`
}

func (p *Person) BeforeCreate(tx *gorm.DB) error {
	if err := p.Validate(); err != nil {
		return err
	}

	uid := uuid.New()
	p.Id = int64(uid.ID())

	return nil
}

func (p *Person) Validate() error {
	if p.UserId == 0 {
		return errors.New("invalid user id")
	}
	if len(p.Name) == 0 || len(p.Surname) == 0 {
		return errors.New("invalid name or surname")
	}
	if !isValidName(p.Name) || !isValidName(p.Surname) {
		return errors.New("name and surname should start with an uppercase letter")
	}
	if _, err := mail.ParseAddress(p.Email); err != nil {
		return errors.New("invalid email")
	}

	return nil
}

func isValidName(name string) bool {
	if len(name) == 0 {
		return false
	}
	return unicode.IsUpper(rune(name[0]))
}
