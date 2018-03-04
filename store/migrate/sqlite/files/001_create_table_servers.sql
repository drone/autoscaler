-- name: create-table-servers

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
,server_ca_key    TEXT
,server_ca_cert   TEXT
,server_tls_key   TEXT
,server_tls_cert  TEXT
,server_created   INTEGER
,server_updated   INTEGER
,server_started   INTEGER
,server_stopped   INTEGER
);

-- name: create-index-server-id

CREATE INDEX IF NOT EXISTS ix_servers_id ON servers (server_id);

-- name: create-index-server-state

CREATE INDEX IF NOT EXISTS ix_servers_state ON servers (server_state);
