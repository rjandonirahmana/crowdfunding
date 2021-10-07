package account

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
	FindByID(id int) (User, error)
	FindByEmail(email string) (User, error)
}

type service struct {
	repository RepositoryUser
}

func NewService(repository RepositoryUser) *service {
	return &service{repository}
}

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz123456789"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {

	ava, err := s.repository.IsEmailAvailable(input.Email)
	if !ava || err != nil {
		return User{}, fmt.Errorf("errors : %v", err)
	}
	user := User{}
	user.Salt = RandStringBytes(len(input.Password) + 10)

	h := sha256.New()
	h.Write([]byte(input.Password + user.Salt))

	id, err := s.repository.LastID()
	if err != nil {
		return User{}, err
	}
	user.ID = id + 1
	user.PasswordHash = fmt.Sprintf("%X", h.Sum(nil))
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation
	user.CreatedAt = time.Now()
	user.Role = "user"
	user.UpdatedAt = time.Now()

	err = s.repository.Save(user)
	if err != nil {
		return User{}, err
	}

	return user, nil

}

func (s *service) Login(input LoginInput) (User, error) {
	account, err := s.repository.FindByEmail(input.Email)
	if err != nil {
		return User{}, err
	}
	password := input.Password + account.Salt
	h := sha256.New()
	h.Write([]byte(password))
	password = fmt.Sprintf("%X", h.Sum(nil))

	if password != account.PasswordHash {
		return User{}, fmt.Errorf("password incorrect")
	}

	return account, nil

}

func (s *service) FindByID(id int) (User, error) {
	user, err := s.repository.FindByID(id)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *service) FindByEmail(email string) (User, error) {
	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
