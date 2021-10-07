package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

func NewAuthentication(secret []byte) *Auth {
	return &Auth{secretKey: secret}
}

type Authentication interface {
	GenerateToken(id int) (string, error)
	ValidateToken(encodedToken string) (int, error)
}

type Auth struct {
	secretKey []byte
}

func (a *Auth) GenerateToken(id int) (string, error) {
	claims := jwt.MapClaims{}
	claims["account_id"] = id

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(a.secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil

}

func (a *Auth) ValidateToken(GetToken string) (int, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(GetToken, claims, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("INVALID ERROR")
		}

		return a.secretKey, nil
	})

	if err != nil {
		return 0, err
	}

	id := int(claims["account_id"].(float64))
	return id, nil

}
