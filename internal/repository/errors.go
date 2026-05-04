package repository

import (
	"errors"

	moderncsqlite "modernc.org/sqlite"
	lib "modernc.org/sqlite/lib"
)

// isUniqueConstraintErr checks if the given error is a SQLite unique constraint violation error.
func isUniqueConstraintErr(err error) bool {
	var sqliteErr *moderncsqlite.Error
	if errors.As(err, &sqliteErr) {
		return sqliteErr.Code() == lib.SQLITE_CONSTRAINT_UNIQUE
	}
	return false
}
