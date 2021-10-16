package admin

type Admin struct {
	ID       uint   `db:"id" json:"-"`
	Name     string `db:"name" json:"-"`
	Email    string `db:"email" validate:"contains=@rahp.com"`
	Password string `db:"password" validate:"required,min=7,endswith=rahp"`
	JobID    uint   `db:"jobdesk_id" json:"job_id"`
	Secret   string `db:"secret"`
}

type InputAdmin struct {
	Name            string `validate:"contains=rahp,required"`
	Email           string `db:"email" validate:"contains=@rahp.com"`
	Password        string `db:"password" validate:"required,min=7,endswith=rahp"`
	ConfirmPassword string `validate:"eqfield=Password,required"`
}
