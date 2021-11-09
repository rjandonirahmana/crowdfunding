package usecase

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"funding/model"
	"funding/repository"
	"log"
	"math/rand"
)

type serviceAdmin struct {
	repository repository.RepositoryAdmin
}

type ServiceAdmin interface {
	Register(input model.InputAdmin) (model.Admin, error)
}

func NewServiceAdmin(repo repository.RepositoryAdmin) *serviceAdmin {
	return &serviceAdmin{repository: repo}
}

func RandStringRunes(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		log.Print(err)
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(b)
}

func (s *serviceAdmin) Register(input model.InputAdmin) (model.Admin, error) {

	err := s.repository.IsEmailAvailable(input.Email)
	if err != nil {
		return model.Admin{}, err
	}

	secret := RandStringRunes(10)
	h := sha256.New()
	h.Write([]byte(input.Password + secret))
	password := fmt.Sprintf("%X", h.Sum([]byte(secret)))
	admin := model.Admin{
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
