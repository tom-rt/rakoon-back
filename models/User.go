package models

import "time"

// User ...
type User struct {
	Username  string `gorm:"username"`
	Password  string `gorm:"password"`
	Salt      string `gorm:"salt"`
	LastLogin time.Time `gorm:"last_login"`
}
