package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"funding/model"

	"github.com/jmoiron/sqlx"
)

type adminRepo struct {
	db *sqlx.DB
}

type RepositoryAdmin interface {
	CreateAdmin(admin *model.Admin) (*model.Admin, error)
	FindAdmin(id uint) (model.Admin, error)
	IsEmailAvailable(email *string) error
}

func NewRepositoryAdmin(db *sqlx.DB) *adminRepo {
	return &adminRepo{db: db}
}

func (r *adminRepo) CreateAdmin(admin *model.Admin) (*model.Admin, error) {
	querry := `INSERT INTO admin (name, email, password, jobdesk_id, secret) VALUES($1, $2, $3, $4, $5) RETURNING id`

	err := r.db.QueryRowx(querry, admin.Name, admin.Email, admin.Password, admin.JobID, admin.Secret).Scan(&admin.ID)

	if err != nil {
		fmt.Println("error disini ga?")
		return nil, err
	}
	return admin, nil
}

func (r *adminRepo) FindAdmin(id uint) (model.Admin, error) {
	querry := `SELECT * FROM admin WHERE id = $1`

	var admin model.Admin
	err := r.db.Get(&admin, querry, id)
	if err != nil {
		return admin, err
	}

	return admin, nil

}

func (r *adminRepo) IsEmailAvailable(email *string) error {
	var id uint
	querry := `SELECT id FROM admin WHERE email = $1`
	err := r.db.QueryRowx(querry, email).Scan(&id)
	if err == sql.ErrNoRows && id == 0 {
		return nil
	}
	if err != nil {
		return err
	}

	return errors.New("email has been used by another")
}
