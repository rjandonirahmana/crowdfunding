package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"funding/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type repositoryUser struct {
	db *sqlx.DB
}

type RepositoryUser interface {
	Save(user model.User) (*model.User, error)
	FindByEmail(email *string) (*model.User, error)
	FindByID(ID uint) (*model.User, error)
	// UpdateUser(user model.User) error
	IsEmailAvailable(email *string) error
}

func NewRepositoryUser(db *sqlx.DB) *repositoryUser {
	return &repositoryUser{db: db}
}

func (r *repositoryUser) Save(user model.User) (*model.User, error) {
	querry := `	INSERT INTO 
	users
	(
		name, email, occupation, password, role, created_at, updated_at, salt, avatar
	) 
	VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id
	`
	err := r.db.QueryRowx(querry, user.Name, user.Email, user.Occupation, user.PasswordHash, user.Role, time.Now(), time.Now(), user.Salt, user.Avatar).Scan(&user.ID)

	if err != nil {
		return &user, err
	}

	return &user, nil

}

func (r *repositoryUser) FindByEmail(email *string) (*model.User, error) {
	querry := ` 
	SELECT * FROM users WHERE email = $1`

	var user model.User
	err := r.db.Get(&user, querry, email)
	if err != nil {
		return &user, fmt.Errorf("email not found or error :%v", err)
	}

	return &user, nil
}

func (r *repositoryUser) FindByID(ID uint) (*model.User, error) {
	querry := `
	SELECT * FROM users WHERE id = $1
	`
	var user model.User
	err := r.db.Get(&user, querry, ID)

	if err != nil {
		return &model.User{}, err
	}

	return &user, nil
}

// func (r *repositoryUser) UpdateUser(user model.User) error {

// 	querry := `
// 			UPDATE
// 				users
// 			SET
// 				avatar = $1
// 			WHERE
// 				id = $2
// 		`

// 	_, err := r.db.Exec(querry, user.Avatar, user.ID)
// 	if err != nil {
// 		return err
// 	}

// 	return nil

// }

func (r *repositoryUser) IsEmailAvailable(email *string) error {
	var id uint
	querry := `SELECT id FROM users WHERE email = $1`
	err := r.db.QueryRowx(querry, email).Scan(&id)
	if err == sql.ErrNoRows && id == 0 {
		return nil
	}
	if err != nil {
		return err
	}

	return errors.New("email has been used by another")
}
