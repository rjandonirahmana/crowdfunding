package account

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

type RepositoryUser interface {
	Save(user User) (uint, error)
	FindByEmail(email string) (User, error)
	FindByID(ID int) (User, error)
	UpdateUser(user User) error
	IsEmailAvailable(email string) (bool, error)
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(user User) (uint, error) {
	querry := `	INSERT INTO 
	users
	(
		name, email, occupation, password, role, created_at, updated_at, salt, avatar
	) 
	VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id
	`
	var id uint
	err := r.db.QueryRowx(querry, user.Name, user.Email, user.Occupation, user.PasswordHash, user.Role, time.Now(), time.Now(), user.Salt, user.Avatar).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil

}

func (r *Repository) FindByEmail(email string) (User, error) {
	querry := ` 
	SELECT * FROM users WHERE email = $1`

	var user User
	err := r.db.Get(&user, querry, email)
	if err != nil {
		return User{}, fmt.Errorf("email not found or error :%v", err)
	}

	return user, nil
}

func (r *Repository) FindByID(ID int) (User, error) {
	querry := `
	SELECT * FROM users WHERE id = $1
	`
	var user User
	err := r.db.Get(&user, querry, ID)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *Repository) UpdateUser(user User) error {

	querry := `
			UPDATE
				users
			SET 
				avatar = $1
			WHERE
				id = $2
		`

	_, err := r.db.Exec(querry, user.Avatar, user.ID)
	if err != nil {
		return err
	}

	return nil

}

func (r *Repository) IsEmailAvailable(email string) (bool, error) {
	querry := ` 
	SELECT 
		id
	FROM 
		users
	WHERE 
		email = $1
	`

	var id int
	err := r.db.Get(&id, querry, email)
	if err != sql.ErrNoRows || id != 0 {
		return false, fmt.Errorf("error %v or email has been used", err)
	}

	return true, nil
}
