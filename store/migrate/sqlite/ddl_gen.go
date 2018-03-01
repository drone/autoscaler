package sqlite

import (
	"database/sql"
)

var migrations = []struct {
	name string
	stmt string
}{
	{
		name: "create-table-servers",
		stmt: createTableServers,
	},
	{
		name: "create-index-server-id",
		stmt: createIndexServerId,
	},
	{
		name: "create-index-server-state",
		stmt: createIndexServerState,
	},
}

// Migrate performs the database migration. If the migration fails
// and error is returned.
func Migrate(db *sql.DB) error {
	if err := createTable(db); err != nil {
		return err
	}
	completed, err := selectCompleted(db)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	for _, migration := range migrations {
		if _, ok := completed[migration.name]; ok {

			continue
		}

		if _, err := db.Exec(migration.stmt); err != nil {
			return err
		}
		if err := insertMigration(db, migration.name); err != nil {
			return err
		}

	}
	return nil
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(migrationTableCreate)
	return err
}

func insertMigration(db *sql.DB, name string) error {
	_, err := db.Exec(migrationInsert, name)
	return err
}

func selectCompleted(db *sql.DB) (map[string]struct{}, error) {
	migrations := map[string]struct{}{}
	rows, err := db.Query(migrationSelect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		migrations[name] = struct{}{}
	}
	return migrations, nil
}

//
// migration table ddl and sql
//

var migrationTableCreate = `
CREATE TABLE IF NOT EXISTS migrations (
 name VARCHAR(255)
,UNIQUE(name)
)
`

var migrationInsert = `
INSERT INTO migrations (name) VALUES (?)
`

var migrationSelect = `
SELECT name FROM migrations
`

//
// 001_create_table_servers.sql
//

var createTableServers = `
CREATE TABLE IF NOT EXISTS servers (
 server_name      TEXT PRIMARY KEY
,server_id        TEXT
,server_provider  TEXT
,server_state     TEXT
,server_image     TEXT
,server_region    TEXT
,server_size      TEXT
,server_address   TEXT
,server_capacity  INTEGER
,server_secret    TEXT
,server_error     TEXT
,server_created   INTEGER
,server_updated   INTEGER
,server_started   INTEGER
,server_stopped   INTEGER
);
`

var createIndexServerId = `
CREATE INDEX IF NOT EXISTS ix_servers_id ON servers (server_id);
`

var createIndexServerState = `
CREATE INDEX IF NOT EXISTS ix_servers_state ON servers (server_state);
`
