package models

import (
	"rakoon/rakoon-back/db"
	"time"

	"github.com/jmoiron/sqlx"
)

//User object
type User struct {
	Name      string    `db:"name" gorm:"name" json:"name" binding:"required"`
	Password  string    `db:"password" gorm:"password" json:"password" binding:"required"`
	Salt      string    `db:"salt" gorm:"salt"`
	Reauth    bool      `db:"reauth" gorm:"reauth"`
	LastLogin time.Time `db:"last_login" gorm:"last_login"`
	CreatedOn time.Time `db:"created_on" gorm:"created_on"`
}

//Username obj
type Username struct {
	Name string `db:"name" gorm:"name" json:"name" binding:"required"`
}

//GetUser func
func GetUser(name string) (User, error) {
	var user User
	err := db.DB2.Get(&user,
		"SELECT id, name, password, salt, reauth, created_on, last_login FROM users where name = $1",
		name)
	return user, err
}

//SetReauth func
func SetReauth(userName string, value bool) *sqlx.Row {
	ret := db.DB2.QueryRowx("UPDATE users SET reauth = $1 WHERE name = $2", value, userName)
	return ret
}

func CreateUser(user User) {
	tx := db.DB2.MustBegin()
	tx.MustExec("INSERT INTO users (name, password, salt, reauth) VALUES ($1, $2, $3, $4)", user.Name, user.Password, user.Salt, user.Reauth)
	tx.Commit()
}
