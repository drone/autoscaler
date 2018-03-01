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
,server_error     MEDIUMTEXT
,server_created   INTEGER
,server_updated   INTEGER
,server_started   INTEGER
,server_stopped   INTEGER
,INDEX(server_id)
,INDEX(server_state)
);
