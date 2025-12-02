package migrations

import (
	"testing"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/mmrgreenteaa/user-management-service/internal/auth/datebase/postgresql"

	"github.com/stretchr/testify/require"
)

func TestM1224_create_rerefresh_tokens(t *testing.T) {

	cfg := postgresql.DbConfig{
		Host:       "localhost",
		Port:       "5432",
		User:       "postgres",
		Pass:       "QWERTY",
		Name:       "auth_service",
		SearchPath: "tokens_info",
	}
	db := postgresql.Connect(&cfg)

	m := gormigrate.New(db.DB, gormigrate.DefaultOptions, []*gormigrate.Migration{M1224CreateRerefreshTokens()})

	err := m.Migrate()

	require.NoError(t, err)

}
