package sqlcrud

import "database/sql"

// Init contains a configuration for the database connection.
type Init struct {
	Database *sql.DB
}
