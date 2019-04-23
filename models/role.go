package models

import "github.com/jinzhu/gorm"

// Role for user role information
type Role struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	Users       []User
}

// RoleID const defination
const (
	AdminRoleID   = 1
	UserRoleID    = 2
	MonitorRoleID = 3
)

// Roles for default insert
var (
	AdminRole = Role{
		ID:          AdminRoleID,
		Name:        "admin",
		Description: "for admin the website",
	}
	UserRole = Role{
		ID:          UserRoleID,
		Name:        "user",
		Description: "for admin the user file contents",
	}
	MonitorRole = Role{
		ID:          MonitorRoleID,
		Name:        "monitor",
		Description: "for get the basic system metrics",
	}
)

// InsertDefaultRoles insert the default of roles
func InsertDefaultRoles(db *gorm.DB) error {
	var defaultRoles = []Role{
		AdminRole, UserRole, MonitorRole,
	}
	for _, role := range defaultRoles {
		updateRoles := Role{}
		err := db.Where("id = ?", role.ID).First(&updateRoles).Error
		if err == nil {
			continue
		}
		if err != gorm.ErrRecordNotFound {

			return err
		}
		if errForCreate := db.Create(&role).Error; errForCreate != nil {
			return errForCreate
		}
	}
	return nil
}
