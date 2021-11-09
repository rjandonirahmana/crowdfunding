package usecase

import (
	"crypto/sha256"
	"fmt"
	"funding/model"
	"funding/repository"
	"time"
)

type ServiceUser interface {
	RegisterUser(input model.RegisterUserInput) (*model.User, error)
	Login(input model.LoginInput) (*model.User, error)
	FindByID(id uint) (*model.User, error)
}

type serviceUser struct {
	repository repository.RepositoryUser
}

func NewService(repository repository.RepositoryUser) *serviceUser {
	return &serviceUser{repository: repository}
}

func (s *serviceUser) RegisterUser(input model.RegisterUserInput) (*model.User, error) {

	err := s.repository.IsEmailAvailable(input.Email)
	if err != nil {
		return &model.User{}, fmt.Errorf("errors : %v", err)
	}
	user := model.User{}
	user.Salt = RandStringRunes(10)

	h := sha256.New()
	h.Write([]byte(input.Password + user.Salt))

	user.PasswordHash = fmt.Sprintf("%X", h.Sum([]byte(input.Email)))
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation
	user.CreatedAt = time.Now()
	user.Role = "user"
	user.UpdatedAt = time.Now()

	usernew, err := s.repository.Save(user)
	if err != nil {
		return &user, err
	}

	return usernew, nil

}

func (s *serviceUser) Login(input model.LoginInput) (*model.User, error) {
	account, err := s.repository.FindByEmail(input.Email)
	if err != nil {
		return &model.User{}, err
	}
	password := input.Password + account.Salt
	h := sha256.New()
	h.Write([]byte(password))
	password = fmt.Sprintf("%X", h.Sum([]byte(account.Email)))

	if password != account.PasswordHash {
		return &model.User{}, fmt.Errorf("password incorrect")
	}

	return account, nil

}

func (s *serviceUser) FindByID(id uint) (*model.User, error) {
	user, err := s.repository.FindByID(id)
	if err != nil {
		return user, err
	}
	return user, nil
}
