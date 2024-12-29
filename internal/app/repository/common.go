package repository

import (
	"github.com/jackc/pgx/v5/pgconn"
)

func IsPostgresError(err error, code string) bool {
	if err == nil {
		return false
	}
	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return false
	}
	return code == pgErr.Code
}
