package seeders

import (
	"log"

	"github.com/jmoiron/sqlx"
)

// RunSeeders runs all database seeders
func RunSeeders(db *sqlx.DB) error {
	log.Println("")
	log.Println("╔════════════════════════════════════════════════════════════╗")
	log.Println("║              DATABASE SEEDING STARTED                      ║")
	log.Println("╚════════════════════════════════════════════════════════════╝")
	log.Println("")

	seeders := []struct {
		name string
		run  func(*sqlx.DB) error
	}{
		{
			name: "Users",
			run:  SeedUsers,
		},
		{
			name: "Menus",
			run:  SeedMenus,
		},
		// Add more seeders here as needed
		// {
		//     name: "Berita",
		//     run:  SeedBerita,
		// },
	}

	for _, seeder := range seeders {
		log.Printf("Running seeder: %s", seeder.name)
		if err := seeder.run(db); err != nil {
			log.Printf("✗ Seeder %s failed: %v", seeder.name, err)
			return err
		}
		log.Println("")
	}

	log.Println("╔════════════════════════════════════════════════════════════╗")
	log.Println("║              ✓ ALL SEEDERS COMPLETED                      ║")
	log.Println("╚════════════════════════════════════════════════════════════╝")
	log.Println("")

	return nil
}
