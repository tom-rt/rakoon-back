package types

import "time"

// Database User object
type User struct {
	Name      string    `gorm:"name" json:"name" binding:"required"`
	Password  string    `gorm:"password" json:"password" binding:"required"`
	Salt      string    `gorm:"salt"`
	Reauth    bool      `gorm:"reauth"`
	LastLogin time.Time `gorm:"last_login"`
	CreatedOn time.Time `gorm:"created_on"`
}

type Logout struct {
	Name string `gorm:"name" json:"name" binding:"required"`
}
