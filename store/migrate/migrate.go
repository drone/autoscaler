package ddl

import (
	"github.com/drone/autoscaler/store/migrate/mysql"
	"github.com/drone/autoscaler/store/migrate/postgres"
	"github.com/drone/autoscaler/store/migrate/sqlite"

	"github.com/jmoiron/sqlx"
)

// Migrate performs the database migration.
func Migrate(db *sqlx.DB) error {
	switch db.DriverName() {
	case "postgres":
		return postgres.Migrate(db.DB)
	case "mysql":
		return mysql.Migrate(db.DB)
	default:
		return sqlite.Migrate(db.DB)
	}
}
