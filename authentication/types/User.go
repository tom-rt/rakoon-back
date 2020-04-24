package types

import "time"

// User object
type User struct {
	Name  string `gorm:"name" json:"name" binding:"required"`
	Password  string `gorm:"password" json:"password" binding:"required"`
	Salt      string `gorm:"salt"`
	LastLogin time.Time `gorm:"last_login"`
	CreatedOn time.Time `gorm:"created_on"`
}
