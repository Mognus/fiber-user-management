package auth

import "time"

// UserRole represents the role of a user in the system.
type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
	RoleGuest UserRole = "guest"
)

// Role represents a role that can be assigned to users.
type Role struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;uniqueIndex;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName - Custom table name.
func (Role) TableName() string {
	return "roles"
}

func defaultRoles() []Role {
	roles := make([]Role, 0, 3)
	for _, name := range roleNames() {
		roles = append(roles, Role{Name: name})
	}
	return roles
}

func roleNames() []string {
	return []string{
		string(RoleAdmin),
		string(RoleUser),
		string(RoleGuest),
	}
}
