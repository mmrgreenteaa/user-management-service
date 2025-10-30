package postgresql

import (
	"errors"
	"fmt"
	"time"

	"github.com/mmrgreenteaa/user-management-service/internal/auth/models"
	"gorm.io/gorm"
)

func (db *DB) AddRefresh(refresh string, userAgent string, ip string) error {

	token := models.RefreshToken{
		Token:     GenHash(refresh),
		ExpiresAt: time.Now().AddDate(0, 0, 2),
		UserAgent: userAgent,
		Ip:        ip,
	}
	res := db.Create(&token)
	if res.Error != nil {
		return fmt.Errorf("add refresh tocken - %w", res.Error)
	}
	return nil
}

func (db *DB) RefershTokenValid(refresh string) error {
	token := models.RefreshToken{}
	genrefresh := GenHash(refresh)
	res := db.Where("token = ?", genrefresh).First(&token)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("refresh token not found or invalid: %w", res.Error)
		}
		return fmt.Errorf("database query error: %w", res.Error)
	}
	if time.Now().After(token.ExpiresAt) {
		return fmt.Errorf("token has expired")
	}

	return nil

}

func (db *DB) DeleteToken(refresh string) error {

	hashref := GenHash(refresh)
	res := db.Where("token = ?", hashref).Delete(&models.RefreshToken{})
	if res.Error != nil {
		return fmt.Errorf("fail refresh token delete: %w", res.Error)
	}
	return nil
}
