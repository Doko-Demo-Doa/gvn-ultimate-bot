package discordrepos

import (
	"doko/gvn-ultimate-bot/models"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// newTestDB spins up an isolated in-memory SQLite database, migrated with the
// models these repo tests exercise. Each call gets its own fresh database.
func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite db: %v", err)
	}

	if err := db.AutoMigrate(
		&models.DiscordUser{},
		&models.DiscordRole{},
		&models.DiscordMessageAuditLog{},
	); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	return db
}
