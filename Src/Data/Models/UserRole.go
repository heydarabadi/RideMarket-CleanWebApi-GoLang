package Models

import "time"

type UserRole struct {
	UserID uint `gorm:"primaryKey;column:user_id"`
	RoleID uint `gorm:"primaryKey;column:role_id"`

	AssignedAt   time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	AssignedByID uint       `gorm:"index"`
	IsActive     bool       `gorm:"default:true"`
	ExpiresAt    *time.Time `gorm:"default:null"`
}
