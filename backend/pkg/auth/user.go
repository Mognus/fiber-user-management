package auth

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a user with role-based access control.
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"` // Never send password in JSON
	FirstName string    `gorm:"size:100" json:"first_name"`
	LastName  string    `gorm:"size:100" json:"last_name"`
	Role      UserRole  `gorm:"type:varchar(20);default:'user'" json:"role"`
	Active    bool      `gorm:"default:true" json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName - Custom table name.
func (User) TableName() string {
	return "users"
}

// BeforeSave hook - Hash password before creating or updating user.
func (u *User) BeforeSave(tx *gorm.DB) error {
	// Only hash if password is set and not already hashed (bcrypt hashes start with $2a$).
	if u.Password != "" && len(u.Password) > 0 && u.Password[0] != '$' {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}

	// Set default role if not specified.
	if u.Role == "" {
		u.Role = RoleUser
	}

	return nil
}

// CheckPassword compares provided password with stored hash.
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// IsAdmin checks if user has admin role.
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsUser checks if user has user role.
func (u *User) IsUser() bool {
	return u.Role == RoleUser
}

// FullName returns user's full name.
func (u *User) FullName() string {
	if u.FirstName != "" && u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	return u.Email
}

// HashPassword hashes a plain text password.
func (u *User) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
