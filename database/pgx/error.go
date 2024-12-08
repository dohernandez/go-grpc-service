package pgx

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
)

type sqlState interface {
	SQLState() string
}

// IsUniqueViolation checks whether the error is unique violation.
func IsUniqueViolation(err error) bool {
	pgErr, ok := err.(sqlState)
	if !ok {
		return false
	}

	return pgErr.SQLState() == pgerrcode.UniqueViolation
}

// IsForeignKeyViolation checks whether the error is foreign key violation.
func IsForeignKeyViolation(err error) bool {
	pgErr, ok := err.(sqlState)
	if !ok {
		return false
	}

	return pgErr.SQLState() == pgerrcode.ForeignKeyViolation
}

// IsNoRows checks whether the error is sql.ErrNoRows.
func IsNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
