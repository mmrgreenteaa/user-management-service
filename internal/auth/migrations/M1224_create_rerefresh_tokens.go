package migrations

import (
	"fmt"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func M1224CreateRerefreshTokens() *gormigrate.Migration {

	type RefreshToken struct {
		ID        uint   `gorm:"primaryKey"`
		Token     string `gorm:"unique;not null"`
		UserAgent string
		Ip        string
		UserID    string
		CreatedAt time.Time
		ExpiresAt time.Time
	}

	return &gormigrate.Migration{
		ID: "2025159140000",
		Migrate: func(tx *gorm.DB) error {

		
			err := tx.AutoMigrate(RefreshToken{})
			if err != nil {
				return fmt.Errorf("magrate error %w", err)
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&RefreshToken{})
		},
	}
}
