package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime:false"`
	IsDeleted bool       `gorm:"not null"`
	Name      string     `gorm:"not null"`
	Email     string     `gorm:"not null"`
	Password  string     `gorm:"not null"`
}

func NewUser(name, email, password string) *User {
	return &User{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: nil,
		IsDeleted: false,
		Name:      name,
		Email:     email,
		Password:  password,
	}
}
