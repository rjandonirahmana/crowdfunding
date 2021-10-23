package admin

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

type Repository interface {
	CreateAdmin(admin Admin) (uint, error)
	FindAdmin(id uint) (Admin, error)
	IsEmailAvailable(email string) (bool, error)
}

func NewRepositoryAdmin(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) CreateAdmin(admin Admin) (uint, error) {
	querry := `INSERT INTO admin (name, email, password, jobdesk_id, secret) VALUES($1, $2, $3, $4, $5) RETURNING id`

	var id uint
	err := r.db.QueryRowx(querry, admin.Name, admin.Email, admin.Password, admin.JobID, admin.Secret).Scan(&id)

	if err != nil {
		fmt.Println("error disini ga?")
		return 0, err
	}
	return id, nil
}

func (r *repository) FindAdmin(id uint) (Admin, error) {
	querry := `SELECT * FROM admin WHERE id = $1`

	var admin Admin
	err := r.db.Get(&admin, querry, id)
	if err != nil {
		return admin, err
	}

	return admin, nil

}

func (r *repository) IsEmailAvailable(email string) (bool, error) {
	querry := `SELECT id FROM admin WHERE email = $1`

	var id *int
	err := r.db.Get(&id, querry, email)
	if err != sql.ErrNoRows || *id != 0 {
		return false, fmt.Errorf("error %v or email has been used", err)
	}

	return true, nil
}
