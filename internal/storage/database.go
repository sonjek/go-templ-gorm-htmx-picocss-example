package storage

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite" //nolint:blank-imports
)

const defaultDSN = "file:sqlite.db?cache=shared&mode=rwc"

func NewDbStorage() (*gorm.DB, error) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = defaultDSN
	}
	return openDb(dsn)
}

func NewInMemoryDbStorage() (*gorm.DB, error) {
	return openDb(":memory:")
}

func openDb(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.New(sqlite.Config{
		DSN:        dsn,
		DriverName: "sqlite",
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func DBMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Note{})
}

func SeedData(db *gorm.DB) error {
	var count int64
	if err := db.Unscoped().Model(&Note{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		for i := range NotesSeed {
			NotesSeed[i].UpdatedAt = NotesSeed[i].CreatedAt
		}
		return db.Create(&NotesSeed).Error
	}
	return nil
}
