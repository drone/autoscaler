-- name: create-table-servers

CREATE TABLE IF NOT EXISTS servers (
 server_name      VARCHAR(50) PRIMARY KEY
,server_id        VARCHAR(250)
,server_provider  VARCHAR(50)
,server_state     VARCHAR(50)
,server_image     VARCHAR(250)
,server_region    VARCHAR(50)
,server_size      VARCHAR(50)
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

-- name: create-index-server-id

CREATE INDEX IF NOT EXISTS ix_servers_id ON servers (server_id);

-- name: create-index-server-state

CREATE INDEX IF NOT EXISTS ix_servers_state ON servers (server_state);
