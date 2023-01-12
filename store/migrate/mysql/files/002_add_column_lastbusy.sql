-- name: alter-table-servers-add-column-server-lastbusy

ALTER TABLE servers ADD COLUMN server_lastbusy INTEGER NOT NULL DEFAULT 0;
