package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func NewAuthentication(secret, secretAdmin string) *Auth {
	return &Auth{secretKey: []byte(secret), secretKeyAdmin: []byte(secretAdmin)}
}

type Authentication interface {
	GenerateToken(id uint) (string, error)
	ValidateToken(encodedToken string) (uint, error)
	GenerateTokenAdmin(id uint) (string, error)
	ValidateTokenAdmin(encodedToken string) (uint, error)
}

type Auth struct {
	secretKey      []byte
	secretKeyAdmin []byte
}

func (a *Auth) GenerateToken(id uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["account_id"] = id
	claims["expire"] = time.Now().Add(10 * time.Hour).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(a.secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil

}

func (a *Auth) ValidateToken(encodedToken string) (uint, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("INVALID ERROR")
		}

		return []byte(a.secretKey), nil

	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !token.Valid || !ok {
		return 0, errors.New("Unauthorized")
	}

	userID := uint(claims["account_id"].(float64))
	return userID, nil
}

func (a *Auth) GenerateTokenAdmin(id uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["admin_id"] = id
	claims["expire"] = time.Now().Add(2 * time.Hour).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(a.secretKeyAdmin)
	if err != nil {
		return "", err
	}

	return signedToken, nil

}

func (a *Auth) ValidateTokenAdmin(encodedToken string) (uint, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("INVALID ERROR")
		}

		return []byte(a.secretKeyAdmin), nil

	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !token.Valid || !ok {
		return 0, errors.New("Unauthorized")
	}

	userID := uint(claims["admin_id"].(float64))
	return userID, nil
}
