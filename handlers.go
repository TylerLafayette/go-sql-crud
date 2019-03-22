package sqlcrud

import "database/sql"

// Options contains a configuration for the handlers.
type Options struct {
	Database *sql.DB
}
