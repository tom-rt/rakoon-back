package models

import (
	"rakoon/rakoon-back/db"
	"time"
)

// User object
type User struct {
	ID         int       `db:"id" json:"id"`
	Name       string    `db:"name" json:"name" binding:"required"`
	Password   string    `db:"password" json:"password" binding:"required"`
	Salt       string    `db:"salt" json:"salt"`
	Reauth     bool      `db:"reauth" json:"reauth"`
	CreatedOn  time.Time `db:"created_on" json:"created_on"`
	LastLogin  time.Time `db:"last_login" json:"last_login"`
	ArchivedOn time.Time `db:"archived_on" json:"archived_on"`
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

// UserPassword input for user's password
type UserPassword struct {
	ID       string `db:"id" json:"id" binding:"required"`
	Password string `db:"password" json:"password" binding:"required"`
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

// RefreshUserConnection func
func RefreshUserConnection(userName string, value bool) {
	tx := db.DB.MustBegin()
	tx.MustExec("UPDATE users SET reauth = $1, last_login = now() WHERE name = $2", value, userName)
	tx.Commit()
}

// UpdateUserPassword func
func UpdateUserPassword(user UserPassword, salt string) {
	tx := db.DB.MustBegin()
	tx.MustExec("UPDATE users SET password = $1, salt = $2, reauth = true WHERE id = $3", user.Password, salt, user.ID)
	tx.Commit()
}

// SetReauthByName func
func SetReauthByName(userName string, value bool) {
	tx := db.DB.MustBegin()
	tx.MustExec("UPDATE users SET reauth = $1 WHERE name = $2", value, userName)
	tx.Commit()
}

// SetReauthByID func
func SetReauthByID(ID string, value bool) {
	tx := db.DB.MustBegin()
	tx.MustExec("UPDATE users SET reauth = $1 WHERE id = $2", value, ID)
	tx.Commit()
}

// CreateUser function
func CreateUser(user User) {
	tx := db.DB.MustBegin()
	tx.MustExec("INSERT INTO users (name, password, salt, reauth) VALUES ($1, $2, $3, $4)", user.Name, user.Password, user.Salt, user.Reauth)
	tx.Commit()
}

// UpdateUser function
func UpdateUser(update UserUpdate) {
	tx := db.DB.MustBegin()
	tx.MustExec("UPDATE users SET name = $1 WHERE id = $2", update.Name, update.ID)
	tx.Commit()
}

// ArchiveUser function
func ArchiveUser(user UserID) {
	tx := db.DB.MustBegin()
	tx.MustExec("UPDATE users SET archived_on = now(), reauth = true WHERE id = $1", user.ID)
	tx.Commit()
}

// DeleteUser function
func DeleteUser(ID string) {
	tx := db.DB.MustBegin()
	tx.MustExec("DELETE FROM users where id = $1", ID)
	tx.Commit()
}
