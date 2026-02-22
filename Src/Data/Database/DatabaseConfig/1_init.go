package DatabaseConfig

import (
	"RideMarket-CleanWebApi-GoLang/Data/Models"
	"fmt"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MigrateAndSeed handles database migration and initial seeding
func MigrateAndSeed() error {
	db := GetDb()
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	// 1. AutoMigrate - GORM will create or update tables automatically
	err := db.AutoMigrate(
		&Models.Country{},
		&Models.City{},
		&Models.Role{},
		&Models.User{},
		&Models.UserRole{},
	)
	if err != nil {
		return fmt.Errorf("auto migrate failed: %w", err)
	}

	// 2. Seed initial data (idempotent - safe to run multiple times)
	seedRoles(db)
	seedAdminUser(db)

	log.Println("Migration and initial seeding completed successfully")
	return nil
}

// seedRoles creates initial roles if they don't exist
func seedRoles(db *gorm.DB) {
	roles := []Models.Role{
		{
			Name:        "admin",
			DisplayName: "System Administrator",
			Description: "Full access to all features and settings",
			IsSystem:    true,
			IsDefault:   false,
			Priority:    100,
			Color:       "#FF0000",
			Icon:        "shield-lock",
			Scope:       "global",
		},
		{
			Name:        "user",
			DisplayName: "Regular User",
			Description: "Standard user access",
			IsSystem:    false,
			IsDefault:   true,
			Priority:    10,
			Color:       "#4CAF50",
			Icon:        "person",
			Scope:       "global",
		},
		{
			Name:        "support",
			DisplayName: "Support Team",
			Description: "Can view and respond to support tickets",
			IsSystem:    false,
			IsDefault:   false,
			Priority:    50,
			Color:       "#2196F3",
			Icon:        "headset",
			Scope:       "global",
		},
	}

	for _, role := range roles {
		// FirstOrCreate ensures the role is created only if it doesn't exist
		db.Where(Models.Role{Name: role.Name}).FirstOrCreate(&role)
	}
}

// seedAdminUser creates an initial admin user if it doesn't exist
func seedAdminUser(db *gorm.DB) {
	// Replace this with real hashed password in production
	// Example: hashed, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	hashedPassword := "$2a$12$...yourhashedpasswordhere..." // ← CHANGE THIS!

	admin := Models.User{
		Username:        "admin",
		Email:           "admin@example.com",
		Password:        hashedPassword,
		FullName:        "System Admin",
		IsActive:        true,
		IsEmailVerified: true,
	}

	// Create admin user only if it doesn't exist (based on Username)
	result := db.Where(Models.User{Username: admin.Username}).FirstOrCreate(&admin)
	if result.RowsAffected > 0 {
		log.Println("Admin user created successfully")
	} else {
		log.Println("Admin user already exists")
	}

	// Assign 'admin' role to the admin user
	var adminRole Models.Role
	db.Where("name = ?", "admin").First(&adminRole)

	if adminRole.ID != 0 {
		userRole := Models.UserRole{
			UserId: int(admin.ID),
			RoleId: int(adminRole.ID),
		}

		// Upsert to prevent duplicate entries
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&userRole)
	}
}
