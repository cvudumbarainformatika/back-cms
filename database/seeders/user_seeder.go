package seeders

import (
	"fmt"
	"log"

	models "github.com/cvudumbarainformatika/backend/app/Models"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// SeedUsers seeds the database with initial users
func SeedUsers(db *sqlx.DB) error {
	log.Println("🌱 Seeding users...")

	// Check if admin user already exists
	existingAdmin, err := models.FindByEmail(db, "admin@pdpi.co.id")
	if err != nil {
		return fmt.Errorf("failed to check admin user: %w", err)
	}

	if existingAdmin != nil {
		log.Println("⊘ Admin user already exists, skipping seed")
		return nil
	}

	// Define initial users
	users := []struct {
		name     string
		email    string
		password string
		role     string
	}{
		{
			name:     "Admin Pusat",
			email:    "admin@pdpi.co.id",
			password: "AdminPDPI2024!",
			role:     "admin_pusat",
		},
		{
			name:     "Admin Wilayah",
			email:    "admin.wilayah@pdpi.co.id",
			password: "AdminWilayah2024!",
			role:     "admin_wilayah",
		},
		{
			name:     "Admin Cabang",
			email:    "admin.cabang@pdpi.co.id",
			password: "AdminCabang2024!",
			role:     "admin_cabang",
		},
		{
			name:     "Member Test",
			email:    "member@pdpi.co.id",
			password: "Member2024!",
			role:     "member",
		},
	}

	// Create users
	for _, userData := range users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password for %s: %w", userData.email, err)
		}

		user := &models.User{
			Name:     userData.name,
			Email:    userData.email,
			Password: string(hashedPassword),
			Role:     userData.role,
			Status:   "active",
		}

		if err := user.Create(db); err != nil {
			return fmt.Errorf("failed to create user %s: %w", userData.email, err)
		}

		log.Printf("✓ Created user: %s (%s) - Email: %s", user.Name, user.Role, user.Email)
	}

	log.Println("✓ User seeding completed")
	return nil
}
