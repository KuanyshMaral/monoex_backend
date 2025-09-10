package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"monoex_backend/internal/config"
)

var SQLDB *sql.DB

func ConnectSQL() (*sql.DB, error) {
	if SQLDB != nil {
		return SQLDB, nil
	}

	cfg := config.GetConfig()
	db, err := sql.Open(cfg.DB.Driver, cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	SQLDB = db
	log.Println("✅ Database connection established")
	return db, nil
}

// RunMigrations безопасно применяет миграции
func RunMigrations() error {
	db, err := ConnectSQL()
	if err != nil {
		return fmt.Errorf("failed to connect DB for migration: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(getMigrationsPath(), "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Проверяем текущую версию и dirty статус
	version, dirty, verr := m.Version()
	if verr != nil && verr != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", verr)
	}

	if dirty {
		// На продакшене остановить процесс
		if os.Getenv("ENV") == "production" {
			return fmt.Errorf("database is dirty at version %d. Manual intervention required", version)
		}
		// В dev можно форсировать, если ENV != production
		forceVersion := os.Getenv("MIGRATE_FORCE_VERSION")
		if forceVersion == "" {
			return fmt.Errorf("database is dirty at version %d. Set MIGRATE_FORCE_VERSION to force", version)
		}
		log.Printf("⚠️ Dirty database detected. Forcing version to %s in dev environment\n", forceVersion)
		if err := m.Force(int(version)); err != nil {
			return fmt.Errorf("failed to force migration version: %w", err)
		}
	}

	// Применяем миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("✅ Migrations applied successfully")
	return nil
}

// getMigrationsPath возвращает путь к папке migrations в формате file:///
func getMigrationsPath() string {
	path := "./database/migrations"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("Migrations folder not found: %s", path)
	}

	return "file://" + path
}
