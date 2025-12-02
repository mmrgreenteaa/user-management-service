package postgresql

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/mmrgreenteaa/user-management-service/internal/auth/models"
	"gorm.io/gorm"
)

// ErrUserAgentMismatch indicates  that the user agent in the database does not match the current request.
var ErrUserAgentMismatch = errors.New("the user agent does not match")
// NoRecord indicates  about missing records
var NoRecord = errors.New("no record found")

// AddRefresh adds a new refresh token to the database.
func (db *DB) AddRefresh(refresh string, userId string, userAgent string, ip string) error {

	token := models.RefreshToken{
		Token:     GenHash(refresh),
		UserId:    userId,
		ExpiresAt: time.Now().AddDate(0, 0, 2),
		UserAgent: userAgent,
		Ip:        ip,
	}
	res := db.Create(&token)
	if res.Error != nil {

		return fmt.Errorf("faild add refresh token - %w", res.Error)
	}
	return nil
}

// RefreshTokenValid verifies the token for compliance.
// It returns ErrUserAgentMismatch if the user agent does not match.
func (db *DB) RefershTokenValid(refresh string, userAgent string, Ip string, userId string) (string, error) {
	token := models.RefreshToken{}
	genrefresh := GenHash(refresh)
	res := db.Where("token = ?", genrefresh).Where("user_id = ?", userId).First(&token)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("refresh token not found - %w", gorm.ErrRecordNotFound)
		}
		return "", fmt.Errorf("database refresh token valid query error: %w", res.Error)
	}
	if time.Now().After(token.ExpiresAt) {
		return "", fmt.Errorf("failed refresh token has expired")
	}

	if token.Ip != Ip {

		file, err := os.OpenFile("app.json.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			db.logger.Error("failed to open file logs - %w", err)
			return "", fmt.Errorf("failed to open file logs")
		}
		defer file.Close()

		fileHandler := slog.NewJSONHandler(file, nil)
		fLogger := slog.New(fileHandler)
		db.logger.Warn("unknown ip address", "user_id", userId, "refresh token", genrefresh)
		fLogger.Warn("unknown ip address", "user_id", userId, "refresh token", genrefresh)
	}

	if token.UserAgent != userAgent {
		return "", ErrUserAgentMismatch
	}

	return strconv.FormatUint(uint64(token.ID), 10), nil

}
// EditRefreshToken replaces the old refresh token with a new one.
func (db *DB) EditRefreshToken(tokenId string, newRefresh string) error {

	token := models.RefreshToken{}
	res := db.Model(&token).Where("id = ?", tokenId).Update("token", GenHash(newRefresh))
	if res.Error != nil {
		log.Printf("Ошибка при обновлении: %v", res.Error)
		return fmt.Errorf("fail update refresh token - %w", res.Error)
	} else {
		if res.RowsAffected == 0 {
			db.logger.Warn("falied edit refresh roken no record found")
			return fmt.Errorf("no row affcted ")
		}

	}
	return nil
}

// DeleteToken delete refresh token 
func (db *DB) DeleteToken(id string) error {

	res := db.Where("id = ?", id).Delete(&models.RefreshToken{})
	if res.Error != nil {

		return fmt.Errorf(" failed to delete refresh token: %w", res.Error)
	}
	if res.RowsAffected == 0 {

		return NoRecord
	}
	return nil
}
