package service

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"stakeholders/model"
	"stakeholders/proto/auth"
	"stakeholders/repo"
	"time"
)

type AuthService struct {
	UserRepo *repo.UserRepository
	Key      string
}

func (s *AuthService) Login(credentials *auth.Credentials) (*auth.AuthenticationTokens, error) {
	user, err := s.UserRepo.GetByUsername(credentials.Username)
	if err != nil {
		return nil, err
	}

	if user.Password != credentials.Password {
		return nil, errors.New("wrong password")
	}

	person, err := s.UserRepo.GetPersonByUserID(user.Id)
	if err != nil {
		return nil, err
	}

	tokenString, err := generateJWT(*user, *person, s.Key)
	if err != nil {
		return nil, err
	}

	return &auth.AuthenticationTokens{
		Id:          user.Id,
		AccessToken: tokenString,
	}, nil
}

func (s *AuthService) Register(regReq *auth.AccountRegistrationRequest) (*auth.AccountRegistrationResponse, error) {
	var user = &model.User{
		Username: regReq.Username,
		Password: regReq.Password,
		Role:     model.UserRole(regReq.Role),
		IsActive: true,
	}

	if err := s.UserRepo.CreateUser(user); err != nil {
		return &auth.AccountRegistrationResponse{
			Success: false,
			Id:      -1,
		}, err
	}

	person := &model.Person{
		UserId:  user.Id,
		Name:    regReq.Name,
		Surname: regReq.Surname,
		Email:   regReq.Email,
	}

	if err := s.UserRepo.CreatePerson(person); err != nil {
		return &auth.AccountRegistrationResponse{
			Success: false,
			Id:      -1,
		}, err
	}

	return &auth.AccountRegistrationResponse{
		Success: true,
		Id:      user.Id,
	}, nil
}

func generateTokenString(claims jwt.MapClaims, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func generateJWT(user model.User, person model.Person, key string) (string, error) {
	claims := jwt.MapClaims{
		"jti":      uuid.New().String(),
		"id":       user.Id,
		"username": user.Username,
		"personId": person.Id,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	tokenString, err := generateTokenString(claims, key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
