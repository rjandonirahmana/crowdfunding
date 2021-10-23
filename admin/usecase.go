package admin

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
)

type service struct {
	repository Repository
}

type Service interface {
	Register(input InputAdmin) (Admin, error)
}

func NewServiceAdmin(repo Repository) *service {
	return &service{repository: repo}
}

var letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func RandStringRunes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (s *service) Register(input InputAdmin) (Admin, error) {

	ava, err := s.repository.IsEmailAvailable(input.Email)
	if !ava || err != nil {
		return Admin{}, err
	}

	secret := RandStringRunes(10)
	h := sha256.New()
	h.Write([]byte(input.Password + secret))
	password := fmt.Sprintf("%X", h.Sum([]byte(secret)))
	admin := Admin{
		Name:     input.Name,
		Email:    input.Email,
		Secret:   secret,
		Password: password,
		JobID:    input.Job_ID,
	}

	id, err := s.repository.CreateAdmin(admin)
	if err != nil {
		return admin, err
	}

	admin.ID = id
	return admin, nil
}
