package repo

import (
	"stakeholders/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.DB.Where("username = ? AND is_active = true", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetPersonByUserID(userID int64) (*model.Person, error) {
	var person model.Person
	if err := r.DB.Where("user_id = ?", userID).First(&person).Error; err != nil {
		return nil, err
	}
	return &person, nil
}

func (r *UserRepository) CreateUser(user *model.User) error {
	if err := r.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) CreatePerson(person *model.Person) error {
	if err := r.DB.Create(person).Error; err != nil {
		return err
	}
	return nil
}
