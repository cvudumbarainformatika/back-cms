package seeders

import (
	"encoding/json"
	"fmt"
	"log"

	models "github.com/cvudumbarainformatika/backend/app/Models"
	"github.com/jmoiron/sqlx"
)

// SeedMenus seeds the database with initial navigation menus
func SeedMenus(db *sqlx.DB) error {
	log.Println("🌱 Seeding menus...")

	// Check if menus already exist
	var count int
	err := db.Get(&count, `SELECT COUNT(*) FROM menus`)
	if err != nil {
		return fmt.Errorf("failed to check existing menus: %w", err)
	}

	if count > 0 {
		log.Println("⊘ Menus already exist, skipping seed")
		return nil
	}

	// Define roles JSON array
	allRoles := []string{"public", "member", "admin_cabang", "admin_wilayah", "admin_pusat"}
	rolesJSON, _ := json.Marshal(allRoles)

	// Define initial menus
	menus := []struct {
		label    string
		slug     string
		to       string
		icon     string
		parentID *int64
		position string
		order    int
		isActive bool
		isFixed  bool
		roles    string
	}{
		{
			label:    "Beranda",
			slug:     "beranda",
			to:       "/",
			icon:     "i-lucide-home",
			parentID: nil,
			position: "header",
			order:    1,
			isActive: true,
			isFixed:  true,
			roles:    string(rolesJSON),
		},
		{
			label:    "Berita",
			slug:     "berita",
			to:       "/berita",
			icon:     "i-lucide-newspaper",
			parentID: nil,
			position: "header",
			order:    3,
			isActive: true,
			isFixed:  true,
			roles:    string(rolesJSON),
		},
		{
			label:    "Agenda",
			slug:     "agenda",
			to:       "/agenda",
			icon:     "i-lucide-calendar",
			parentID: nil,
			position: "header",
			order:    4,
			isActive: true,
			isFixed:  true,
			roles:    string(rolesJSON),
		},
		{
			label:    "Direktori",
			slug:     "direktori",
			to:       "/direktori",
			icon:     "i-lucide-map-pin",
			parentID: nil,
			position: "header",
			order:    5,
			isActive: true,
			isFixed:  true,
			roles:    string(rolesJSON),
		},
	}

	// Create menus
	for _, menuData := range menus {
		menu := &models.Menu{
			Label:    menuData.label,
			Slug:     menuData.slug,
			To:       menuData.to,
			Icon:     menuData.icon,
			ParentID: menuData.parentID,
			Position: menuData.position,
			Order:    menuData.order,
			IsActive: menuData.isActive,
			IsFixed:  menuData.isFixed,
			Roles:    menuData.roles,
		}

		if err := menu.Create(db); err != nil {
			return fmt.Errorf("failed to create menu %s: %w", menuData.label, err)
		}

		log.Printf("✓ Created menu: %s (order: %d, position: %s)", menu.Label, menu.Order, menu.Position)
	}

	log.Println("✓ Menu seeding completed")
	return nil
}
