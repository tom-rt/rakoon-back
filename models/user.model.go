package models

import (
	"rakoon/rakoon-back/db"
	"time"

	"github.com/jmoiron/sqlx"
)

// User object
type User struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name" binding:"required"`
	Password  string    `db:"password" json:"password" binding:"required"`
	Salt      string    `db:"salt" json:"salt"`
	Reauth    bool      `db:"reauth" json:"reauth"`
	LastLogin time.Time `db:"last_login" json:"last_login"`
	CreatedOn time.Time `db:"created_on" json:"created_on"`
}

// UserPublic object
type UserPublic struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name" binding:"required"`
	Reauth    bool      `db:"reauth" json:"reauth"`
	LastLogin time.Time `db:"last_login" json:"last_login"`
	CreatedOn time.Time `db:"created_on" json:"created_on"`
}

// UserID obj
type UserID struct {
	ID string `db:"id" json:"id" binding:"required"`
}

// UserUpdate input for user updated values
type UserUpdate struct {
	ID   string `db:"id" json:"id" binding:"required"`
	Name string `db:"name" json:"name" binding:"required"`
}

// GetUserByName func
func GetUserByName(name string) (User, error) {
	var user User
	err := db.DB.Get(&user,
		"SELECT id, name, password, salt, reauth, created_on, last_login FROM users where name = $1",
		name)
	return user, err
}

// GetUserByID func model
func GetUserByID(ID string) (User, error) {
	var user User
	err := db.DB.Get(&user,
		"SELECT id, name, password, salt, reauth, created_on, last_login FROM users where id = $1",
		ID)
	return user, err
}

// GetUserPublicByID func model
func GetUserPublicByID(ID string) (UserPublic, error) {
	var user UserPublic
	err := db.DB.Get(&user,
		"SELECT id, name, reauth, created_on, last_login FROM users where id = $1",
		ID)
	return user, err
}

// SetReauthByName func
func SetReauthByName(userName string, value bool) *sqlx.Row {
	ret := db.DB.QueryRowx("UPDATE users SET reauth = $1 WHERE name = $2", value, userName)
	return ret
}

// SetReauthByID func
func SetReauthByID(ID string, value bool) *sqlx.Row {
	ret := db.DB.QueryRowx("UPDATE users SET reauth = $1 WHERE id = $2", value, ID)
	return ret
}

// CreateUser function
func CreateUser(user User) {
	tx := db.DB.MustBegin()
	tx.MustExec("INSERT INTO users (name, password, salt, reauth) VALUES ($1, $2, $3, $4)", user.Name, user.Password, user.Salt, user.Reauth)
	tx.Commit()
}

// CreateUser function
func UpdateUser(update UserUpdate) {
	db.DB.Queryx("UPDATE users SET name = $1 WHERE id = $2", update.Name, update.ID)
}

// DeleteUser function
func DeleteUser(ID string) {
	tx := db.DB.MustBegin()
	tx.MustExec("DELETE FROM users where id = $1", ID)
	tx.Commit()
}
