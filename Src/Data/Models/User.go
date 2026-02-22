package Models

import (
	"time"
)

// User represents a system user with support for authentication, RBAC, security tracking, etc.
type User struct {
	BaseModel
	Username            string     `gorm:"size:64;uniqueIndex;not null;comment:Unique username used for login"`
	Email               string     `gorm:"size:120;uniqueIndex;not null;comment:Unique email address used for login and verification"`
	Password            string     `gorm:"size:255;not null;comment:Hashed password (use bcrypt or argon2)"`
	FullName            string     `gorm:"size:100;comment:User's full name (first + last name)"`
	Phone               string     `gorm:"size:20;comment:Mobile phone number (optional)"`
	AvatarURL           string     `gorm:"size:255;comment:URL to the user's profile picture"`
	IsActive            bool       `gorm:"default:true;index;comment:Whether the account is active"`
	IsEmailVerified     bool       `gorm:"default:false;index;comment:Whether the email has been verified"`
	EmailVerifiedAt     *time.Time `gorm:"comment:Timestamp when the email was verified"`
	LastLoginAt         *time.Time `gorm:"comment:Timestamp of the last successful login"`
	LastLoginIP         string     `gorm:"size:45;comment:IP address of the last successful login (for security tracking)"`
	FailedLoginAttempts int        `gorm:"default:0;comment:Number of failed login attempts (used for brute-force protection)"`
	LockedUntil         *time.Time `gorm:"comment:Timestamp until which the account is locked (after too many failed attempts)"`

	// Many-to-many relationship with Roles (RBAC)
	UserRoles *[]UserRole
}
