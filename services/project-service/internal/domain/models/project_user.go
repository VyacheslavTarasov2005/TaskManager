package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleOwner  UserRole = "OWNER"
	RoleMember UserRole = "MEMBER"
)

type ProjectUser struct {
	UserID    uuid.UUID `gorm:"primaryKey"`
	ProjectID uuid.UUID `gorm:"primaryKey"`
	Role      UserRole  `gorm:"type:varchar(10);not null" json:"role"`
	AddedTime time.Time `gorm:"not null"`
}

func NewProjectUser(userID, projectID uuid.UUID, role string) (*ProjectUser, error) {
	if role != string(RoleOwner) && role != string(RoleMember) {
		return nil, errors.New("invalid role: must be OWNER or MEMBER")
	}
	return &ProjectUser{
		UserID:    userID,
		ProjectID: projectID,
		Role:      UserRole(role),
		AddedTime: time.Now(),
	}, nil
}
