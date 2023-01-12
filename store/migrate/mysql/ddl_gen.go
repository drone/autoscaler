package mysql

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
	{
		name: "alter-table-servers-add-column-server-lastbusy",
		stmt: alterTableServersAddColumnServerLastbusy,
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
CREATE TABLE servers (
 server_name      VARCHAR(50) PRIMARY KEY
,server_id        VARCHAR(250)
,server_provider  VARCHAR(50)
,server_state     VARCHAR(50)
,server_image     VARCHAR(250)
,server_region    VARCHAR(50)
,server_size      VARCHAR(50)
,server_platform  VARCHAR(50)
,server_address   VARCHAR(250)
,server_capacity  INTEGER
,server_secret    VARCHAR(50)
,server_error     BLOB
,server_ca_key    BLOB
,server_ca_cert   BLOB
,server_tls_key   BLOB
,server_tls_cert  BLOB
,server_created   INTEGER
,server_updated   INTEGER
,server_started   INTEGER
,server_stopped   INTEGER
);
`

var createIndexServerId = `
CREATE INDEX ix_servers_id ON servers (server_id);
`

var createIndexServerState = `
CREATE INDEX ix_servers_state ON servers (server_state);
`

//
// 002_add_column_lastbusy.sql
//

var alterTableServersAddColumnServerLastbusy = `
ALTER TABLE servers ADD COLUMN server_lastbusy INTEGER NOT NULL DEFAULT 0;
`
