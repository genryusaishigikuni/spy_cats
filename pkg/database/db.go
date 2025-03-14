package database

import (
	"fmt"
	"github.com/genryusaishigikuni/spy_cats/internal/target"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/postgres" // or whichever driver you use
	"gorm.io/gorm"

	"github.com/genryusaishigikuni/spy_cats/config"
	"github.com/genryusaishigikuni/spy_cats/internal/cat"
	"github.com/genryusaishigikuni/spy_cats/internal/mission"
	"github.com/genryusaishigikuni/spy_cats/internal/note"
)

// Connect opens a GORM DB connection based on the provided config.DBConfig.
func Connect(dbCfg config.DBConfig) (*gorm.DB, error) {
	// Example DSN for Postgres; adjust for MySQL, SQLite, etc.
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Name,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	return db, nil
}

// RunMigrations can do any (or both) of these:
// 1) GORM AutoMigrate
// 2) Raw SQL Migrations
func RunMigrations(db *gorm.DB) error {
	// (1) GORM AutoMigrate
	if err := autoMigrate(db); err != nil {
		return err
	}

	// (2) Run any raw SQL files in "pkg/database/migrations/"
	if err := runSQLMigrations(db, "pkg/database/migrations"); err != nil {
		return err
	}

	return nil
}

// autoMigrate uses GORM's AutoMigrate to create/modify DB tables
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&cat.Cat{},
		&mission.Mission{},
		&target.Target{},
		&note.Note{},
	)
}

// runSQLMigrations reads *.sql files from a given folder and executes them in order
func runSQLMigrations(db *gorm.DB, migrationsDir string) error {
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to read migration folder: %w", err)
	}

	for _, file := range files {
		log.Printf("Applying migration: %s", file)
		contents, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", file, err)
		}

		if err := db.Exec(string(contents)).Error; err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}
	}
	return nil
}
