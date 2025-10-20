package postgresql

import (
	"fmt"
	"time"

	"github.com/mmrgreenteaa/user-management-service/internal/auth/models"
)

func (db *DB) AddRefresh(refresh string, userAgent string,  ip string) error {

	token := models.RefreshToken{
		Token:     refresh,
		ExpiresAt: time.Now().AddDate(0, 0, 2),
		UserAgent: userAgent,
		Ip: ip,
	}
	res := db.Create(&token)
	if res.Error != nil {
		return fmt.Errorf("add refresh tocken - %w", res.Error)
	}
	return nil
}
