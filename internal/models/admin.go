package models

type Admin struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"` // хранить хэш
}
