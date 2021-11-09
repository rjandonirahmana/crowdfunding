package model

type Admin struct {
	ID       uint   `db:"id" json:"-"`
	Name     string `db:"name" json:"-"`
	Email    string `db:"email" validate:"contains=@rahp.com"`
	Password string `db:"password" validate:"required,min=7,endswith=rahp"`
	JobID    uint   `db:"jobdesk_id" json:"job_id"`
	Secret   string `db:"secret"`
}

type InputAdmin struct {
	Name            string `json:"name" validate:"containsany=rahp,required"`
	Email           string `json:"email" validate:"containsany=@rahp.com"`
	Password        string `json:"password" db:"password" validate:"required,min=7,endswith=rahp"`
	ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password,required"`
	Job_ID          uint   `json:"job_id" validate:"required"`
}
