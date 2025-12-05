package model

import "time"

type User struct {
	ID        int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement;not null"`
	Username  string    `json:"username" gorm:"column:username;unique;not null"`
	Password  string    `json:"-" gorm:"column:password;not null"`
	Email     string    `json:"email" gorm:"column:email"`
	Role      string    `json:"role" gorm:"column:role;default:admin"` // admin
	IsActive  int       `json:"is_active" gorm:"column:is_active;default:1"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}
