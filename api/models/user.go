package models

import (
	"gorm.io/gorm"
	"time"
)

type PermissionLevel int64

const (
	PermissionLevelEmployee PermissionLevel = iota
	PermissionLevelManager
	PermissionLevelAdmin
)

func (p PermissionLevel) String() string {
	switch p {
	case PermissionLevelEmployee:
		return "Employee"
	case PermissionLevelManager:
		return "Manager"
	case PermissionLevelAdmin:
		return "Administrator"
	}
	return "unknown"
}

type User struct {
	gorm.Model

	FirstName       string          `json:"first_name"`
	LastName        string          `json:"last_name"`
	PermissionLevel PermissionLevel `json:"permission_level"`
	Email           string          `json:"email" gorm:"uniqueIndex"`
	Password        string          `json:"-"`

	OrganizationID uint         `json:"organization_id"`
	Organization   Organization `json:"-"`
}

type Token struct {
	ID        string `json:"key" gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserID    uint
	User      User
}
