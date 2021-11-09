package model

import "time"

type User struct {
	ID           uint      `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Occupation   string    `db:"occupation" json:"occupation"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password" json:"-"`
	Avatar       string    `db:"avatar" json:"avatar"`
	Salt         string    `db:"salt" json:"-"`
	Role         string    `db:"role" json:"role"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type RegisterUserInput struct {
	Name            string `json:"name" validate:"required"`
	Occupation      string `json:"occupation" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6,max=15"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CheckEmailInput struct {
	Email string `json:"email" binding:"required,email"`
}
