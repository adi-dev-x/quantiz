package migrations

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"log"
	"myproject/pkg/config"
)

func RunMigrations(cnf config.Config) error {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Jakarta", cnf.PGHost, cnf.PGUserName, cnf.PGPassword, cnf.PGDBName, cnf.PgPort)
	log.Println("this is the database for migration ", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get *sql.DB:", err)
	}
	defer sqlDB.Close()
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("could not set goose dialect: %w", err)
	}

	if err := applyMigrations(sqlDB, "pkg/migrations/backups"); err != nil {
		return fmt.Errorf("could not apply  migrations: %v", err)
	}

	return nil
}
func applyMigrations(db *sql.DB, migrationPath string) error {
	if err := goose.Up(db, migrationPath); err != nil {
		return fmt.Errorf("could not apply migrations from %s: %v", migrationPath, err)
	}
	return nil
}
